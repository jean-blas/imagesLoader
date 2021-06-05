/**
* Load some images from pexels.com
*
* https://api.pexels.com/v1/search?query=people
*
* Run with command line options
* go run main.go parseJson.go request.go checks.go config.go fileutil.go -q people -p 2 -n 5 -s Small -z small -f /tmp -l info
*
* Run with YAML config file
* go run main.go parseJson.go request.go checks.go config.go fileutil.go -c config.yml
 */
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const URL string = "https://api.pexels.com"

var Usage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), "Utility program used to download some photos automatically from the %s website\n\n", URL)
	flag.PrintDefaults()
}

// TODO read more images than the per_page which is limited to 80 (increase automatically the "page")

func main() {

	flag.Usage = Usage
	orientation := flag.String("o", "landscape", "Orientation (landscape, portrait or square)")
	out_folder := flag.String("f", "/tmp", "Output folder to store the downloaded images")
	page := flag.Int("p", 1, "Page to display")
	per_page := flag.Int("n", 20, "Number of results per page")
	query := flag.String("q", "people", "Query (people, nature, ...)")
	size := flag.String("z", "small", "Minimum photo size (large (24MP), medium(12MP), small(4MP))")
	src_size := flag.String("s", "Landscape", "Size of the photo to download (Original, Large2x, Large, Medium, Small, Portrait, Landscape, Tiny)")
	conf := flag.String("c", "", "Configuration file to use")
	level := flag.String("l", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	c := new(Conf)
	if *conf != "" {
		logerr(checkExists(*conf))
		var err error = nil
		c, err = c.parseConf(*conf)
		logerr(err)
	} else {
		qry := strings.Split(*query, ",")
		c.updateConf(*orientation, *out_folder, *size, *src_size, *level, *conf, qry, *page, *per_page)
	}
	logerr(c.checkConf())

	loginit(c.Level)

	request(c)
}

// Initialize the log level
func loginit(level string) {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

// Utility to write an error and exit
func logerr(err error) {
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// Launch the queries in //
func request(c *Conf) {
	var wg sync.WaitGroup
	wg.Add(len(c.Query))
	for _, q := range c.Query {
		go func(q string) {
			download(c, q)
			wg.Done()
		}(q)
	}
	wg.Wait()
}

// Download some photos in // according to the options
// q : the query keyword to search for
func download(c *Conf, q string) {
	// Build the request according to the command line options
	req, err := buildRequest(URL, q, c.Orientation, c.Size, c.Page, c.Per_page)
	logerr(err)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	logerr(err)
	defer resp.Body.Close()

	// Decode the request response into a json Photo object
	photos, err := parseJson(resp.Body)
	logerr(err)

	// Get the photos urls to download
	urls, err := getPhotosUrls(photos, c.Src_size)
	logerr(err)

	// Get the photos ids for filenames to save on disk
	ids, err := getPhotosIds(photos)
	logerr(err)

	// Simple closure to build the image title
	title := func(i int) string {
		return c.Out_folder + "/" + strings.Join(strings.Split(q, " "), "_") + "_p" + strconv.Itoa(c.Page) + "_" + strconv.Itoa(ids[i]) + ".jpg"
	}

	// Download the photos in parallel
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for i := range urls {
		go func(i int) {
			req, err := buildLoadRequest(urls[i])
			if err == nil {
				res, err := client.Do(req)
				if err == nil {
					defer res.Body.Close()
					data, err := ioutil.ReadAll(res.Body)
					if err == nil {
						log.WithFields(log.Fields{"i": i, "title": title(i)}).Info("Saving file")
						ioutil.WriteFile(title(i), data, 0666)
					} else {
						log.WithFields(log.Fields{"title": title(i)}).Warn(err)
					}
				} else {
					log.WithFields(log.Fields{"title": title(i)}).Warn(err)
				}
			} else {
				log.WithFields(log.Fields{"title": title(i)}).Warn(err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
