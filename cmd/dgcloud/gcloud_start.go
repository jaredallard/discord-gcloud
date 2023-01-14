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
	"fmt"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rgst-io/discord-gcloud/internal/gcloud"
)

// gcloudStartCommand returns the start command.
func gcloudStartCommand(gcli *gcloud.Client) *disgolf.Command {
	return &disgolf.Command{
		Name:        "start",
		Description: "Start an instance",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "instance",
				Description: "The instance to start",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "zone",
				Description: "The zone to start the instance in",
				Required:    true,
			},
		},
		Handler: disgolf.HandlerFunc(gcloudStartFunc(gcli)),
	}
}

// gcloudStartFunc is the handler for when the start command is ran.
func gcloudStartFunc(gcli *gcloud.Client) func(ctx *disgolf.Ctx) {
	return func(ctx *disgolf.Ctx) {
		instanceOpt, ok := ctx.Options["instance"]
		if !ok {
			errorReplyInteraction(ctx, fmt.Errorf("no instance provided"))
			return
		}
		instance, ok := instanceOpt.Value.(string)
		if !ok {
			errorReplyInteraction(ctx, fmt.Errorf("invalid instance provided (not a string)"))
			return
		}

		zoneOpt, ok := ctx.Options["zone"]
		if !ok {
			errorReplyInteraction(ctx, fmt.Errorf("no zone provided"))
			return
		}
		zone, ok := zoneOpt.Value.(string)
		if !ok {
			errorReplyInteraction(ctx, fmt.Errorf("invalid zone provided (not a string)"))
			return
		}

		if err := ctx.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		}); err != nil {
			errorReplyInteraction(ctx, errors.Wrap(err, "failed to inform discord we'll reply later, instance has been started but we can't inform you when it's done starting"))
		}

		if _, err := gcli.StartInstance(context.TODO(), instance, zone); err != nil {
			errorReplyInteraction(ctx, err)
			return
		}

		resp := fmt.Sprintf("Started instance %s in zone %s", instance, zone)
		if _, err := ctx.InteractionResponseEdit(ctx.Interaction, &discordgo.WebhookEdit{
			Content: &resp,
		}); err != nil {
			errorReplyInteraction(ctx, errors.Wrap(err, "failed to respond to interaction"))
			return
		}
	}
}
