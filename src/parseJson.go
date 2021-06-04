package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

type Src struct {
	Original  string `json:"original"`
	Large2x   string `json:"large2x"`
	Large     string `json:"large"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type Photo struct {
	Id               int    `json:"id"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	Url              string `json:"url"`
	Photographer     string `json:"photographer"`
	Photographer_url string `json:"photographer_url"`
	Photographer_id  int    `json:"photographer_id"`
	Avg_color        string `json:"avg_color"`
	Src              Src    `json:"src"`
	Liked            bool   `json:"liked"`
}

type Photos struct {
	Page          int     `json:"page"`
	Per_page      int     `json:"per_page"`
	Photos        []Photo `json:"photos"`
	Total_results int     `json:"total_results"`
	Next_page     string  `json:"next_page"`
}

/** Parse the json response */
func parseJson(body io.ReadCloser) (Photos, error) {
	var photos Photos
	err := json.NewDecoder(body).Decode(&photos)
	log.WithFields(log.Fields{"Photos": photos}).Debug("parseJson")
	return photos, err
}

/** Pretty print the photos object */
func prettyPrintPhotos(photos Photos) {
	prettyJSON, err := json.MarshalIndent(photos, "", "    ")
	logerr(err)
	fmt.Printf("%s\n", string(prettyJSON))
}

/** Extract the urls of photos according to options */
func getPhotosUrls(photos Photos, src_size string) ([]string, error) {
	urls := make([]string, 0)
	for i := 0; i < photos.Per_page; i++ {
		switch src_size {
		case "Original":
			urls = append(urls, photos.Photos[i].Src.Original)
		case "Large2x":
			urls = append(urls, photos.Photos[i].Src.Large2x)
		case "Large":
			urls = append(urls, photos.Photos[i].Src.Large)
		case "Medium":
			urls = append(urls, photos.Photos[i].Src.Medium)
		case "Small":
			urls = append(urls, photos.Photos[i].Src.Small)
		case "Portrait":
			urls = append(urls, photos.Photos[i].Src.Portrait)
		case "Landscape":
			urls = append(urls, photos.Photos[i].Src.Landscape)
		case "Tiny":
			urls = append(urls, photos.Photos[i].Src.Tiny)
		default:
			return nil, errors.New("Unknown src_size format")
		}
	}
	log.WithFields(log.Fields{"urls": urls}).Debug("getPhotosUrls")
	return urls, nil
}

/** Extract the ids of photos according to options */
func getPhotosIds(photos Photos) ([]int, error) {
	ids := make([]int, 0)
	for i := 0; i < photos.Per_page; i++ {
		ids = append(ids, photos.Photos[i].Id)
	}
	log.WithFields(log.Fields{"ids": ids}).Debug("getPhotosIds")
	return ids, nil
}
