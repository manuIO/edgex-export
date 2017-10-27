package distro

import (
	"github.com/drasko/edgex-export"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"strconv"
)

type mqttSender struct {
	mqttClient MQTT.Client
}

const clientID = "edgex"

func NewMqttSender(addr export.Addressable) Sender {
	opts := MQTT.NewClientOptions()
	// CHN: Should be added protocol from Addressable instead of include it the address param.
	// CHN: We will maintain this behaviour for compatibility with Java
	broker := addr.Address + ":" + strconv.Itoa(addr.Port)
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetUsername(addr.User)
	opts.SetPassword(addr.Password)

	var sender mqttSender

	sender.mqttClient = MQTT.NewClient(opts)
	if token := sender.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		// FIXME
		panic(token.Error())
	}
	logger.Info("Sample Publisher Started")

	return sender
}

func (sender mqttSender) Send(data string) {
	token := sender.mqttClient.Publish("FCR", 0, false, data)
	// FCR could be removed? set of tokens?
	token.Wait()
	logger.Info("Sent data: ", zap.String("data", data))
}