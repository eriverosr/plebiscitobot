package main

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/silentsokolov/go-vimeo/vimeo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// BOT_NAME is the Bot Name
const BOT_NAME = "plebiscitobot"
const DateFormat = "2006-01-02 15:04:05"

func init() {
	viper.SetConfigName(fmt.Sprintf("%s-config", BOT_NAME))
	viper.AddConfigPath("/etc/tg-bots/")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config file: %s \n", err))
	}
}

func main() {
	var config Config
	viper.Unmarshal(&config)

	// Last video date
	latestDate, err := time.Parse(DateFormat, config.Vimeo.LatestDate)
	if err != nil {
		latestDate, _ = time.Parse(DateFormat, "2020-09-25 00:00:00")
	}
	// Telegram
	tgBot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		panic(err)
	}
	tgBot.Debug = true

	// Vimeo
	if config.Vimeo.Active {
		tc := oauth2.NewClient(oauth2.NoContext, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: config.Vimeo.Token},
		))
		vimeoClient := vimeo.NewClient(tc, nil)
		for {
			log.Printf("Sleeping %d minutes...", config.Vimeo.CheckInterval)
			time.Sleep(time.Duration(config.Vimeo.CheckInterval) * time.Minute)
			log.Printf("Looking for new videos...")
			videos, _, err := vimeoClient.Users.ListVideo(config.Vimeo.UserID, vimeo.OptSort("date"), vimeo.OptDirection("desc"))
			if err != nil {
				log.Printf("error: %s", err)
				continue
			}
			toPublish := make([]*vimeo.Video, 0)
			for _, video := range videos {
				if video.ReleaseTime.After(latestDate) {
					toPublish = append(toPublish, video)
					continue
				}
				break
			}
			log.Printf("%d videos found! Publishing...", len(toPublish))
			correct := 0
			for i := len(toPublish) - 1; i >= 0; i-- {
				msg := tgbotapi.NewMessageToChannel(config.Telegram.Channel, toPublish[i].Link)
				msg.DisableNotification = true
				m, err := tgBot.Send(msg)
				if err != nil {
					log.Printf("cannot send message: %s", err)
					log.Printf("Response: %+v", m)
					continue
				}
				latestDate = toPublish[i].ReleaseTime
				correct++
				viper.Set("vimeo.latestDate", latestDate.Format(DateFormat))
				viper.WriteConfig()
			}
			log.Printf("Published %d videos", correct)
			log.Printf("New latest date: %s", latestDate.Format(DateFormat))
		}
	}
}
