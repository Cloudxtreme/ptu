package config

import (
	"errors"
	"flag"
	"os"
	"regexp"
)

// Config is a container for ptu configuration
type Config struct {
	SSHServer   string
	SSHUsername string
	SSHPassword string
	SSHUseAgent bool
	TargetHost  string
	ExposedBind string
	ExposedPort int
	ExposedHost string
	ConnectTo   string
}

// IsListEmpty checks if no command line arguments were passed
func IsListEmpty() bool {
	return len(os.Args) < 2
}

// IsHelpRequested checks if help was requested (by passing -h|--help as an argument)
func IsHelpRequested() bool {
	var helpArgumentRegexp = regexp.MustCompile(`^(-h|--help)$`)
	return helpArgumentRegexp.MatchString(os.Args[1])
}

// ParseArguments parses command line arguments, performs some initial validation and variable mutation
func ParseArguments(d *Config) (*Config, error) {
	var sshServer = flag.String("s", d.SSHServer, "SSH server (host[:port]) to connect")
	var sshUsername = flag.String("u", d.SSHUsername, "username to connect SSH server")
	var sshPassword = flag.String("p", d.SSHPassword, "password to authenticate against SSH server (do not use, please)")
	var targetHost = flag.String("t", d.TargetHost, "target host:port we will forward connections to")
	var exposedBind = flag.String("b", d.ExposedBind, "bind (listener) to expose on the SSH server side")
	var exposedPort = flag.Int("e", d.ExposedPort, "port to expose and forward on the SSH server side")

	flag.Parse()

	if !isSSHServerSet(*sshServer) {
		return nil, errors.New("SSH server not defined")
	}

	if !isTCPPortValid(*exposedPort) {
		return nil, errors.New("exposed TCP port number is invalid")
	}

	if !isHostWithPort(*sshServer) {
		*sshServer = joinHostPort(*sshServer, defaultSSHPort)
	}

	if !isHostWithPort(*targetHost) {
		*targetHost = joinHostPort(*targetHost, defaultTargetPort)
	}

	config := &Config{
		SSHServer:   *sshServer,
		SSHUsername: *sshUsername,
		SSHPassword: *sshPassword,
		SSHUseAgent: !isSSHPasswordSet(*sshPassword),
		TargetHost:  *targetHost,
		ExposedBind: *exposedBind,
		ExposedPort: *exposedPort,
		ExposedHost: joinHostPort(*exposedBind, *exposedPort),
		ConnectTo:   mergeHostPort(*sshServer, *exposedPort),
	}

	return config, nil
}
