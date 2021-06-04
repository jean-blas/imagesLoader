package main

import (
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// The request API key
const AUTHORIZATION = "563492ad6f9170000100000193cb293dca544b4c86d74483e6f233e6"

/** Build the request according to options */
func buildRequest(url, query, orientation, size string, page, per_page int) (*http.Request, error) {
	// Parse the command line and create the request
	qry := url
	if strings.TrimSpace(query) != "" {
		qry += "/v1/search?query=" + strings.TrimSpace(query)
		// Add the page
		qry += "&page=" + strconv.Itoa(page)
		// Add the number of results per page
		qry += "&per_page=" + strconv.Itoa(per_page)
		// Add the orientation
		qry += "&orientation=" + orientation
		// Add the minimum photo size
		qry += "&size=" + size
		log.WithFields(log.Fields{"query": query}).Debug("buildRequestt")
		req, err := http.NewRequest("GET", qry, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Authorization", AUTHORIZATION)
		log.WithFields(log.Fields{"request": req}).Debug("buildRequest")
		return req, err
	}
	return nil, nil
}

func buildLoadRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", AUTHORIZATION)
	return req, err
}
