package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func subscribe(client mqtt.Client) (err error) {
	var (
		callback mqtt.MessageHandler
		done     chan struct{}
		timeout  <-chan time.Time
		topics   = make(map[string]byte)
	)

	for _, t := range cfg.Topics {
		topics[t] = byte(cfg.QOS) //#nosec:G115 // QOS is expected to be 0,1,2 - fine to convert
	}

	switch cfg.OutputFormat {
	case "log":
		callback = subscribeCallbackLog

	case "csv":
		fmt.Println("Topic,QOS,Retained,Message") //nolint:forbidigo // fine for CSV print
		callback = subscribeCallbackCSV

	case "jsonl":
		callback = subscribeCallbackJSONL

	default:
		return fmt.Errorf("invalid output format %q", cfg.OutputFormat)
	}

	if cfg.ReceiveOnce {
		var once sync.Once

		done = make(chan struct{})
		outputCallback := callback
		callback = func(client mqtt.Client, msg mqtt.Message) {
			once.Do(func() {
				outputCallback(client, msg)
				close(done)
			})
		}
	}

	if err = mqttTokenToError(client.SubscribeMultiple(topics, callback)); err != nil {
		return fmt.Errorf("subscribing topics: %w", err)
	}

	if cfg.Timeout > 0 {
		timeout = time.NewTimer(cfg.Timeout).C
	}

	select {
	case <-done:
	case <-timeout:
	}

	return nil
}

func subscribeCallbackCSV(_ mqtt.Client, msg mqtt.Message) {
	//nolint:forbidigo // fine for CSV print
	fmt.Printf(
		"%s,%d,%v,%q\n",
		msg.Topic(),
		msg.Qos(),
		msg.Retained(),
		string(msg.Payload()),
	)
}

func subscribeCallbackJSONL(_ mqtt.Client, msg mqtt.Message) {
	jsonMessage := struct {
		Topic    string `json:"topic"`
		QOS      byte   `json:"qos"`
		Retained bool   `json:"retained"`
		Message  string `json:"message"`
	}{
		msg.Topic(),
		msg.Qos(),
		msg.Retained(),
		string(msg.Payload()),
	}

	if err := json.NewEncoder(os.Stdout).Encode(jsonMessage); err != nil {
		logrus.WithError(err).Fatal("marshaling message into jsonl format")
	}
}

func subscribeCallbackLog(_ mqtt.Client, msg mqtt.Message) {
	logrus.WithFields(logrus.Fields{
		"topic":    msg.Topic(),
		"qos":      msg.Qos(),
		"retained": msg.Retained(),
		"message":  string(msg.Payload()),
	}).Info("message received")
}
