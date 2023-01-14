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

// Package permissions implements permissions for the bot.
package permissions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// groupToRoleID is a map of groups to role IDs.
var groupToRoleID = map[string]string{
	"superuser": "",
	"instances": "",
}

// SetGroupToRoleID updates the group to role ID map.
func SetGroupToRoleID(group, roleID string) {
	groupToRoleID[group] = roleID
}

// commandToGroups is a map of commands to groups that can run them.
// Wildcards in the command name evaluate to all commands.
var commandToGroups = map[string][]string{
	"*":     {"superuser"},
	"list":  {"instances"},
	"start": {"instances"},
	"stop":  {"instances"},
}

// GetGroupForRoleID returns a group for a role ID.
func GetGroupForRoleID(roleID string) string {
	for g, r := range groupToRoleID {
		// Skip empty groups
		if r == "" {
			continue
		}

		if r == roleID {
			return g
		}
	}

	return ""
}

// GetGroupsForUser returns a list of groups for a user.
func GetGroupsForUser(user *discordgo.Member) []string {
	var groups []string

	for _, r := range user.Roles {
		g := GetGroupForRoleID(r)
		if g == "" {
			continue
		}

		groups = append(groups, g)
	}

	return groups
}

// CanRunCommand checks if a group can run a command.
func CanRunCommand(groups []string, command string) bool {
	usersGroups := map[string]struct{}{}
	for _, g := range groups {
		usersGroups[g] = struct{}{}
	}

	// Evaluate wildcard permissions first
	if groups, ok := commandToGroups["*"]; ok {
		for _, g := range groups {
			if _, ok := usersGroups[g]; ok {
				log.Debug().Str("group", g).Str("command", command).Msg("User has wildcard permission through group")
				return true
			}
		}
	}

	// Evaluate specific permissions
	if groups, ok := commandToGroups[command]; ok {
		for _, g := range groups {
			if _, ok := usersGroups[g]; ok {
				log.Debug().Str("group", g).Str("command", command).Msg("User has permission through group")
				return true
			}
		}
	}

	return false
}
