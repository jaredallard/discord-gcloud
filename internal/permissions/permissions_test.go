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

import "testing"

func TestCanRunCommand(t *testing.T) {
	type args struct {
		groups  []string
		command string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "superuser can run all commands",
			args: args{
				groups:  []string{"superuser"},
				command: "list",
			},
			want: true,
		},
		{
			name: "instances can run list command",
			args: args{
				groups:  []string{"instances"},
				command: "list",
			},
			want: true,
		},
		{
			name: "unknown group cannot run list command",
			args: args{
				groups:  []string{"unknown"},
				command: "list",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanRunCommand(tt.args.groups, tt.args.command); got != tt.want {
				t.Errorf("CanRunCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
