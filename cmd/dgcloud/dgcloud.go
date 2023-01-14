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

// Package main implements a Discord bot that can be used to
// to run gcloud commands with nice formatting and other features.
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/FedorLap2006/disgolf"
	"github.com/getoutreach/gobox/pkg/cfg"
	"github.com/rgst-io/discord-gcloud/internal/permissions"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// main is the entry point for the dgcloud command.
func main() {
	token := cfg.SecretData(os.Getenv("DISCORD_TOKEN"))
	guildID := os.Getenv("DISCORD_GUILD_ID")
	appID := os.Getenv("DISCORD_APP_ID")
	logFormat := os.Getenv("LOG_FORMAT")

	superusersRoleID := os.Getenv("SUPERUSERS_ROLE_ID")
	instancesRoleID := os.Getenv("INSTANCES_ROLE_ID")
	if superusersRoleID != "" {
		permissions.SetGroupToRoleID("superuser", superusersRoleID)
	}
	if instancesRoleID != "" {
		permissions.SetGroupToRoleID("instances", instancesRoleID)
	}

	if logFormat != "json" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	bot, err := disgolf.New(string(token))
	if err != nil {
		panic(err)
	}

	registerGCloudCommands(bot.Router)

	bot.AddHandler(bot.Router.HandleInteraction)

	log.Info().Msg("Starting bot...")
	if err := bot.Open(); err != nil {
		panic(err)
	}
	defer bot.Close()

	if err := bot.Router.Sync(bot.Session, appID, guildID); err != nil {
		panic(err)
	}

	stchan := make(chan os.Signal, 1)
	signal.Notify(stchan, syscall.SIGTERM, os.Interrupt, syscall.SIGSEGV)

	// exit on signals
	log.Info().Msg("Bot started")
	<-stchan

	log.Info().Msg("Exiting...")
}
