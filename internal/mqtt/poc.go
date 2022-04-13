package mqtt

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type sarasa struct {
	StationID string `json:"station_id"`
	OutputID  string `json:"output_id"`
	Command   string `json:"command"`
}

func main() {
	fmt.Println("Hello, World!")
	o := MQTT.NewClientOptions()
	o.AddBroker("tcp://192.168.1.100:1883")
	o.SetClientID("iot-dispatcher")
	o.SetUsername("dispatcher")
	o.SetPingTimeout(1 * time.Second)

	messageJSON, err := json.Marshal(sarasa{
		StationID: "1",
		OutputID:  "2",
		Command:   "3",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := MQTT.NewClient(o)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := client.Publish("commands/STA01010", 0, false, messageJSON)
	token.Wait()
	client.Disconnect(250)
}
