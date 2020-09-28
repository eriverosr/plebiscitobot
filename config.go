package main

type Config struct {
	Telegram TelegramConfig
	Vimeo    VimeoConfig
	Servel   ServelConfig
}

type TelegramConfig struct {
	Token   string
	Channel string
}

type VimeoConfig struct {
	Active        bool
	Token         string
	UserID        string
	CheckInterval int
	LatestDate    string
}

type ServelConfig struct {
	Active        bool
	URL           string
	CheckInterval int
}
