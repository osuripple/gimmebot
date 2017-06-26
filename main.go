package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var config = struct {
	BotToken string
	MasterID string
	Except   []string
	GuildID  string
}{
	Except: []string{"Ripple Developer"},
}

var bot *discordgo.Session
var roles []*discordgo.Role

func main() {
	loadConfig()

	var err error
	bot, err = discordgo.New("Bot " + config.BotToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		fmt.Println(m.Content)
		if m.Content == "!loadRoles" {
			if m.Author.ID != config.MasterID {
				bot.ChannelMessageSend(m.ChannelID, "You are not my master. How dare you!")
				return
			}
			loadRoles()
			bot.ChannelMessageSend(m.ChannelID, "Roles have been reloaded")
			return
		}
		if !strings.HasPrefix(m.Content, "!gimme ") {
			return
		}
		wantRole := strings.ToLower(m.Content[7:])
		for _, ex := range config.Except {
			if strings.ToLower(ex) == wantRole {
				bot.ChannelMessageSend(m.ChannelID, "Can't give you that, sorry.")
				return
			}
		}
		for _, role := range roles {
			if strings.ToLower(role.Name) == wantRole {
				bot.GuildMemberRoleAdd(config.GuildID, m.Author.ID, role.ID)
				bot.ChannelMessageSend(m.ChannelID, "Should have been given!")
				return
			}
		}
		bot.ChannelMessageSend(m.ChannelID, "Couldn't find that role, sorry.")
	})
	err = bot.Open()
	if err != nil {
		fmt.Println(err)
	}
	for {
		time.Sleep(time.Hour * 24 * 365)
	}
}

func loadConfig() {
	f, err := os.Open("gimmebot.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewDecoder(f).Decode(&config)
	f.Close()
}

func loadRoles() {
	roles, _ = bot.GuildRoles(config.GuildID)
}
