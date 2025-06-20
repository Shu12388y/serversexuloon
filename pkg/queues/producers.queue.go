package pkg

import (
	"encoding/json"
	"time"
)

type ProducerJob struct {
	Payload any
	Event   string
	State   bool
}

// using queue by go routines
func Producers(message ProducerJob, ch chan string) {

	data := message

	jsonData, err := json.Marshal(data)
	if err != nil {
		// Handle error
		return
	}
	ch <- string(jsonData)
	time.Sleep(time.Millisecond * 10)

	close(ch)

}
