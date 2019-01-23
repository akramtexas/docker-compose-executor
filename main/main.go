package main

import (
	"fmt"
	"github.com/edgexfoundry/docker-compose-executor/interfaces"
	"os"
	"time"
)

// Global variables
var executorClient interface{}
var err error

const (
	START           = "start"
	STOP            = "stop"
	RESTART         = "restart"
	ApplicationName = "docker-compose-executor"
	AppOpenMsg      = "This is the docker-compose-executor application!"
)

func main() {

	start := time.Now()

	fmt.Println("Application dependencies resolved...")
	fmt.Println(fmt.Sprintf("Starting the %s application...", ApplicationName))
	fmt.Println(AppOpenMsg)

	// Time it took to start service
	fmt.Println("Application started in: " + time.Since(start).String())

	var service = ""
	var operation = ""

	if len(os.Args) > 2 {
		service = os.Args[1]
		operation = os.Args[2]

		switch operation {
		case START:
			if starter, ok := executorClient.(interfaces.ServiceStarter); ok {
				err := starter.Start(service)
				if err != nil {
					//msg := fmt.Sprintf("error starting service \"%s\": %v", service, err)
					//fmt.Println.Error(msg)
					fmt.Println("error starting service: ", service)
				}
			} else {
				//fmt.Println.Info(fmt.Sprintf("starting service {%s} succeeded", service))
				fmt.Println("success in starting service: ", service)
			}
			break

		case STOP:
			if stopper, ok := executorClient.(interfaces.ServiceStopper); ok {
				err := stopper.Stop(service)
				if err != nil {
					//msg := fmt.Sprintf("error stopping service \"%s\": %v", service, err)
					//fmt.Println.Error(msg)
					fmt.Println("error stopping service: ", service)
				}
			} else {
				//fmt.Println.Info(fmt.Sprintf("stopping service {%s} succeeded", service))
				fmt.Println("success in stopping service: ", service)
			}
			break

		case RESTART:
			if restarter, ok := executorClient.(interfaces.ServiceRestarter); ok {
				err := restarter.Restart(service)
				if err != nil {
					//msg := fmt.Sprintf("error restarting service \"%s\": %v", service, err)
					//fmt.Println.Error(msg)
					fmt.Println("error restarting service: ", service)
				}
			} else {
				//fmt.Println.Info(fmt.Sprintf("restarting service {%s} succeeded", service))
				fmt.Println("success in restarting service: ", service)
			}
			break

		default:
			fmt.Println("unknown operation was requested: ", operation)
			break

		}
	}
}
