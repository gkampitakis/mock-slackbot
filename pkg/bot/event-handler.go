package bot

import "log"

func eventHandler(event interface{}) {
	log.Println(event)
}
