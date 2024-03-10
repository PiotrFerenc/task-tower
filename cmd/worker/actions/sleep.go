package actions

import (
	"log"
	"time"
)

type Sleep struct {
}

func (receiver Sleep) Execute(_ ActionContext) string {

	for i := 0; i < 5; i++ {
		time.Sleep(5 * time.Second)
		log.Printf("sleep: %d", i)
	}
	return ""
}
