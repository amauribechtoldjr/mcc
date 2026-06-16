package scryfall

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/amauribechtoldjr/mcc/internal/core/domain"
	"github.com/amauribechtoldjr/mcc/internal/core/port"
	"github.com/google/uuid"
)

const (
	bulkDataURL  = "https://api.scryfall.com/bulk-data"
	bulkDataType = "default_cards"
	importDir    = "C:/dev/mcc/imports"
)

type bulkDataResponse struct {
	Data []bulkItem `json:"data"`
}

type bulkItem struct {
	Type        string `json:"type"`
	DownloadUri string `json:"download_uri"`
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
	ID             string                `json:"id"`
	OracleID       string                `json:"oracle_id"`
	Name           string                `json:"name"`
	Lang           string                `json:"lang"`
	ReleasedAt     string                `json:"released_at"`
	Layout         string                `json:"layout"`
	ImageURIs      scryfallCardImageURIs `json:"image_uris"`
	CMC            float32               `json:"cmc"`
	Colors         []string              `json:"colors"`
	ColorIdentity  []string              `json:"color_identity"`
	ColorIndicator []string              `json:"color_indicator"`
}

type cardSource struct {
	client    *http.Client
	userAgent string
}

func NewCardSource(client *http.Client, userAgent string) port.CardSource {
	return &cardSource{client: client, userAgent: userAgent}
}

func (s *cardSource) bulkFilePath() string {
	return filepath.Join(importDir, fmt.Sprintf("scryfall-%s.json", bulkDataType))
}

func (s *cardSource) GetBulkFileIfExists() (string, bool) {
	filePath := s.bulkFilePath()

	_, err := os.Stat(filePath)
	if err != nil {
		return "", false
	}

	return filePath, true
}

func (s *cardSource) Download(ctx context.Context) (string, error) {
	downloadURI, err := s.bulkDownloadURL(ctx, bulkDataType)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURI, nil)
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

	finalPath := s.bulkFilePath()

	tmp, err := os.CreateTemp(importDir, fmt.Sprintf("scryfall-%s-*.json.tmp", bulkDataType))
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

func (s *cardSource) bulkDownloadURL(ctx context.Context, bulkType string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, bulkDataURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "application/json;q=0.9")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("scryfall bulk-data error: %s", resp.Status)
	}

	var bulk bulkDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&bulk); err != nil {
		return "", err
	}

	for _, item := range bulk.Data {
		if item.Type == bulkType {
			return item.DownloadUri, nil
		}
	}

	return "", fmt.Errorf("scryfall bulk type %q not found", bulkType)
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

		oracleID, err := uuid.Parse(sc.OracleID)
		if err != nil {
			return nil, err
		}

		card := domain.Card{
			OracleID: oracleID,
		}

		mtgCard := domain.MTGCard{
			Name:           sc.Name,
			Layout:         sc.Layout,
			CMC:            sc.CMC,
			ColorIdentity:  sc.ColorIdentity,
			ColorIndicator: sc.ColorIndicator,
			Colors:         sc.Colors,
			ImgSmallURI:    sc.ImageURIs.Small,
			ImgNormalURI:   sc.ImageURIs.Normal,
		}

		cards = append(cards, domain.ImportCard{Card: card, MTGCard: mtgCard})
	}

	return cards, nil
}
