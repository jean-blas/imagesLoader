package main

/**
 * Read a YAML file
 */

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	Orientation string `yaml:"orientation"`
	Out_folder  string `yaml:"out_folder"`
	Query       string `yaml:"query"`
	Size        string `yaml:"size"`
	Src_size    string `yaml:"src_size"`
	Level       string `yaml:"level"`
	Page        int    `yaml:"page"`
	Per_page    int    `yaml:"per_page"`
	Config      string
}

// Update the config object according to the parameters
func (c *Conf) updateConf(orientation, out_folder, query, size, src_size, level, config string, page, per_page int) {
	c.Orientation = orientation
	c.Out_folder = out_folder
	c.Query = query
	c.Page = page
	c.Per_page = per_page
	c.Size = size
	c.Src_size = src_size
	c.Level = level
	c.Config = config
	log.WithFields(log.Fields{"orientation": c.Orientation, "out_folder": c.Out_folder, "query": c.Query,
		"page": c.Page, "per_page": c.Per_page, "size": c.Size, "src_size": c.Src_size, "config": c.Config}).Debug("Options")
}

// Parse the YAML file and extract the options into a config object
func (c *Conf) parseConf(filename string) (*Conf, error) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{"orientation": c.Orientation, "out_folder": c.Out_folder, "query": c.Query,
		"page": c.Page, "per_page": c.Per_page, "size": c.Size, "src_size": c.Src_size, "config": c.Config}).Debug("Configuration file")
	return c, nil
}

// Check that the config fields are valid with respect to the Pexels site
func (c Conf) checkConf() error {
	if err := checkInterval(c.Per_page, 1, 80, "Bad number of results per page (option -n):"); err != nil {
		return err
	}
	if err := checkPage(c.Page); err != nil {
		return err
	}
	if err := checkEnum(c.Src_size, "Src_size", []string{"Original", "Large2x", "Large", "Medium", "Small", "Portrait", "Landscape", "Tiny"}); err != nil {
		return err
	}
	if err := checkEnum(c.Orientation, "Orientation", []string{"landscape", "portrait", "square"}); err != nil {
		return err
	}
	if err := checkEnum(c.Size, "Size", []string{"large", "medium", "small"}); err != nil {
		return err
	}
	if err := checkNotEmpty(c.Out_folder, "Output folder (option -f)"); err != nil {
		return err
	}
	if err := checkEnum(c.Level, "Level", []string{"debug", "info", "warn", "error"}); err != nil {
		return err
	}
	return nil
}
