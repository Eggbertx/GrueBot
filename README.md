# GrueBot
A Discord bot that allows you to use dfrotz to play Z-machines

## Building
To build GrueBot, run `make build`. To build with debugging symbols, run `make build-debug`.

## Running
Gruebot's usage is `./gruebot /path/to/game`.

Here is an example `gruebot.json`. Gruebot needs one in order to run
```JSON
{
	"frotzPath": "/path/to/dfrotz",
	"discordToken": "yourdiscordtokenhere",
	"serverID": "getserveridandputithere",
	"gameChannel": "channelinserver",
	"testingChannel": "channelfortestingiftestingModeistrue",
	"testingMode": false,
	"botOwnerID": "yourDiscordID"
}
```
Note that GrueBot specifically requires a "dumb" build of Frotz, without ncurses support. Using the mainline build will likely not work.
