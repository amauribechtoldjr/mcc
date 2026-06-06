package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type BulkDataResponse struct {
	Object  string     `json:"object"`
	Data    []BulkItem `json:"data"`
	HasMore bool       `json:"has_more"`
}

type BulkItem struct {
	Id              string `json:"id"`
	Uri             string `json:"uri"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	DownloadUri     string `json:"download_uri"`
	UpdatedAt       string `json:"updated_at"`
	Size            int    `json:"size"`
	ContentType     string `json:"content_type"`
	ContentEncoding string `json:"content_encoding"`
}

func main() {
	downloadURI, err := getBulkDownloadURL("default_cards")
	if err != nil {
		panic(err)
	}

	fmt.Println("downloading from:", downloadURI)

	err = downloadFile(downloadURI, "cards.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("file saved at ./cards.json")
}

func getBulkDownloadURL(bulkType string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(
		"GET",
		"https://api.scryfall.com/bulk-data",
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set(
		"User-Agent",
		"my-magic-collection/0.1 (contact: amauribechtoldjr@gmail.com)",
	)
	req.Header.Set(
		"Accept",
		"application/json;q=0.9",
	)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %s", response.Status)
	}

	var bulkData BulkDataResponse

	err = json.NewDecoder(response.Body).Decode(&bulkData)
	if err != nil {
		return "", err
	}

	for _, item := range bulkData.Data {
		if item.Type == bulkType {
			return item.DownloadUri, nil
		}
	}

	return "", fmt.Errorf("bulk type %s not found", bulkType)
}

func downloadFile(url, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
