package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// TODO: Externalize these by putting in a properties file.
const (
	ConfigSeedServiceKey            = "edgex-config-seed"
	CoreCommandServiceKey           = "edgex-core-command"
	CoreDataServiceKey              = "edgex-core-data"
	CoreMetaDataServiceKey          = "edgex-core-metadata"
	ExportClientServiceKey          = "edgex-export-client"
	ExportDistroServiceKey          = "edgex-export-distro"
	SupportLoggingServiceKey        = "edgex-support-logging"
	SupportNotificationsServiceKey  = "edgex-support-notifications"
	SystemManagementAgentServiceKey = "edgex-sys-mgmt-agent"
	SupportSchedulerServiceKey      = "edgex-support-scheduler"
)

// TODO: Externalize these by putting in a properties file.
var services = map[string]string{
	SupportNotificationsServiceKey: "Notifications",
	CoreCommandServiceKey:          "Command",
	CoreDataServiceKey:             "CoreData",
	CoreMetaDataServiceKey:         "Metadata",
	ExportClientServiceKey:         "Export",
	ExportDistroServiceKey:         "Distro",
	SupportLoggingServiceKey:       "Logging",
	SupportSchedulerServiceKey:     "Scheduler",
}

type ExecuteDockerCompose struct {
}

func (ec *ExecuteDockerCompose) Start(service string) error {
	error := ExecuteDockerCommands(service, "start")
	return error
}

func (ec *ExecuteDockerCompose) Stop(service string) error {
	error := ExecuteDockerCommands(service, "stop")
	return error
}

func (ec *ExecuteDockerCompose) Restart(service string) error {
	error := ExecuteDockerCommands(service, "restart")
	return error
}

func findDockerContainerStatus(service string, status string) bool {

	var (
		cmdOut []byte
		err    error
	)
	cmdName := "docker"
	cmdArgs := []string{"ps"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Println("error running the docker command", "error message", err.Error())
		os.Exit(1)
	}
	dockerOutput := string(cmdOut)
	// Find whether the container to start has started.
	for _, line := range strings.Split(strings.TrimSuffix(dockerOutput, "\n"), "\n") {
		if strings.Contains(line, service) {

			if status == "Up" {
				if strings.Contains(line, "Up") {
					fmt.Println("container started", "service name", service, "details", line)
					return true
				} else {
					fmt.Println("container NOT started", "service name", service)
					return false
				}
			} else if status == "Exited" {
				if strings.Contains(line, "Exited") {
					fmt.Println("container stopped", "service name", service, "details", line)
					return true
				} else {
					fmt.Println("container NOT stopped", "service name", service)
					return false
				}
			}
		}
	}
	return false
}

func ExecuteDockerCommands(service string, operation string) error {
	_, knownService := services[service]

	if knownService {
		runDockerCommands(service, services[service], operation)

		return nil
	} else {
		newError := fmt.Errorf("unknown service: %v", service)
		fmt.Println(newError.Error())

		return newError
	}
}

func runDockerCommands(service string, dockerService string, operation string) {

	var (
		err    error
		cmdDir string
	)

	cmdName := "docker"

	// Retry finding the path to where the docker command will be run.
	err = Do(func(attempt int) (bool, error) {
		var err error
		cmdDir, err = findPathToRunDocker()
		// Try 5 times
		return attempt < 5, err
	})
	if err != nil {
		fmt.Println("unable to find the path to where the docker command will be run", "error message", err.Error())
	}

	cmdArgs := []string{operation, dockerService}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = cmdDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("docker command failed", "error message", err.Error())
		fmt.Println("associated ouptut", "error message", string(out))
	}

	if operation == "start" {
		if !findDockerContainerStatus(service, "Up") {
			fmt.Println("docker start operation failed", "service name", service)
		}
	} else if operation == "stop" {
		if !findDockerContainerStatus(service, "Exited") {
			fmt.Println("docker stop operation failed", "service name", service)
		}
	}
}

func findPathToRunDocker() (string, error) {

	// Determine the directory (in the deployed filesystem) from where docker will be executed.
	cmdName := "pwd"

	cmd := exec.Command(cmdName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("exec.Command(cmdName) failed", "error message", err.Error())
	}
	pathOutput := string(out)

	path := strings.TrimSuffix(pathOutput, "\n")

	return path, err
}
