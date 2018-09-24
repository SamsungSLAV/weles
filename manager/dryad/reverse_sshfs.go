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
	"fmt"
	"os/exec"

	"golang.org/x/crypto/ssh"

	"github.com/SamsungSLAV/slav/logger"
)

// reverseSSHFS will start a process on the remote and the local hosts and a goroutine.
// It is not safe for concurrent use.
type reverseSSHFS struct {
	ctx        context.Context
	sftpCancel context.CancelFunc
	session    *ssh.Session
	// TODO (amistewicz): replace with a single buffer wrapped in logger
	sshfsStderr bytes.Buffer
	sftpStderr  bytes.Buffer
	// results of sftp and sshfs commands. Always non-nil values are passed.
	errors        chan error
	local, remote string
}

func newReverseSSHFS(ctx context.Context, local, remote string) *reverseSSHFS {
	return &reverseSSHFS{
		ctx:    ctx,
		errors: make(chan error),
		local:  local,
		remote: remote,
	}
}

// open mounts local directory on remote.
//
// It connects stdin of locally running sftp-server to stdout of remotely running sshfs,
// same with stdout of sftp-server and stdin of sshfs. It starts two goroutines.
//
// Provided session must not be used for Run, Start or Shell calls.
// If an error is returned the caller is responsible for closing session. open may be called
// again with new instance of ssh.Session if the previous call failed.
//
// To ensure that the filesystem is mounted one should use check().
func (sshfs *reverseSSHFS) open(session *ssh.Session) (err error) {
	// remote:  sshfs -o slave :path path
	//                ^use stdin and stdout
	// local:   /usr/lib/openssh/sftp-server -e
	//                                       ^log on stderr
	ctx, cancel := context.WithCancel(sshfs.ctx)

	// gas/gosec returns error here about subprocess launching with variable. This is intended.
	sftp := exec.CommandContext(ctx, "/usr/lib/openssh/sftp-server", "-e", "-l", "INFO") //nolint: gas, gosec,lll
	session.Stdin, err = sftp.StdoutPipe()
	if err != nil {
		logger.WithError(err).Error("Failed to get sftp stdout pipe.")
		cancel()
		return
	}
	session.Stdout, err = sftp.StdinPipe()
	if err != nil {
		logger.WithError(err).Error("Failed to get sftp stdin pipe.")
		cancel()
		return
	}

	// Collect stderr
	session.Stderr = &sshfs.sshfsStderr
	sftp.Stderr = &sshfs.sftpStderr

	err = sftp.Start()
	if err != nil {
		logger.WithError(err).Error("Failed to start reverse SSHFS.")
		cancel()
		return
	}
	// TODO(amistewicz): add gid translation

	// Start sshfs command in the provided session. It will run in the foreground and it will
	// not exit even if mount fails.
	err = session.Start(fmt.Sprintf(
		"mkdir -p \"%s\" && sshfs -o idmap=user -o slave \":%s\" \"%s\"",
		sshfs.remote, sshfs.local, sshfs.remote))
	if err != nil {
		logger.WithError(err).Error("Failed to start reverse SSHFS.")
		cancel()
		return
	}

	sshfs.session = session
	sshfs.sftpCancel = cancel
	go sshfs.sshfsWait()
	go sshfs.sftpWait(sftp)
	return
}

// sshfsWait calls Wait on *ssh.Session.
// It is intended to be called in a go statement. It always sends an error to errors channel.
func (sshfs *reverseSSHFS) sshfsWait() {
	err := sshfs.session.Wait()
	sshfs.sftpCancel()
	if err != nil {
		logger.WithError(err).WithProperty("sshfsStderr", sshfs.sshfsStderr.String()).
			Errorf("SSHFS process exited with error.")
		err = fmt.Errorf("sshfs process exited: %s: %s", err, sshfs.sshfsStderr.String())
	} else {
		err = fmt.Errorf("sshfs process exited with success: %s", sshfs.sshfsStderr.String())
	}
	sshfs.errors <- err
}

// sftpWait calls Wait on *exec.Cmd.
// It is intended to be called in a go statement. It always sends an error to errors channel.
func (sshfs *reverseSSHFS) sftpWait(sftp *exec.Cmd) {
	err := sftp.Wait()
	if err != nil {
		logger.WithError(err).WithProperty("sshfsStderr", sshfs.sftpStderr.String()).
			Errorf("SFTP process exited with error.")
		err = fmt.Errorf("sftp process exited: %s: %s", err, sshfs.sftpStderr.String())
	} else {
		err = fmt.Errorf("sftp process exited with success: %s", sshfs.sftpStderr.String())
	}
	sshfs.errors <- err
}

// close terminates sshfs session.
func (sshfs *reverseSSHFS) close() (err error) {
	if sshfs.session != nil {
		// TODO gracefully shut down sshfs
		err = sshfs.session.Close()
		// drain error channel
		<-sshfs.errors
		<-sshfs.errors
	}
	close(sshfs.errors)
	return
}

// check matches the name of the mountpoint to /proc/mounts lists.
//
// It should be used when the user requires the filesystem access
// to ensure that it is active.
//
// The caller is responsible for closing the session. If an error was returned,
// this function should not be called again and either open() or close() used.
func (sshfs *reverseSSHFS) check(session *ssh.Session) (err error) {
	// Synchronization with something that is to be mounted is hard. inotifywait can't be used to
	// listen on file which does not exist. Therefore, it listens on the parent directory for file
	// creation and access. Timeout parameter is gradually changed. It works just a little better
	// than sleep.
	const mountChecker = `check() { grep -q "%s" /proc/mounts && exit 0; }
	notify() {
		if [ -x "$(which inotifywait)" ]; then
			f="$(dirname "%s")"
			mkdir -p "$f"
			inotifywait -qq -e access -t "${1:-12}" "$f"
		else
			sleep "${1:-12}"
		fi
	}
	check
	for i in $(seq 0 5); do
		notify "$i"
		check
	done
	exit 1`

	select {
	case err = <-sshfs.errors:
		// Only the first error is interesting to us
		<-sshfs.errors
		return
	default:
		err = session.Run(fmt.Sprintf(mountChecker, sshfs.remote, sshfs.remote))
		if err != nil {
			logger.WithError(err).Error("Filesystem not mounted. Session failed.")
			if err := sshfs.session.Close(); err != nil {
				logger.WithError(err).Error("Failed to close session.")
				return fmt.Errorf("filesystem not mounted: failed to close session: %s", err)
			}
			// Drain channels to free goroutines.
			<-sshfs.errors
			<-sshfs.errors
			sshfs.session = nil
			return ErrNotMounted
		}
		return
	}
}
