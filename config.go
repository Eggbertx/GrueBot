package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bwmarrin/discordgo"
)

type grueBotConfig struct {
	FrotzPath             string `json:"frotzPath"`
	DiscordToken          string `json:"discordToken"`
	ServerID              string `json:"serverID"`
	GameChannel           string `json:"gameChannel"`
	TestingChannel        string `json:"testingChannel"`
	OwnerID               string `json:"botOwnerID"`
	TestingMode           bool   `json:"testingMode"`
	useChannel            string
	discordgameChannel    *discordgo.Channel
	discordTestingChannel *discordgo.Channel
}

func (c *grueBotConfig) Read() error {
	jsonBytes, err := ioutil.ReadFile("gruebot.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, c)
	if err != nil {
		return err
	}

	if _, err = os.Stat(c.FrotzPath); err != nil {
		fmt.Printf("frotz not found at given path (%s)\n", c.FrotzPath)
		os.Exit(1)
	}

	if c.TestingMode {
		c.useChannel = c.TestingChannel
	} else {
		c.useChannel = c.GameChannel
	}

	return nil
}

func (c *grueBotConfig) Write() error {
	jsonBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("gruebot.json", jsonBytes, 0777)
}
