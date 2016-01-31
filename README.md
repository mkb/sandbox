# sandbox
--
    import "github.com/jhunt/sandbox"

Package sandbox provides a safe and hassle-free way of running arbitrarily
complex commands, with potentially harmful side-effects, without incurring the
risk to the host system.

## Usage

#### func  Run

```go
func Run(opt Options) (int, string, error)
```
Run spins up a new sandbox container, issues the command(s), and then waits for
the exit code and output (both standard output and standard error, intermixed).

#### type Options

```go
type Options struct {
	// Endpoint specifies the URL at which to contact Docker daemon
	// for running containers.
	// Defaults to `unix:///var/run/docker.sock`
	Endpoint string

	// Image identifies the name (or URL) of the Docker image to
	// instantiate as a sandbox container.
	// Defaults to `ubuntu`
	Image string

	// Command is the set of commands that will be fed to
	// `bash -lc`.  Multiple commands can be separated by newlines,
	// or other Bash command separators, like ';' and '&&'.
	// It is an error to omit this option.
	Command string

	// How long to allow the command to execute before forcibly
	// timing it out (in seconds).
	// Defaults to 30.
	Timeout int
}
```

Options is used to specify configuration values for the Run function.
