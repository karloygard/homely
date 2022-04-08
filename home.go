package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2/clientcredentials"
)

type BoolValue struct {
	LastUpdated time.Time
	Value       bool
}

type FloatValue struct {
	LastUpdated time.Time
	Value       float32
}

type StringValue struct {
	LastUpdated time.Time
	Value       string
}

type AlarmStates struct {
	Alarm  BoolValue
	Tamper BoolValue
}

type TemperatureStates struct {
	Temperature FloatValue
}

type BatteryStates struct {
	Defect  BoolValue
	Low     BoolValue
	Voltage FloatValue
}

type Diagnostic struct {
	NetworklinkAddress  StringValue
	NetworklinkStrength FloatValue
}

type Temperature struct {
	States TemperatureStates
}

type Alarm struct {
	States AlarmStates
}

type Battery struct {
	States BatteryStates
}

type Features struct {
	Alarm       Alarm
	Battery     Battery
	Temperature Temperature
	Diagnostic  Diagnostic
}

type Device struct {
	Name         string
	Id           uuid.UUID
	SerialNumber string
	Location     string
	Online       bool
	ModelId      uuid.UUID
	ModelName    string
	Features     Features
}

type Home struct {
	Name               string
	LocationId         uuid.UUID
	GatewaySerial      string
	AlarmState         string
	UserRoleAtLocation string
	Devices            []Device
}

func homeCommandLine() *cli.Command {
	return &cli.Command{
		Name:  "home",
		Usage: "get home",
		Action: func(c *cli.Context) error {
			return credentials(c, home)
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

func home(ctx context.Context, cliContext *cli.Context,
	cfg *clientcredentials.Config) error {

	r, err := cfg.Client(ctx).Get(
		fmt.Sprintf("https://sdk.iotiliti.cloud/homely/home/%s",
			cliContext.String("location")),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	defer r.Body.Close()

	h := Home{}
	if err := json.NewDecoder(r.Body).Decode(&h); err != nil {
		return errors.WithStack(err)
	}

	bytes, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Println(string(bytes))

	return nil
}
