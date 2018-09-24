/*
 *  Copyright (c) 2017-2018 Samsung Electronics Co., Ltd All Rights Reserved
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
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/SamsungSLAV/slav/logger"
	"github.com/SamsungSLAV/weles"
)

const (
	stmCommand = "/usr/local/bin/stm"
)

type sshClient struct {
	config *ssh.ClientConfig
	client *ssh.Client
}

// sessionProvider implements SessionProvider interface.
// FIXME: When the connection is broken after it is established, all client functions stall.
// This provider has to be rewritten.
type sessionProvider struct {
	SessionProvider
	dryad      weles.Dryad
	connection *sshClient
	sshfs      *reverseSSHFS
	log        *os.File
}

func prepareSSHConfig(userName string, key rsa.PrivateKey) *ssh.ClientConfig {
	signer, err := ssh.NewSignerFromKey(&key)
	if err != nil {
		logger.WithError(err).Error("Failed to create signer from received ssh key.")
		// TODO: If there is a problem with parsing ssh key, job should fail.
	}

	return &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// TODO: Below will accept any host key. This should change in the future.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint:gosec
		Timeout:         30 * time.Second,            // TODO: Use value from config.
	}
}

func (d *sessionProvider) connect() (err error) {
	d.connection.client, err = ssh.Dial("tcp", d.dryad.Addr.String(), d.connection.config)
	if err != nil {
		logger.WithError(err).WithProperty("addr", d.dryad.Addr.String()).
			Error("Failed to dial dryad.")
		return err
	}
	session, err := d.connection.client.NewSession()
	if err != nil {
		logger.WithError(err).WithProperty("addr", d.dryad.Addr.String()).
			Error("Failed to create connection with dryad.")
		return err
	}
	return d.sshfs.open(session)
}

func (d *sessionProvider) newSession() (*ssh.Session, error) {
	if d.connection.client == nil {
		err := d.connect()
		if err != nil {
			logger.WithError(err).Error("Failed to connect to dryad.")
			return nil, err
		}
	}

	session, err := d.connection.client.NewSession()
	if err != nil {
		logger.WithError(err).Error("Failed to create session with dryad.")
		return nil, err
	}

	return session, nil
}

func (d *sessionProvider) executeRemoteCommand(cmd string) ([]byte, []byte, error) {
	session, err := d.newSession()
	if err != nil {
		logger.WithError(err).Error("Failed to create session with dryad.")
		return nil, nil, err
	}
	defer func() {
		if err = session.Close(); err != nil {
			logger.WithError(err).Error("Failed to close session with dryad.")
		}
	}()

	var stdout, stderr bytes.Buffer
	session.Stdout = io.MultiWriter(&stdout, os.Stderr)
	session.Stderr = io.MultiWriter(&stderr, os.Stderr)

	err = session.Run(cmd)
	return stdout.Bytes(), stderr.Bytes(), err
}

// NewSessionProvider returns new instance of SessionProvider.
func NewSessionProvider(dryad weles.Dryad, workdir string) SessionProvider {
	cfg := prepareSSHConfig(dryad.Username, dryad.Key)

	return &sessionProvider{
		dryad: dryad,
		connection: &sshClient{
			config: cfg,
		},
		sshfs: newReverseSSHFS(context.Background(), workdir, workdir),
	}
}

// Exec is a part of SessionProvider interface.
// cmd parameter is used as is. Quotations should be added by the user as needed.
func (d *sessionProvider) Exec(cmd ...string) ([]byte, []byte, error) {
	session, err := d.newSession()
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err = session.Close(); err != nil {
			logger.WithError(err).Error("Failed to close session with dryad.")
		}
	}()

	err = d.sshfs.check(session)
	if err != nil {
		logger.WithError(err).Error("Filesystem not avaliable.")
		return nil, nil, err
	}

	return d.executeRemoteCommand(strings.Join(cmd, " "))
}

// DUT is a part of SessionProvider interface.
// This function requires 'stm' binary on MuxPi's NanoPi.
func (d *sessionProvider) DUT() error {
	_, stderr, err := d.executeRemoteCommand(stmCommand + " -dut")
	if err != nil {
		logger.WithError(err).WithProperty("stderr", stderr).Error("DUT command failed.")
		return fmt.Errorf("DUT command failed: %s: %s", err, stderr)
	}
	return nil
}

// TS is a part of SessionProvider interface.
// This function requires 'stm' binary on MuxPi's NanoPi.
func (d *sessionProvider) TS() error {
	_, stderr, err := d.executeRemoteCommand(stmCommand + " -ts")
	if err != nil {
		logger.WithError(err).WithProperty("stderr", stderr).Error("TS command failed.")
		return fmt.Errorf("TS command failed: %s: %s", err, stderr)
	}
	return nil
}

// PowerTick is a part of SessionProvider interface.
// This function requires 'stm' binary on MuxPi's NanoPi.
func (d *sessionProvider) PowerTick() error {
	_, stderr, err := d.executeRemoteCommand(stmCommand + " -tick")
	if err != nil {
		logger.WithError(err).WithProperty("stderr", stderr).Error("PowerTick command failed.")
		return fmt.Errorf("PowerTick command failed: %s: %s", err, stderr)
	}
	return nil
}

// Close is a part of SessionProvider interface.
func (d *sessionProvider) Close() error {
	if d.connection.client == nil {
		return nil
	}

	if err := d.sshfs.close(); err != nil {
		logger.WithError(err).Error("Failed to close SSHFS connection with dryad.")
	}
	err := d.connection.client.Close()
	d.connection.client = nil
	return err
}
