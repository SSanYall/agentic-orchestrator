//go:build !windows

// Copyright 2026 DoorDash, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package selfupdate

import (
	"errors"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func lockFD(fd int) error {
	if err := unix.Flock(fd, unix.LOCK_EX|unix.LOCK_NB); err != nil {
		if errors.Is(err, unix.EWOULDBLOCK) || errors.Is(err, unix.EAGAIN) {
			return errLockBusy
		}
		return err
	}
	return nil
}

func unlockFD(fd int) error { return unix.Flock(fd, unix.LOCK_UN) }

func closeFD(fd int) error { return unix.Close(fd) }

func accessWritable(path string) bool { return unix.Access(path, unix.W_OK) == nil }

func setCloseOnExec(fd int, set bool) error {
	if set {
		_, err := unix.FcntlInt(uintptr(fd), unix.F_SETFD, unix.FD_CLOEXEC)
		return err
	}
	_, err := unix.FcntlInt(uintptr(fd), unix.F_SETFD, 0)
	return err
}

func isCloseOnExec(fd int) (bool, error) {
	flags, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0)
	if err != nil {
		return false, err
	}
	return flags&unix.FD_CLOEXEC != 0, nil
}

func fileOwnerUID(info os.FileInfo) (int, bool) {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return int(st.Uid), true
	}
	return 0, false
}

func fileIdentity(info os.FileInfo) (uint64, uint64, int, int, bool) {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return uint64(st.Dev), uint64(st.Ino), int(st.Uid), int(st.Gid), true
	}
	return 0, 0, 0, 0, false
}

func fileNLink(info os.FileInfo) (uint32, bool) {
	if st, ok := info.Sys().(*syscall.Stat_t); ok {
		return st.Nlink, true
	}
	return 0, false
}

func processAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	err := unix.Kill(pid, 0)
	return err == nil || errors.Is(err, unix.EPERM)
}

func syncFD(fd int) error { return unix.Fsync(fd) }
