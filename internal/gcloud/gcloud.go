// Copyright (C) 2023 Jared Allard <jared@rgst.io>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package gcloud is a wrapper around the gcloud cli
package gcloud

import (
	"context"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Client struct{}

// New returns a new gcloud client
func New() *Client {
	return &Client{}
}

// runCommand runs a gcloud command
func (c *Client) runCommand(ctx context.Context, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "gcloud", args...)
	out, err := cmd.CombinedOutput()
	log.Info().Str("cmd", cmd.String()).Str("output", string(out)).Msg("Ran gcloud command")
	if err != nil {
		return "", errors.Wrapf(err, "failed to run gcloud command: %s", string(out))
	}

	return string(out), nil
}

// ListInstances lists all instances
func (c *Client) ListInstances(ctx context.Context) (string, error) {
	out, err := c.runCommand(ctx, "compute", "instances", "list")
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// StartInstance starts an instance
func (c *Client) StartInstance(ctx context.Context, instance, zone string) (string, error) {
	out, err := c.runCommand(ctx, "compute", "instances", "start", instance, "--zone", zone)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

// StopInstance stops an instance
func (c *Client) StopInstance(ctx context.Context, instance, zone string) (string, error) {
	out, err := c.runCommand(ctx, "compute", "instances", "stop", instance, "--zone", zone)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
