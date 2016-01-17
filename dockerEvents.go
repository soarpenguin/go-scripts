package main

import (
	"log"

	dockerapi "github.com/fsouza/go-dockerclient"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	docker, err := dockerapi.NewClient("tcp://127.0.0.1:2375")
	assert(err)

	// Start event listener before listing containers to avoid missing anything
	events := make(chan *dockerapi.APIEvents)
	assert(docker.AddEventListener(events))
	log.Println("Listening for Docker events ...")

	// Process Docker events
	for msg := range events {
		switch msg.Status {
		case "start":
			log.Println("add container: ", msg.ID)
		case "die":
			log.Println("remove container: ", msg.ID)
		case "stop", "kill":
			log.Println("stop container: ", msg.ID)
		}
	}
}
