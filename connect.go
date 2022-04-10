package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type Change struct {
	Feature    string
	LastUpdate time.Time
	StateName  string
	Value      any
}

type Data struct {
	Changes        []Change
	DeviceId       uuid.UUID
	GatewayId      uuid.UUID
	LocationId     uuid.UUID
	ModelId        uuid.UUID
	RootLocationId uuid.UUID
}

type Event struct {
	Data Data
	Type string
}

func connectCommandLine() *cli.Command {
	return &cli.Command{
		Name:  "connect",
		Usage: "connect websocket",
		Action: func(c *cli.Context) error {
			return credentials(c, connect)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "location",
				Aliases: []string{"l"},
				Usage:   "Location id",
			},
		},
	}
}

func connect(ctx context.Context, cliContext *cli.Context,
	cfg *clientcredentials.Config) error {

	token, err := cfg.Token(ctx)
	if err != nil {
		return err
	}
	t := transport.GetDefaultWebsocketTransport()
	t.PingInterval = 20 * time.Second

	c, err := gosocketio.Dial(
		fmt.Sprintf("%s&locationId=%s&token=Bearer%%20%s",
			gosocketio.GetUrl("sdk.iotiliti.cloud", 443, true),
			cliContext.String("location"), token.AccessToken),
		t,
	)
	if err != nil {
		return err
	}
	defer c.Close()

	if err := c.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Connected")
	}); err != nil {
		return err
	}

	if err := c.On(gosocketio.OnError, func(c *gosocketio.Channel) {
		log.Printf("error: %v", c)
	}); err != nil {
		return err
	}

	if err := c.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	}); err != nil {
		return err
	}

	if err := c.On("event", func(c *gosocketio.Channel, arg any) {
		t, _ := json.MarshalIndent(arg, "", "  ")
		log.Printf("event: %v", string(t))
	}); err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}
