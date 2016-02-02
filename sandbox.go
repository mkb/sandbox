package sandbox

import (
	"fmt"
	"bytes"
	docker "github.com/fsouza/go-dockerclient"
)

// Options is used to specify configuration values for the Run
// function.
type Options struct {
	// Endpoint specifies the URL at which to contact Docker daemon
	// for running containers.
	// Defaults to `unix:///var/run/docker.sock`
	Endpoint string

	// Image identifies the name (or URL) of the Docker image to
	// instantiate as a sandbox container.
	// Defaults to `ubuntu`
	Image    string

	// Command is the set of commands that will be fed to
	// `bash -lc`.  Multiple commands can be separated by newlines,
	// or other Bash command separators, like ';' and '&&'.
	// It is an error to omit this option.
	Command  string

	// How long to allow the command to execute before forcibly
	// timing it out (in seconds).
	// Defaults to 30.
	Timeout  int
}

// Run spins up a new sandbox container, issues the command(s), and
// then waits for the exit code and output (both standard output
// and standard error, intermixed).
func Run(opt Options) string {
	if opt.Endpoint == "" {
		opt.Endpoint = "unix:///var/run/docker.sock"
	}
	if opt.Image == "" {
		opt.Image = "ubuntu"
	}
	if opt.Timeout <= 0 {
		opt.Timeout = 30
	}
	if opt.Command == "" {
		return ">> error: no command specified\n"
	}

	client, err := docker.NewClient(opt.Endpoint)
	if err != nil {
		return fmt.Sprintf(">> error: %s\n", err)
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: opt.Image,
			Cmd:   []string{"timeout", "--preserve-status", fmt.Sprintf("%ds", opt.Timeout), "/bin/bash", "-l", "-c", opt.Command},
			Tty:   true,
		},
	})
	if err != nil {
		return fmt.Sprintf(">> error: %s\n", err)
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
		return fmt.Sprintf(">> error: %s\n", err)
	}

	_, err := client.WaitContainer(c.ID)
	if err != nil {
		return fmt.Sprintf(">> error: %s\n", err)
	}

	client.RemoveContainer(docker.RemoveContainerOptions{ID: c.ID})
	return <-output
}
