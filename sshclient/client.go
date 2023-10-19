package sshclient

import (
	"bytes"
	"log"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Connection struct {
	client *ssh.Client
}
type SshCreds struct {
	User    string
	Host    string
	KeyFile string
}

const sshCommandRetries = 5

func ConnectWithKey(creds SshCreds) Connection {
	// parse the user's private key:
	pvtKeyBts, err := os.ReadFile(creds.KeyFile)
	if err != nil {
		log.Fatalln(err)
	}

	signer, err := ssh.ParsePrivateKey(pvtKeyBts)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := ssh.Dial("tcp", creds.Host, &ssh.ClientConfig{
		User: creds.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// use OpenSSH's known_hosts file if you care about host validation
			return nil
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	return Connection{client: client}
}

func (c *Connection) newSession(buf *bytes.Buffer) *ssh.Session {
	session, err := c.client.NewSession()
	if err != nil {
		log.Fatalln(err)
	}
	if buf == nil {
		session.Stdout = os.Stdout
		session.Stderr = os.Stderr
	} else {
		session.Stdout = buf
		session.Stderr = buf
	}

	if err != nil {
		log.Fatal(err)
	}
	return session
}

func (c *Connection) ExecCommand(command string) {
	session := c.newSession(nil)
	defer session.Close()
	if err := session.Run(command); err != nil {
		log.Fatalln(err)
	}
}

// Executes provided command on ssh server and returns its output.

// Sometimes ssh client just returns empty execution output without error however output is expected.
// shouldRetry function provides condition in that case should we retry command
func (c *Connection) ExecCommandWithOutput(command string, shouldRetry func(string) bool) string {
	execute := func() string {
		var buf bytes.Buffer
		session := c.newSession(&buf)

		defer session.Close()
		if err := session.Run(command); err != nil {
			log.Fatalln(err)
		}
		return strings.TrimSuffix(buf.String(), "\n")
	}
	for attempt := 0; attempt < sshCommandRetries; attempt++ {
		out := execute()
		if !shouldRetry(out) {
			return out
		}
	}
	log.Fatalf("Failed to execute ssh command: %s\n", command)
	// Unreachable code. Satisfy complier.
	return ""
}
