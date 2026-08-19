package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

func publish(client mqtt.Client) (err error) {
	for _, t := range cfg.Topics {
		logger := logrus.WithField("topic", t)

		//#nosec:G115 // QoS is expected to be 0,1,2 - fine to convert
		if err = mqttTokenToError(client.Publish(t, byte(cfg.QOS), cfg.Retain, cfg.Message)); err != nil {
			return fmt.Errorf("publishing message to %q: %w", t, err)
		}

		logger.Info("message published")
	}

	return nil
}
