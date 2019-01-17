package interfaces

// TODO: The abstraction which should be accessed via a global var.
// For now, we have this as a duplicate of the (same) abstraction as exists in the edgex-go repo;
// TODO: Eventually, turn this into a shared lib, perhaps...

type ExecutorClient interface {
	ServiceStarter(service string) error
	ServiceStopper(service string) error
	ServiceRestarter(service string) error
}
