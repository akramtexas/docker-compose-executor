package main

import (
	"github.com/edgexfoundry/docker-compose-executor/executor"
	"github.com/edgexfoundry/docker-compose-executor/interfaces"
	"github.com/edgexfoundry/docker-compose-executor/startup"
	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
)

// Global variables
var Configuration *interfaces.ConfigurationStruct
var Conf = &interfaces.ConfigurationStruct{}
var err error

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

func Retry() {

	if Configuration == nil {
		Configuration, err = initializeConfiguration()
	}
	return
}

func initializeConfiguration() (*interfaces.ConfigurationStruct, error) {
	//We currently have to load configuration from filesystem first in order to obtain ConsulHost/Port
	err := startup.LoadFromFile(Conf)
	if err != nil {
		return nil, err
	}

	return Conf, nil
}

func setLoggingTarget() string {
	logTarget := Configuration.LoggingRemoteURL
	if !Configuration.EnableRemoteLogging {
		return Configuration.LoggingFile
	}
	return logTarget
}
