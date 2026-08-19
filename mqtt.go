package main

import (
	"errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mqttTokenToError(tok mqtt.Token) error {
	if !tok.WaitTimeout(cfg.MQTTTimeout) {
		return errors.New("token wait timed out")
	}

	return tok.Error() //nolint:wrapcheck // fine in this case, only used in logging
}
