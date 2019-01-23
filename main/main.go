package main

import (
	"encoding/json"
	"fmt"
	"github.com/edgexfoundry/docker-compose-executor/executor"
	"github.com/edgexfoundry/docker-compose-executor/interfaces"
	"github.com/edgexfoundry/docker-compose-executor/startup"
	"github.com/edgexfoundry/edgex-go/pkg/clients/logging"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// Global variables
var Configuration *interfaces.ConfigurationStruct
var Conf = &interfaces.ConfigurationStruct{}
var err error

var usageStr = `
Usage: %s [options]
Server Options:
    -s, --service                    Indicates the service that should be targeted 
    -o, --operation                  Indicate the operation (start, stop, restart) to be performed
Common Options:
    -h, --help                       Show this message
`

//var executorClient interface{}

const (
	START   = "start"
	STOP    = "stop"
	RESTART = "restart"
)

var executorClient interface{}

func main() {

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
					//executor.LoggingClient.Error(msg)
					fmt.Println("error starting service: ", service)
				}
			} else {
				//executor.LoggingClient.Info(fmt.Sprintf("starting service {%s} succeeded", service))
				fmt.Println("success in starting service: ", service)
			}
			break

		case STOP:
			if stopper, ok := executorClient.(interfaces.ServiceStopper); ok {
				err := stopper.Stop(service)
				if err != nil {
					//msg := fmt.Sprintf("error stopping service \"%s\": %v", service, err)
					//executor.LoggingClient.Error(msg)
					fmt.Println("error stopping service: ", service)
				}
			} else {
				//executor.LoggingClient.Info(fmt.Sprintf("stopping service {%s} succeeded", service))
				fmt.Println("success in stopping service: ", service)
			}
			break

		case RESTART:
			if restarter, ok := executorClient.(interfaces.ServiceRestarter); ok {
				err := restarter.Restart(service)
				if err != nil {
					//msg := fmt.Sprintf("error restarting service \"%s\": %v", service, err)
					//executor.LoggingClient.Error(msg)
					fmt.Println("error restarting service: ", service)
				}
			} else {
				//executor.LoggingClient.Info(fmt.Sprintf("restarting service {%s} succeeded", service))
				fmt.Println("success in restarting service: ", service)
			}
			break

		default:
			executor.LoggingClient.Info("unknown operation was requested: ", operation)
			break

		}
	}
}

// TODO: clean up main() version below.

/*
func main() {

	start := time.Now()

	flag.Usage = HelpCallback
	flag.Parse()

	var useProfile string
	params := startup.BootParams{UseProfile: useProfile, BootTimeout: interfaces.BootTimeoutDefault}
	startup.Bootstrap(params, Retry, logBeforeInit)

	ok := Init()
	if !ok {
		logBeforeInit(fmt.Errorf("%s: service bootstrap failed", "docker-compose-executor"))
		os.Exit(1)
	}

	executor.LoggingClient.Info("Application dependencies resolved...")
	executor.LoggingClient.Info(fmt.Sprintf("Starting the %s application...", "docker-compose-executor"))

	http.TimeoutHandler(nil, time.Millisecond*time.Duration(Configuration.ServiceTimeout), "Request timed out")
	executor.LoggingClient.Info(Configuration.AppOpenMsg)

	errs := make(chan error, 2)
	listenForInterrupt(errs)
	startHttpServer(errs, Configuration.ServicePort)

	// Time it took to start service
	executor.LoggingClient.Info("Application started in: " + time.Since(start).String())
	executor.LoggingClient.Info("Listening on port: " + strconv.Itoa(Configuration.ServicePort))
	c := <-errs
	executor.LoggingClient.Warn(fmt.Sprintf("terminating: %v", c))

	argsWithProg := os.Args
	argsWithoutProg := os.Args[1:]
	arg := os.Args[3]

	fmt.Println(argsWithProg)
	fmt.Println(argsWithoutProg)
	fmt.Println(arg)

	os.Exit(0)
}
*/
func Retry(useConsul bool, useProfile string, timeout int, wait *sync.WaitGroup, ch chan error) {
	until := time.Now().Add(time.Millisecond * time.Duration(timeout))
	for time.Now().Before(until) {
		var err error
		// When looping, only handle configuration if it hasn't already been set.
		if Configuration == nil {
			Configuration, err = initializeConfiguration()
			if err != nil {
				ch <- err
				if !useConsul {
					//Error occurred when attempting to read from local filesystem. Fail fast.
					close(ch)
					wait.Done()
					return
				}
			} else {
				// Setup Logging
				logTarget := setLoggingTarget()
				executor.BuildLoggingClient(Configuration, logTarget)
			}
		}

		// Exit the loop if the dependencies have been satisfied.
		if Configuration != nil {
			//executorClient, err = newExecutorClient()
			break
		}
		time.Sleep(time.Second * time.Duration(1))
	}
	close(ch)
	wait.Done()

	return
}

func startHttpServer(errChan chan error, port int) {
	go func() {
		r := LoadRestRoutes()
		errChan <- http.ListenAndServe(":"+strconv.Itoa(port), r)
	}()
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}

// usage will print out the flag options for the server.
func HelpCallback() {
	msg := fmt.Sprintf(usageStr, os.Args[0])
	fmt.Printf("%s\n", msg)
	os.Exit(0)
}

func logBeforeInit(err error) {
	l := logger.NewClient("docker-compose-executor", false, "", logger.InfoLog)
	l.Error(err.Error())
}

func Init() bool {
	if Configuration == nil {
		return false
	}
	return true
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

func LoadRestRoutes() *mux.Router {
	r := mux.NewRouter()
	b := r.PathPrefix("/api/v1").Subrouter()

	//b.HandleFunc("/operation", operationHandler).Methods(http.MethodPost)
	//b.HandleFunc("/config/{services}", configHandler).Methods(http.MethodGet)
	//b.HandleFunc("/metrics/{services}", metricsHandler).Methods(http.MethodGet)

	// Ping Resource
	// /api/v1/ping
	b.HandleFunc("/ping", pingHandler).Methods(http.MethodGet)

	return r
}

// Helper function for encoding things for returning from REST calls
func encode(i interface{}, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")

	enc := json.NewEncoder(w)
	err := enc.Encode(i)

	if err != nil {
		executor.LoggingClient.Error("error during encoding", "error message", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

// Test if the service is working
func pingHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("pong"))
}

func ProcessResponse(response string) map[string]interface{} {
	rsp := make(map[string]interface{})
	err := json.Unmarshal([]byte(response), &rsp)
	if err != nil {
		executor.LoggingClient.Error("error unmarshalling response from JSON", "error message", err.Error())
	}
	return rsp
}
