package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func checkInterval(value, min, max int, msg string) error {
	if value < min || value > max {
		return errors.New("Allowed values in [" + strconv.Itoa(min) + ", " + strconv.Itoa(max) + "] Found: " + strconv.Itoa(value))
	}
	return nil
}

func checkPage(page int) error {
	if page < 1 {
		return errors.New("page number must be greater than 1. Found " + strconv.Itoa(page))
	}
	return nil
}

func checkEnum(value, message string, allowed []string) error {
	for _, v := range allowed {
		if value == v {
			return nil
		}
	}
	return errors.New(message + ": unknown value. Wanted: " + strings.Join(allowed, ", ") + ". Found: " + value)
}

func checkNotEmpty(value, msg string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(msg + " must not be empty")
	}
	return nil
}

func checkExists(config string) error {
	if strings.TrimSpace(config) == "" {
		return nil
	}
	info, err := os.Stat(config)
	if os.IsNotExist(err) {
		return err
	}
	if err != nil {
		return err
	}
	if !IsOwnerReadable(info) {
		return errors.New("File is not readable : " + config)
	}
	return nil
}
