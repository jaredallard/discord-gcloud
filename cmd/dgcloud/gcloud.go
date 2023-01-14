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
	"fmt"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"github.com/rgst-io/discord-gcloud/internal/gcloud"
	"github.com/rgst-io/discord-gcloud/internal/permissions"
	"github.com/rs/zerolog/log"
)

// errorReplyWithoutInteraction replies to the interaction with an error.
// This should not be used once the interaction has been responded to,
// e.g. a deferred response.
func errorReplyWithoutInteraction(ctx *disgolf.Ctx, err error) {
	log.Error().Err(err).Msg("Responded with error while handling interaction")
	if err := ctx.Respond(&discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Failed to handle interaction:\n```" + err.Error() + "```",
		},
	}); err != nil {
		log.Error().Err(err).Msg("Failed to respond with error")
	}
}

// errorReplyInteraction replies to the interaction with an error.
func errorReplyInteraction(ctx *disgolf.Ctx, err error) {
	if ctx.Interaction == nil {
		log.Error().Err(err).Msg("Called errorReplyInteraction with no active interaction")
		return
	}

	errStr := err.Error()
	respStr := "Failed to handle interaction:\n```" + errStr + "```"
	log.Error().Err(err).Msg("Responded with error while handling interaction")
	if _, err := ctx.InteractionResponseEdit(ctx.Interaction, &discordgo.WebhookEdit{
		Content: &respStr,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to respond with error")
	}
}

// getUsersGroups returns the groups for a user.
func getUsersGroups(ctx *disgolf.Ctx) ([]string, error) {
	if ctx.Interaction.Member == nil {
		return nil, fmt.Errorf("DMs are not supported")
	}

	user := ctx.Interaction.Member
	if user == nil {
		return nil, fmt.Errorf("no user found")
	}

	return permissions.GetGroupsForUser(user), nil
}

// registerGCloudCommands registers the gcloud commands.
func registerGCloudCommands(router *disgolf.Router) {
	gcli := gcloud.New()

	router.Register(&disgolf.Command{
		Name:        "gcloud",
		Description: "Run a gcloud command",
		Middlewares: []disgolf.Handler{
			disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
				groups, err := getUsersGroups(ctx)
				if err != nil {
					errorReplyWithoutInteraction(ctx, errors.Wrap(err, "Failed to determine your groups"))
					return
				}

				var userId string
				var userName string
				if ctx.Interaction.Member != nil {
					userId = ctx.Interaction.Member.User.ID
					userName = ctx.Interaction.Member.User.Username
				} else {
					if ctx.Interaction.User == nil {
						userId = ctx.Interaction.User.ID
						userName = ctx.Interaction.User.Username
					} else {
						errorReplyWithoutInteraction(ctx, fmt.Errorf("Failed to determine your user id :("))
						return
					}
				}

				log.Info().
					Strs("groups", groups).
					Str("userId", userId).
					Str("userName", userName).
					Msg("Handling interaction")
				if !permissions.CanRunCommand(groups, ctx.Caller.Name) {
					errorReplyWithoutInteraction(ctx, fmt.Errorf("You do not have permission to run this command"))
					return
				}

				if err := ctx.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				}); err != nil {
					errorReplyWithoutInteraction(ctx, errors.Wrap(err, "Failed to inform Discord of our deferred response. This is likely a Discord API issue. Error"))
					return
				}

				ctx.Next()
			}),
		},
		SubCommands: disgolf.NewRouter([]*disgolf.Command{
			gcloudListCommand(gcli),
			gcloudStartCommand(gcli),
			gcloudStopCommand(gcli),
		}),
	})
}
