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
	dclient, err := dockerapi.NewClient("tcp://127.0.0.1:2376")
	assert(err)

	if containers, err := dclient.ListContainers(dockerapi.ListContainersOptions{All: true}); err != nil {
		log.Fatalf("Get container error: %v", err)
	} else {
		for _, container := range containers {
			if c, err := dclient.InspectContainer(container.ID); err != nil {
				log.Printf("%s: %v\n", container.ID, err)
			} else {
				if c.State.Running {
					continue
				} else {
					log.Println(c.ID)
				}
			}
		}
	}
}
