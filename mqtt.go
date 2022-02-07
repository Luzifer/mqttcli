package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func mqttTokenToError(tok mqtt.Token) error {
	if !tok.WaitTimeout(cfg.MQTTTimeout) {
		return errors.New("token wait timed out")
	}

	return tok.Error()
}

func publish(client mqtt.Client) error {
	for _, t := range cfg.Topics {
		logger := log.WithField("topic", t)

		if err := mqttTokenToError(client.Publish(t, byte(cfg.QOS), cfg.Retain, cfg.Message)); err != nil {
			logger.WithError(err).Fatal("Unable to publish message")
		}

		logger.Info("Message published")
	}

	return nil
}

func subscribe(client mqtt.Client) error {
	var (
		callback mqtt.MessageHandler
		topics   = map[string]byte{}
	)

	for _, t := range cfg.Topics {
		topics[t] = byte(cfg.QOS)
	}

	switch cfg.OutputFormat {
	case "log":
		callback = subscribeCallbackLog

	case "csv":
		fmt.Println("Topic,QOS,Retained,Message")
		callback = subscribeCallbackCSV

	case "jsonl":
		callback = subscribeCallbackJSONL

	default:
		log.WithField("format", cfg.OutputFormat).Fatal("Invalid output format specified")
	}

	if err := mqttTokenToError(client.SubscribeMultiple(topics, callback)); err != nil {
		log.WithError(err).Fatal("Unable to subscribe topics")
	}

	for {
		select {}
	}
}

func subscribeCallbackLog(client mqtt.Client, msg mqtt.Message) {
	log.WithFields(log.Fields{
		"topic":    msg.Topic(),
		"qos":      msg.Qos(),
		"retained": msg.Retained(),
		"message":  string(msg.Payload()),
	}).Info("Message received")
}

func subscribeCallbackCSV(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("%s,%d,%v,%q\n",
		msg.Topic(),
		msg.Qos(),
		msg.Retained(),
		string(msg.Payload()),
	)
}

func subscribeCallbackJSONL(client mqtt.Client, msg mqtt.Message) {
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
		log.WithError(err).Fatal("Unable to marshal message into jsonl")
	}
}
