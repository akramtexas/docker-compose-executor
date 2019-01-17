/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package main

import (
	"github.com/edgexfoundry/docker-compose-executor/interfaces"
	"github.com/edgexfoundry/docker-compose-executor/startup"
)

// Global variables
var Configuration *interfaces.ConfigurationStruct
var Conf = &interfaces.ConfigurationStruct{}
var err error

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
