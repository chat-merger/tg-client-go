package telegrambot

import (
	"context"
	"fmt"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"time"
)

func (c *Client) Run(ctx context.Context) error {
	go c.contextCancelHandler(ctx)
	// Start receiving updates.
	err := c.updater.StartPolling(c.bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("start polling: %s", err)
	}
	err = c.listenServerMessages()
	if err != nil {
		return fmt.Errorf("listing Server: %s", err)
	}
	//c.updater.Idle()
	return nil
}

func (c *Client) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		c.updater.Stop()
	}
}
