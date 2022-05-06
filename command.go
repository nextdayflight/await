package main

import (
	"context"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

type commandResource struct {
	url.URL
}

func (r *commandResource) Await(ctx context.Context) error {
	cmdString, err := url.QueryUnescape(r.URL.Path)
	if err != nil {
		return err
	}
	// TODO(uwe): Splitting by space is brittle
	cmdParts := strings.SplitN(cmdString, " ", 2)
	if len(cmdParts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := cmdParts[0]
	args := cmdParts[1:]

	if err := exec.CommandContext(ctx, cmd, args...).Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &unavailabilityError{exitErr}
		}
		return err
	}

	return nil
}
