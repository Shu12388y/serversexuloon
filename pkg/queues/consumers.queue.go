package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)


type ConsumerJob struct {
	Payload any
	Event   string
	State   bool
}

func Consumers(ch chan string) {
	for msg := range ch {
		var data ConsumerJob
		err := json.Unmarshal([]byte(msg), &data)
		if err != nil{
			log.Fatal(err)
		}
		fmt.Printf("%v", data.Event)
		fmt.Printf("%v", data.Payload.(map[string]interface{})["data"])
		fmt.Printf("%v", data.State)
		time.Sleep(time.Millisecond * 10)
	}

}
