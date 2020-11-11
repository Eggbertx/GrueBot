package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	fp             *frotzPlayer
	cfg            *grueBotConfig
	bot            *discordgo.Session
	quitChannel    string
	mentionMatcher = regexp.MustCompile(`^<@!?(.+)>\s(.+)`)
)

const outputCutset = "\x00 	\n"

func isBadError(err error) bool {
	return err != nil && err != os.ErrClosed && err != io.EOF
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s /path/to/zmachine.z3\n", os.Args[0])
		os.Exit(1)
	}
	var err error

	cfg = new(grueBotConfig)
	if err = cfg.Read(); err != nil {
		fmt.Println("Error reading configuration:", err.Error())
		os.Exit(1)
	}

	if bot, err = discordgo.New("Bot " + cfg.DiscordToken); err != nil {
		fmt.Println("Error setting up Discord bot:", err.Error())
		os.Exit(1)
	}

	if err = bot.Open(); err != nil {
		fmt.Println("Error connecting to Discord: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to Discord as", bot.State.User)

	fp, err := newFrotzPlayer(cfg.FrotzPath, os.Args[1])
	if err != nil {
		bot.ChannelMessageSend(cfg.useChannel, "Error starting z-machine: "+err.Error())
		fmt.Println("Error starting z-machine:", err.Error())
		os.Exit(1)
	}

	defer func() {
		if fp.IsRunning() {
			fp.Kill()
		}
		if bot != nil {
			bot.Close()
		}
	}()

	go func() {
		for {
			if !fp.IsRunning() {
				continue
			}
			output, oErr, eErr := fp.Output()
			if oErr != nil && oErr != os.ErrClosed {
				fmt.Println("oErr: ", oErr.Error())
			}
			if eErr != nil && eErr != os.ErrClosed {
				fmt.Println("eErr: ", eErr.Error())
			}
			output = strings.Trim(output, outputCutset)
			if output == "" {
				continue
			}

			bot.ChannelMessageSend(cfg.useChannel, output)
			fmt.Println(output)
			if quitChannel != "" {
				bot.ChannelMessageSend(quitChannel, "\n\nThank you for playing, I hope you had fun :hearts:")
				bot.Close()
				os.Exit(0)
			}
		}
	}()

	bot.AddHandler(func(session *discordgo.Session, msg *discordgo.MessageCreate) {
		if msg.Author.ID == session.State.User.ID {
			// the bot sent the message, ignore it
			return
		}

		matches := mentionMatcher.FindAllStringSubmatch(msg.Content, 2)
		if len(matches) < 1 || len(matches[0]) < 2 {
			// the message isn't a proper GrueBot Z-Machine command, ignore it
			return
		}

		target := matches[0][1]
		if target != session.State.User.ID {
			// it wasn't addressed to us
			return
		}
		if msg.GuildID != cfg.ServerID {
			fmt.Printf("message from unauthorized server (got %s, should be %s)\n", msg.GuildID, cfg.ServerID)
			return
		}
		authorMention := "<@" + msg.Author.ID + ">"
		command := strings.TrimSpace(matches[0][2])
		words := strings.Split(strings.ToLower(command), " ")
		if !fp.IsRunning() {
			session.ChannelMessageSend(msg.ChannelID, "Something appears to have gone wrong, Frotz isn't running.")
			return
		}
		if len(words) > 0 && words[0] == "quit" {
			if msg.Author.ID != cfg.OwnerID {
				session.ChannelMessageSend(msg.ChannelID,
					authorMention+" You aren't my master. If you want me to quit (e.g. to switch to a different interactive fiction game, then tell <@"+cfg.OwnerID+">")
				return
			}
			fp.Input("quit")
			quitChannel = msg.ChannelID
		}
		fmt.Println("command:", command)
		err := fp.Input(command)
		if err != nil {
			bot.ChannelMessageSend(msg.ChannelID, "Error sending input to z-machine: "+err.Error())
		}
	})

	if err = fp.Run(); err != nil {
		fmt.Println(err.Error())
	}
}
