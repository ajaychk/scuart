package lcn

import (
	"log"
	"sync"

	"github.com/eclipse/paho.mqtt.golang"
)

var wg sync.WaitGroup
var cl mqtt.Client

func init() {
	go initClient()
}

func initClient() {
	opt := mqtt.NewClientOptions()
	opt = opt.AddBroker("tcp://localhost:1883")
	opt = opt.SetUsername("senra")
	opt = opt.SetPassword("sc2havellssenra")

	cl = mqtt.NewClient(opt)

	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	wg.Add(1)

	if token := cl.Subscribe("/lcn/payloads/uplink", 0, ulHandler); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	wg.Wait()
}

// handle received uplink
func ulHandler(cl mqtt.Client, msg mqtt.Message) {
	log.Printf("message received %s ", msg.Payload())
	wg.Done()
}

//  send downlink to node using lcn
func Send(id string, data []byte) (err error) {
	log.Printf("SENDING DATA % x to %s DEVICE.\n", data, id)

	// pl := &struct {
	// 	id   string
	// 	data []byte
	// }{id, data}

	if token := cl.Publish("/lcn/payloads/downlink", 0, true, data); token.Wait() && token.Error() != nil {
		err = token.Error()
	}
	return
}
