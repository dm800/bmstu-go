package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"time"
)

// create human-readable SSH-key strings
func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal()) // e.g. "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTY...."
}

func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {

	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			//log.Printf("WARNING: SSH-key verification is *NOT* in effect: to fix, add this trustedKey: %q", keyString(k))
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}

		return nil
	}
}

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	conn, err := ssh.Dial("tcp", hostname+":9000", config)
	if err != nil {
		fmt.Println("TROUBLE HAPPENED", err)
	}
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	return hostname + ": " + stdoutBuf.String()
}

var tr string

func main() {

	tr = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCjPhGGeTcWblsENotLTbjjxvtI9YjpwLoXkuTTddpCIymIzjtWSfLaTsiwrJ2yaZWz6Rm5+pWUtQBN7A70ZlCh5LR3y5+vOA4xm63IHhBkTGB1TUH+91VAdShkZvriftRcFwSH+LHnTBJeBElUDx3eFB4rwBA+f97kSd6Q1yD6sbDM2s1fhH+XDcJERpn9FEvsfrv4sMGS5ru01+VzbxR0dqbIsZIpJP1XAX8Ubkp+n1ygs8566A/wg+tQXnuDuJrYTOxx+4kK+Wk7qyGpiDEMOnyeiG2i7jloN4mOmys40sbhkYlupAYZYRQcx8KG6/BqonaMKX0A6ui+9JjtWlD/"
	tr = ""
	cmd := os.Args[1]
	hosts := os.Args[2:]
	results := make(chan string, 10)
	timeout := time.After(5 * time.Second)

	pkey, err := os.ReadFile("/home/dm900/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(pkey)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	config := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: trustedHostKeyCallback(tr),
	}

	for _, hostname := range hosts {
		fmt.Println(hostname)
		go func(hostname string) {
			results <- executeCmd(cmd, hostname, config)
		}(hostname)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
}
