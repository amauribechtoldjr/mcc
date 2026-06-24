package scryfall

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

const (
	scryfallAPIUrl = "https://api.scryfall.com"
	importDir      = "C:/dev/mcc/imports"
	dateFileFormat = "2006-01-02_15-04-05.000"
)

type bulkDataResponse struct {
	Data []bulkItem `json:"data"`
}

type bulkItem struct {
	Type        string `json:"type"`
	DownloadUri string `json:"download_uri"`
	UpdatedAt   string `json:"updated_at"`
}

type setsResponse struct {
	Data []scryfallSet `json:"data"`
}

type scryfallSet struct {
	ID            string `json:"id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	SetType       string `json:"set_type"`
	ReleasedAt    string `json:"released_at"`
	ParentSetCode string `json:"parent_set_code"`
	CardCount     int16  `json:"card_count"`
	PrintedSize   int16  `json:"printed_size"`
	Digital       bool   `json:"digital"`
	IconSvgURI    string `json:"icon_svg_uri"`
}

type scryfallCardImageURIs struct {
	Small      string `json:"small"`
	Normal     string `json:"normal"`
	Large      string `json:"large"`
	Png        string `json:"png"`
	ArtCrop    string `json:"art_crop"`
	BorderCrop string `json:"border_crop"`
}

type scryfallCard struct {
	ID              string                `json:"id"`
	OracleID        string                `json:"oracle_id"`
	Name            string                `json:"name"`
	Lang            string                `json:"lang"`
	ReleasedAt      string                `json:"released_at"`
	Layout          string                `json:"layout"`
	ImageURIs       scryfallCardImageURIs `json:"image_uris"`
	CMC             float32               `json:"cmc"`
	Colors          []string              `json:"colors"`
	ColorIdentity   []string              `json:"color_identity"`
	ColorIndicator  []string              `json:"color_indicator"`
	Set             string                `json:"set"`
	SetID           string                `json:"set_id"`
	Games           []string              `json:"games"`
	CollectorNumber string                `json:"collector_number"`
	PrintedName     string                `json:"printed_name"`
	PrintedTypeLine string                `json:"printed_type_line"`
	PrintedText     string                `json:"printed_text"`
	FlavorText      string                `json:"flavor_text"`
}

type cardSource struct {
	client    *http.Client
	userAgent string
}

func NewCardSource(client *http.Client, userAgent string) port.ScryfallCardSource {
	return &cardSource{client: client, userAgent: userAgent}
}

func (s *cardSource) GetSets(ctx context.Context) ([]domain.MTGSet, error) {
	resp, err := s.execScryfallRequest(ctx, scryfallAPIUrl+"/sets")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sets setsResponse
	if err := json.NewDecoder(resp.Body).Decode(&sets); err != nil {
		return nil, err
	}

	mtgSets := make([]domain.MTGSet, 0, len(sets.Data))
	for _, ss := range sets.Data {
		id, err := uuid.Parse(ss.ID)
		if err != nil {
			return nil, err
		}

		if ss.Digital {
			continue
		}

		parsedReleaseAt, err := time.Parse(time.DateOnly, ss.ReleasedAt)
		if err != nil {
			return nil, err
		}

		mtgSets = append(mtgSets, domain.MTGSet{
			ImportID:      id,
			Code:          ss.Code,
			Name:          ss.Name,
			ReleasedAt:    parsedReleaseAt,
			ParentSetCode: ss.ParentSetCode,
			CardCount:     int32(ss.CardCount),
		})
	}

	return mtgSets, nil
}

func (s *cardSource) GetBulkFileIfExists(updatedAt time.Time) (string, bool) {
	filePath := s.bulkFilePath(updatedAt)

	_, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}

	return filePath, true
}

func (s *cardSource) Download(ctx context.Context, bulkData domain.ScryfallBulkData) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkData.DonwloadURI, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", s.userAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scryfall download error: %s", resp.Status)
	}

	finalPath := s.bulkFilePath(bulkData.UpdatedAt)

	tmp, err := os.CreateTemp(importDir, fmt.Sprintf("scryfall_%s_*.json.tmp", bulkData.UpdatedAt.Format(dateFileFormat)))
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return "", err
	}

	if err := tmp.Close(); err != nil {
		return "", err
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		return "", err
	}

	return finalPath, nil
}

func (s *cardSource) ReadCards(ctx context.Context, filePath string, limit int) ([]domain.ImportCard, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Consume the opening '[' of the array.
	if _, err := decoder.Token(); err != nil {
		return nil, err
	}

	cards := make([]domain.ImportCard, 0, limit)
	for decoder.More() && len(cards) < limit {
		var sc scryfallCard
		if err := decoder.Decode(&sc); err != nil {
			return nil, err
		}

		if sc.OracleID == "" {
			continue
		}

		// só cartas fisicamente colecionáveis (descarta produtos digitais: arena/mtgo)
		if !slices.Contains(sc.Games, "paper") {
			continue
		}

		oracleID, err := uuid.Parse(sc.OracleID)
		if err != nil {
			return nil, err
		}

		card := domain.Card{}

		setId := uuid.MustParse(sc.SetID)

		name := sc.Name
		if sc.PrintedName != "" {
			name = sc.PrintedName
		}

		mtgCard := domain.MTGCard{
			OracleID:        oracleID,
			Name:            name,
			PrintedTypeLine: sc.PrintedTypeLine,
			PrintedText:     sc.PrintedText,
			FlavorText:      sc.FlavorText,
			Lang:            sc.Lang,
			CollectorNumber: sc.CollectorNumber,
			Layout:          sc.Layout,
			CMC:             sc.CMC,
			ColorIdentity:   sc.ColorIdentity,
			ColorIndicator:  sc.ColorIndicator,
			Colors:          sc.Colors,
			ImgSmallURI:     sc.ImageURIs.Small,
			ImgNormalURI:    sc.ImageURIs.Normal,
			SetID:           setId,
			SetCode:         sc.Set,
		}

		cards = append(cards, domain.ImportCard{Card: card, MTGCard: mtgCard})
	}

	return cards, nil
}

func (s *cardSource) execScryfallRequest(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "application/json;q=0.9")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scryfall request error: %s", resp.Status)
	}

	return resp, nil
}

func (s *cardSource) GetBulkData(ctx context.Context) ([]domain.ScryfallBulkData, error) {
	resp, err := s.execScryfallRequest(ctx, scryfallAPIUrl+"/bulk-data")
	if err != nil {
		return nil, err
	}

	var bulk bulkDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&bulk); err != nil {
		return nil, err
	}

	bulkData := make([]domain.ScryfallBulkData, len(bulk.Data))
	for i, item := range bulk.Data {

		updatedAt, err := time.Parse(time.RFC3339Nano, item.UpdatedAt)
		if err != nil {
			return nil, err
		}

		bulkData[i] = domain.ScryfallBulkData{
			Type:        item.Type,
			UpdatedAt:   updatedAt,
			DonwloadURI: item.DownloadUri,
		}
	}

	return bulkData, nil
}

func (s *cardSource) bulkFilePath(updatedAt time.Time) string {
	return filepath.Join(importDir, fmt.Sprintf("scryfall_%s.json", updatedAt.Format(dateFileFormat)))
}
