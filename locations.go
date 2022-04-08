package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type Location struct {
	Name          string
	Role          string
	UserId        uuid.UUID
	LocationId    uuid.UUID
	GatewaySerial string
}

func (l Location) String() string {
	return fmt.Sprintf(
		"Name         : %s\nRole         : %s\nUserId       : %s\nLocationId   : %s\nGatewaySerial: %s",
		l.Name, l.Role, l.UserId, l.LocationId, l.GatewaySerial)
}

func locationsCommandLine() *cli.Command {
	return &cli.Command{
		Name:  "locations",
		Usage: "get locations",
		Action: func(c *cli.Context) error {
			return credentials(c, locations)
		},
	}
}

func locations(ctx context.Context, cliContext *cli.Context,
	cfg *clientcredentials.Config) error {

	r, err := cfg.Client(ctx).Get("https://sdk.iotiliti.cloud/homely/locations")
	if err != nil {
		return errors.WithStack(err)
	}
	defer r.Body.Close()

	l := []Location{}
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("Locations\n")

	for i := range l {
		fmt.Printf("\n%s\n", l[i])
	}

	return nil
}
