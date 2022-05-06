package main

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
)

type fileResource struct {
	url.URL
}

func (r *fileResource) Await(context.Context) error {
	// Unify absolute and relative file paths
	filePath := filepath.Join(r.URL.Host, r.URL.Path)

	opts := parseFragment(r.URL.Fragment)

	_, err := os.Stat(filePath)
	if _, ok := opts["absent"]; ok {
		if err == nil {
			return &unavailabilityError{errors.New("file exists")}
		} else if os.IsNotExist(err) {
			return nil
		}
	} else {
		if err == nil {
			return nil
		} else if os.IsNotExist(err) {
			return &unavailabilityError{err}
		}
	}

	return err
}
