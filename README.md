# discord-gcloud

A simple Discord bot to run gcloud commands.

## Configuration 

This is configured using environment variables, which can be seen in `.env` in the root of this repo.

* `DISCORD_TOKEN` - The token for your Discord bot.
* `DISCORD_APP_ID` - The ID of your Discord bot.
* `DISCORD_GUILD_ID` - The ID of the Discord guild you want to use this bot in.

## Usage

This bot is designed to be used in a single Discord server with a single project. It is not designed to be used in multiple servers, or with multiple projects.

This bot uses Discord slash commands, which are only available in servers with the "applications.commands" permission enabled.

### Permissions

Permissions for this bot are governed using Discord roles which are mapped to internal user permissions. The following groups are available:

| Group       | Description                |
| ----------- | -------------------------- |
| `superuser` | Can run any command.       |
| `instances` | Can run instance commands. |

### Commands

| Command        | Description            | Permissions |
| -------------- | ---------------------- | ----------- |
| `/gcloud`      | Runs a gcloud command. | `superuser` |
| `/gcloud list` | Lists all instances    | `instances` |

## License

GPL-3.0