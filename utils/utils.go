package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type conf struct {
	Telegram Telegram `toml:"telegram"`
	Mastodon Mastodon `toml:"mastodon"`
}

type Telegram struct {
	Token  string `toml:"token"`
	ChatID int64  `toml:"chat_id"`
}

type Mastodon struct {
	Server      string `toml:"base_url"`
	ClientID    string `toml:"key"`
	CientSecret string `toml:"secret"`
	AccessToken string `toml:"token"`
}

type Msg struct {
	Text  string
	Media [][]byte
}

var Config conf

func LoadConfig() {
	f := "config.toml"
	if _, err := os.Stat(f); os.IsNotExist(err) {
		log.Fatal("Config file not found")
	}

	_, err := toml.DecodeFile(f, &Config)
	if err != nil {
		log.Fatal(err)
	}
}

func PrintPreamable() {
	license := `Copyright (C) 2024  Estela ad Astra "saeziae"
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions.
`
	fmt.Print(license)
}
