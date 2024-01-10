package vkontaktebot

import (
	"context"
	"fmt"
	"log"
)

func (c *Client) Run(ctx context.Context) error {
	go c.contextCancelHandler(ctx)
	// Start receiving updates.
	// get information about the group

	// Run Bots Long Poll
	go func() {
		log.Println("Start Long Poll")
		if err := c.lp.Run(); err != nil {
			log.Fatalf("lp.Run: %s", err)
		}
	}()

	err := c.listenServerMessages()
	if err != nil {
		return fmt.Errorf("listing Server: %s", err)
	}
	//c.updater.Idle()
	return nil
}

func (c *Client) contextCancelHandler(ctx context.Context) {
	select {
	case <-ctx.Done():
		c.lp.Shutdown()
	}
}
