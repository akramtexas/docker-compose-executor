package main

import (
	"github.com/edgexfoundry/docker-compose-executor/executor"
)

func main() {
	// TODO: Setup Logging (Refine this admittedly kludgy setup..)
	logTarget := setLoggingTarget()
	executor.BuildLoggingClient(Configuration, logTarget)
}
