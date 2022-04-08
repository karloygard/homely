package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	appVersion = "0.0.1"
)

func credentials(cliContext *cli.Context,
	fn func(context.Context, *cli.Context, *clientcredentials.Config) error) error {

	ctx := context.Background()

	cfg := &clientcredentials.Config{
		TokenURL: "https://sdk.iotiliti.cloud/homely/oauth/token",
		EndpointParams: url.Values{
			"username": {cliContext.String("username")},
			"password": {cliContext.String("password")},
		},
	}

	return fn(ctx, cliContext, cfg)
}

func main() {
	app := cli.NewApp()

	app.Version = appVersion
	app.Usage = "a Homely test client"
	app.Commands = []*cli.Command{
		locationsCommandLine(),
		homeCommandLine(),
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "username",
			Aliases: []string{"u"},
			Usage:   "Username",
		},
		&cli.StringFlag{
			Name:    "password",
			Aliases: []string{"p"},
			Usage:   "Password",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
