package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
)

type httpResource struct {
	url.URL
}

func (r *httpResource) Await(ctx context.Context) error {
	var client *http.Client

	if skipTLSVerification(r) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}

	// IDEA(uwe): Use fragment to set method

	req, err := http.NewRequest("GET", r.URL.String(), nil)
	if err != nil {
		return err
	}

	// IDEA(uwe): Use k/v pairs in fragment to set headers

	req = req.WithContext(ctx)

	req.Header.Set("User-Agent", "await/"+version)

	resp, err := client.Do(req)
	if err != nil {
		return &unavailabilityError{err}
	}
	defer func() { _ = resp.Body.Close() }()

	// IDEA(uwe): Use fragment to set tolerated status code

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return &unavailabilityError{errors.New(resp.Status)}
}

func skipTLSVerification(r *httpResource) bool {
	opts := parseFragment(r.URL.Fragment)
	vals, ok := opts["tls"]
	return ok && r.URL.Scheme == "https" && len(vals) == 1 && vals[0] == "skip-verify"
}
