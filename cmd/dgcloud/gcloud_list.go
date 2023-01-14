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

package main

import (
	"context"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
	"github.com/rgst-io/discord-gcloud/internal/gcloud"
)

// gcloudListCommand returns the list command.
func gcloudListCommand(gcli *gcloud.Client) *disgolf.Command {
	return &disgolf.Command{
		Name:        "list",
		Description: "List all instances",
		Handler:     disgolf.HandlerFunc(gcloudListFunc(gcli)),
	}
}

// gcloudListFunc is the handler for when the list command is ran.
func gcloudListFunc(gcli *gcloud.Client) func(ctx *disgolf.Ctx) {
	return func(ctx *disgolf.Ctx) {
		out, err := gcli.ListInstances(context.TODO())
		if err != nil {
			errorReplyInteraction(ctx, err)
			return
		}

		resp := "```" + out + "```"
		ctx.InteractionResponseEdit(ctx.Interaction, &discordgo.WebhookEdit{
			Content: &resp,
		})
	}
}
