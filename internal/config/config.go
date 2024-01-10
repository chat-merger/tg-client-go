package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	TelegramBotToken string
	ServerHost       string
	TelegramChat     int64
	XApiKey          string
}

// Flag-feature part:

// FlagSet is Config factory
type FlagSet struct {
	cfg Config
	fs  *flag.FlagSet
}

func InitFlagSet() *FlagSet {
	cfgFs := new(FlagSet)
	cfgFs.fs = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cfgFs.fs.StringVar(&cfgFs.cfg.TelegramBotToken, flagTelegramBotToken, "", "a string or line consisting of letters and numbers which is require in order for authorizing the bot and for sending requests to the Bot API")
	cfgFs.fs.StringVar(&cfgFs.cfg.ServerHost, flagServerHost, "", "host for connect to chat-merger-server")
	cfgFs.fs.Int64Var(&cfgFs.cfg.TelegramChat, flagTelegramChat, 0, "id of telegram chat, int64")
	cfgFs.fs.StringVar(&cfgFs.cfg.XApiKey, flagXApiKey, "", "api key of chat-merger client (adapter)")
	return cfgFs
}

// cleanLastCfg clean parsed values
func (c *FlagSet) cleanLastCfg() {
	c.cfg.TelegramBotToken = ""
	c.cfg.ServerHost = ""
}

// Flag names:

const (
	flagTelegramBotToken = "token"
	flagServerHost       = "server-host"
	flagTelegramChat     = "tg-chat-id"
	flagXApiKey          = "x-api-key"
)

// Usage printing "how usage flags" message
func (c *FlagSet) Usage() { c.fs.Usage() }

// Parse is Config factory method
func (c *FlagSet) Parse(args []string) (*Config, error) {
	missingArgExit := func(argName string) error {
		return fmt.Errorf("missing `%s` argument: %w", argName, WrongArgumentError)
	}

	err := c.fs.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("parse given config arguments: %w", err)
	}
	newCfg := c.cfg // copy parsed values
	c.cleanLastCfg()

	if newCfg.TelegramBotToken == "" {
		return nil, missingArgExit(flagTelegramBotToken)
	}

	if newCfg.ServerHost == "" {
		return nil, missingArgExit(flagServerHost)
	}

	return &newCfg, nil
}

var WrongArgumentError = errors.New("wrong argument")
