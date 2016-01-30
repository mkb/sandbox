package main

import (
	"bytes"
	docker "github.com/fsouza/go-dockerclient"
)

func Run(command string) (int, string, error) {
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return -1, "", err
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: "ubuntu",
			Cmd:   []string{"/bin/bash", "-l", "-c", command},
			Tty:   true,
		},
	})
	if err != nil {
		return -1, "", err
	}

	output := make(chan string, 0)
	go func() {
		var b bytes.Buffer
		client.AttachToContainer(docker.AttachToContainerOptions{
			Container: c.ID,
			Stream:    true,

			Stdout:       true,
			OutputStream: &b,

			Stderr:      true,
			ErrorStream: &b,

			RawTerminal: true,
		})

		output <- b.String()
	}()

	err = client.StartContainer(c.ID, nil)
	if err != nil {
		return -1, "", err
	}

	rc, err := client.WaitContainer(c.ID)
	if err != nil {
		return -1, "", err
	}

	client.RemoveContainer(docker.RemoveContainerOptions{ID: c.ID})
	return rc, <-output, nil
}
