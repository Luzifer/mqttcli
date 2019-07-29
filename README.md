[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/mqttcli)](https://goreportcard.com/report/github.com/Luzifer/mqttcli)
![](https://badges.fyi/github/license/Luzifer/mqttcli)
![](https://badges.fyi/github/downloads/Luzifer/mqttcli)
![](https://badges.fyi/github/latest-release/Luzifer/mqttcli)

# Luzifer / mqttcli

`mqttcli` is a small CLI util to interface with a MQTT server. It can be used to publish messages to oder subscribe messages from a MQTT server.

At the moment it is intended to connect to simple setups using plain `tcp`, `ssl` with trusted certificates or websockets (`ws`). More options to come later.

## Usage

```console
# mqttcli --help
Usage of mqttcli:
      --log-level string        Log level (debug, info, warn, error, fatal) (default "info")
  -m, --message string          
  -b, --mqtt-broker string      Broker URI to connect to scheme://host:port (scheme is one of 'tcp', 'ssl', or 'ws') (default "tcp://localhost:1883")
      --mqtt-client-id string   Client ID to use when connecting, must be unique (default "21064ab7-c296-445e-b8b6-d1bced77853c")
  -p, --mqtt-pass string        Password to identify against the broker
      --mqtt-timeout duration   How long to wait for the client to complete operations (default 10s)
  -u, --mqtt-user string        Username to identify against the broker
  -o, --output-format string    How to ouptut received messages (One of 'log', 'csv', 'jsonl') (default "log")
      --qos int                 QOS to use (0 - Only Once, 1 - At Least Once, 2 - Only Once) (default 1)
      --retain                  Retain message on topic
  -t, --topic strings           Topic to subscribe / publish to
      --version                 Prints current version and exits
```

## Examples

```console
# envrun -- mqttcli sub -t 'mysensor/+'
INFO[0001] Message received                              message=4058 qos=1 retained=false topic=mysensor/co2

# envrun -- mqttcli sub -t 'mysensor/+' -o csv
Topic,QOS,Retained,Message
mysensor/co2,1,false,"3978"

# envrun -- mqttcli sub -t 'mysensor/+' -o jsonl
{"topic":"mysensor/co2","qos":1,"retained":false,"message":"3972"}
```

```console
# envrun -- mqttcli pub -t mysensor/test -m 'ohai?'

# envrun -- mqttcli sub -t 'mysensor/+'
INFO[0001] Message received                              message="ohai?" qos=1 retained=false topic=mysensor/test
```
