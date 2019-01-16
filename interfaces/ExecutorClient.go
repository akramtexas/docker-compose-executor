package interfaces

// TODO: The abstraction which should be accessed via a global var.
// Duplicate of the abstraction in the mono repo; make it a shared lib,perhaps...

type ExecutorClient interface {
	ServiceStarter(service string) error
	ServiceStopper(service string) error
	ServiceRestarter(service string) error
}
