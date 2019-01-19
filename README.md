# docker-compose-executor
New repo (for an SMA-related reference implementation) that will soon be looking for a (more) permanent home!

- The docker-compose-executor starts up, as in this logging
>level=INFO ts=2019-01-19T22:46:19.307063Z app=docker-compose-executor source=main.go:52 msg="Application dependencies resolved..."
level=INFO ts=2019-01-19T22:46:19.307101Z app=docker-compose-executor source=main.go:53 msg="Starting the docker-compose-executor application..."
level=INFO ts=2019-01-19T22:46:19.307109Z app=docker-compose-executor source=main.go:56 msg="This is the docker-compose-executor application!"
level=INFO ts=2019-01-19T22:46:19.307121Z app=docker-compose-executor source=main.go:63 msg="Application started in: 500.908µs"

- Investigating how the *docker-compose-executor* can accept incoming command (from the SMA) containing the information about which service to target, and which operation (start, stop, restart) to apply, as in this incoming command from the SMA:  `syscall.Exec("main", cmdArgs, env)`
