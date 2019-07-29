package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/Luzifer/rconfig/v2"
)

var (
	cfg = struct {
		LogLevel       string        `flag:"log-level" default:"info" description:"Log level (debug, info, warn, error, fatal)"`
		Message        string        `flag:"message,m" default:"" description:""`
		MQTTBroker     string        `flag:"mqtt-broker,b" default:"tcp://localhost:1883" description:"Broker URI to connect to scheme://host:port (scheme is one of 'tcp', 'ssl', or 'ws')"`
		MQTTClientID   string        `flag:"mqtt-client-id" vardefault:"client-id" description:"Client ID to use when connecting, must be unique"`
		MQTTUser       string        `flag:"mqtt-user,u" default:"" description:"Username to identify against the broker" validate:"nonzero"`
		MQTTPass       string        `flag:"mqtt-pass,p" default:"" description:"Password to identify against the broker" validate:"nonzero"`
		MQTTTimeout    time.Duration `flag:"mqtt-timeout" default:"10s" description:"How long to wait for the client to complete operations"`
		OutputFormat   string        `flag:"output-format,o" default:"log" description:"How to ouptut received messages (One of 'log', 'csv', 'jsonl')"`
		QOS            int           `flag:"qos" default:"1" description:"QOS to use (0 - Only Once, 1 - At Least Once, 2 - Only Once)"`
		Retain         bool          `flag:"retain" default:"false" description:"Retain message on topic"`
		Topics         []string      `flag:"topic,t" default:"" description:"Topic to subscribe / publish to"`
		VersionAndExit bool          `flag:"version" default:"false" description:"Prints current version and exits"`
	}{}

	version = "dev"
)

func init() {
	rconfig.AutoEnv(true)

	rconfig.SetVariableDefaults(map[string]string{
		"client-id": uuid.Must(uuid.NewV4()).String(),
	})

	if err := rconfig.ParseAndValidate(&cfg); err != nil {
		log.Fatalf("Unable to parse commandline options: %s", err)
	}

	if cfg.VersionAndExit {
		fmt.Printf("mqttcli %s\n", version)
		os.Exit(0)
	}

	if l, err := log.ParseLevel(cfg.LogLevel); err != nil {
		log.WithError(err).Fatal("Unable to parse log level")
	} else {
		log.SetLevel(l)
	}
}

func main() {
	var cmd string
	if len(rconfig.Args()) > 1 {
		cmd = rconfig.Args()[1]
	}

	if len(cfg.Topics) == 0 || (len(cfg.Topics) == 1 && cfg.Topics[0] == "") {
		log.Fatal("No topic(s) specified")
	}

	client := mqtt.NewClient(
		mqtt.NewClientOptions().
			AddBroker(cfg.MQTTBroker).
			SetUsername(cfg.MQTTUser).
			SetPassword(cfg.MQTTPass),
	)

	tok := client.Connect()
	if !tok.WaitTimeout(cfg.MQTTTimeout) {
		log.Fatal("Unable to connect within timeout")
	}
	if err := tok.Error(); err != nil {
		log.WithError(err).Fatal("Unable to connect to broker")
	}

	switch cmd {
	case "pub":
		if cfg.Message == "" {
			log.Fatal("Empty message on publish")
		}

		if err := publish(client); err != nil {
			log.WithError(err).Fatal("Failed to publish message")
		}

	case "sub":
		if err := subscribe(client); err != nil {
			log.WithError(err).Fatal("Failed to subscribe and listen")
		}

	default:
		log.Fatal("No command specified. Usage: mqttcli [opts] <pub|sub>")
	}
}
