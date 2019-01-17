package main

import (
	"github.com/edgexfoundry/docker-compose-executor/executor"
	"github.com/edgexfoundry/docker-compose-executor/interfaces"
	"github.com/edgexfoundry/docker-compose-executor/startup"
	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
)

func main() {

	var useProfile string

	params := startup.BootParams{UseProfile: useProfile, BootTimeout: interfaces.BootTimeoutDefault}

	startup.Bootstrap(params, Retry, logBeforeInit)

	// TODO: Setup Logging (Refine this admittedly kludgy setup..)
	logTarget := setLoggingTarget()
	executor.BuildLoggingClient(Configuration, logTarget)
}

func logBeforeInit(err error) {
	l := logger.NewClient(executor.SystemManagementAgentServiceKey, false, "", logger.InfoLog)
	l.Error(err.Error())
}
