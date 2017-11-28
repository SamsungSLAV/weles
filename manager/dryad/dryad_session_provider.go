/*
 *  Copyright (c) 2017 Samsung Electronics Co., Ltd All Rights Reserved
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License
 */

package dryad

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crypto/rsa"

	. "git.tizen.org/tools/weles"
	"golang.org/x/crypto/ssh"
)

const (
	stmCommand = "/usr/local/stm"
)

type sshClient struct {
	config *ssh.ClientConfig
	client *ssh.Client
}

// sessionProvider implements SessionProvider interface.
// FIXME: When the connection is broken after it is established, all client functions stall. This provider has to be rewritten.
type sessionProvider struct {
	SessionProvider
	dryad      Dryad
	connection *sshClient
}

func prepareSSHConfig(userName string, key rsa.PrivateKey) *ssh.ClientConfig {
	signer, _ := ssh.NewSignerFromKey(&key)

	return &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second, // TODO: Use value from config when such appears.
	}
}

func (d *sessionProvider) connect() (err error) {
	d.connection.client, err = ssh.Dial("tcp", d.dryad.Addr.String(), d.connection.config)
	return
}

func (d *sessionProvider) newSession() (*ssh.Session, error) {
	if d.connection.client == nil {
		err := d.connect()
		if err != nil {
			return nil, err
		}
	}

	session, err := d.connection.client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (d *sessionProvider) executeRemoteCommand(cmd string) ([]byte, []byte, error) {
	session, err := d.newSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	return stdout.Bytes(), stderr.Bytes(), err
}

// NewSessionProvider returns new instance of SessionProvider.
func NewSessionProvider(dryad Dryad) SessionProvider {
	cfg := prepareSSHConfig(dryad.Username, dryad.Key)

	return &sessionProvider{
		dryad: dryad,
		connection: &sshClient{
			config: cfg,
		},
	}
}

// Exec is a part of SessionProvider interface.
// FIXME: Exec function checks every argument and if contains space (except surrounding ones) surrounds the argument
// with double quotes. Caller must be aware of such functionality because it may break some special arguments.
func (d *sessionProvider) Exec(cmd []string) (stdout, stderr []byte, err error) {
	joinedCommand := cmd[0] + " "
	for i := 1; i < len(cmd); i++ {
		if strings.Contains(strings.Trim(cmd[i], " "), " ") {
			joinedCommand += `"` + cmd[i] + `" `
		} else {
			joinedCommand += cmd[i] + " "
		}
	}
	return d.executeRemoteCommand(joinedCommand)
}

// DUT is a part of SessionProvider interface.
// This function requires 'stm' binary on MuxPi's NanoPi.
func (d *sessionProvider) DUT() error {
	_, stderr, err := d.executeRemoteCommand(stmCommand + " -dut")
	if err != nil {
		return fmt.Errorf("DUT command failed: %s : %s", err, stderr)
	}
	return nil
}

// TS is a part of SessionProvider interface.
// This function requires 'stm' binary on MuxPi's NanoPi.
func (d *sessionProvider) TS() error {
	_, stderr, err := d.executeRemoteCommand(stmCommand + " -ts")
	if err != nil {
		return fmt.Errorf("TS command failed: %s : %s", err, stderr)
	}
	return nil
}

// PowerTick is a part of SessionProvider interface.
// This function requires 'stm' binary on MuxPi's NanoPi.
func (d *sessionProvider) PowerTick() error {
	_, stderr, err := d.executeRemoteCommand(stmCommand + " -tick")
	if err != nil {
		return fmt.Errorf("PowerTick command failed: %s : %s", err, stderr)
	}
	return nil
}

// Close is a part of SessionProvider interface.
func (d *sessionProvider) Close() error {
	if d.connection.client == nil {
		return nil
	}

	err := d.connection.client.Close()
	d.connection.client = nil
	return err
}

// SendFile is a part of SessionProvider interface.
func (d *sessionProvider) SendFile(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return err
	}

	session, err := d.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	filename := filepath.Base(dst)
	directory := filepath.Dir(dst)

	w, err := session.StdinPipe()
	if err != nil {
		return err
	}
	defer w.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	// Trigger SCP sink mode
	err = session.Start("scp -t " + directory)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, "C0755", s.Size(), filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, "\x00")
	if err != nil {
		return err
	}

	err = session.Wait()

	// FIXME: unexpected <newline> is reported by scp every time the transfer is finished properly. Needs to be solved.
	// Bellow we have a very lousy trick. I hope it will be fixed in the future.
	// I don't know what is the reason or how to fix it. Has to wait a little bit. Or maybe someone else will find the solution.
	// First candidate sshfs -o slave
	if strings.Contains(stdout.String(), "unexpected <newline>") {
		return nil
	}
	return err
}

// ReceiveFile is a part of SessionProvider interface.
func (d *sessionProvider) ReceiveFile(src, dst string) error {
	session, err := d.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	r, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	err = session.Start("scp " + src + " /dev/stdout")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	return session.Wait()
}
