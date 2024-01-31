package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type Config struct {
	ServerHost string
	TgBotToken string
	TgChat     int64
	TgXApiKey  string
	VkBotToken string
	VkPeer     int
	VkXApiKey  string
	DbFile     string
	RedisUrl   string
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
	cfgFs.fs.StringVar(&cfgFs.cfg.ServerHost, flagServerHost, "", "host for connect to chat-merger-server")
	cfgFs.fs.StringVar(&cfgFs.cfg.TgBotToken, flagTgBotToken, "", "a string or line consisting of letters and numbers which is require in order for authorizing the bot and for sending requests to the Bot API")
	cfgFs.fs.Int64Var(&cfgFs.cfg.TgChat, flagTgChat, 0, "id of telegram chat")
	cfgFs.fs.StringVar(&cfgFs.cfg.TgXApiKey, flagTgXApiKey, "", "api key of chat-merger telegram adapter")
	cfgFs.fs.StringVar(&cfgFs.cfg.VkBotToken, flagVkBotToken, "", "a string or line consisting of letters and numbers which is require in order for authorizing the bot and for sending requests to the Bot API")
	cfgFs.fs.IntVar(&cfgFs.cfg.VkPeer, flagVkPeer, 0, "id of vkontakte peer")
	cfgFs.fs.StringVar(&cfgFs.cfg.VkXApiKey, flagVkXApiKey, "", "api key of chat-merger vkontakte adapter")
	cfgFs.fs.StringVar(&cfgFs.cfg.DbFile, flagDbFile, "", "path to sqlite database source")
	cfgFs.fs.StringVar(&cfgFs.cfg.RedisUrl, flagRedisUrl, "", "url for connect to redis")
	return cfgFs
}

// cleanLastCfg clean parsed values
func (c *FlagSet) cleanLastCfg() {
	c.cfg.ServerHost = ""
	c.cfg.TgBotToken = ""
	c.cfg.TgChat = 0
	c.cfg.TgXApiKey = ""
	c.cfg.VkBotToken = ""
	c.cfg.VkPeer = 0
	c.cfg.VkXApiKey = ""
	c.cfg.DbFile = ""
	c.cfg.RedisUrl = ""
}

// Flag names:

const (
	flagServerHost = "host"
	flagTgBotToken = "tg-token"
	flagTgChat     = "tg-chat-id"
	flagTgXApiKey  = "tg-x-api-key"
	flagVkBotToken = "vk-token"
	flagVkPeer     = "vk-peer-id"
	flagVkXApiKey  = "vk-x-api-key"
	flagDbFile     = "db"
	flagRedisUrl   = "redis-url"
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

	// check what all fields defined
	switch {
	case newCfg.ServerHost == "":
		return nil, missingArgExit(flagServerHost)
	case newCfg.TgBotToken == "":
		return nil, missingArgExit(flagTgBotToken)
	case newCfg.TgChat == 0:
		return nil, missingArgExit(flagTgChat)
	case newCfg.TgXApiKey == "":
		return nil, missingArgExit(flagTgXApiKey)
	case newCfg.VkBotToken == "":
		return nil, missingArgExit(flagVkBotToken)
	case newCfg.VkPeer == 0:
		return nil, missingArgExit(flagVkPeer)
	case newCfg.VkXApiKey == "":
		return nil, missingArgExit(flagVkXApiKey)
	case newCfg.DbFile == "":
		return nil, missingArgExit(flagDbFile)
	case newCfg.RedisUrl == "":
		return nil, missingArgExit(flagRedisUrl)
	}

	return &newCfg, nil
}

var WrongArgumentError = errors.New("wrong argument")
