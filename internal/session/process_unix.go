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

package session

import (
	"os"
	"os/exec"
	"syscall"
)

func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killProcessGroup(pgid int, sig syscall.Signal) error {
	if pgid <= 0 {
		return nil
	}
	if err := syscall.Kill(-pgid, sig); err != nil && !isNoSuchProcess(err) {
		return err
	}
	return nil
}

func waitProcessNoHang(pid int) (int, error) {
	var ws syscall.WaitStatus
	return syscall.Wait4(pid, &ws, syscall.WNOHANG, nil)
}

func isNoSuchProcess(err error) bool { return err == syscall.ESRCH }

func processAlive(pid int) bool {
	if pid <= 0 { return false }
	if err := syscall.Kill(-pid, 0); err == nil || err == syscall.EPERM { return true }
	return false
}

func stopSessionProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil { return }
	_ = cmd.Process.Kill()
	_ = cmd.Wait()
}

func processExists(pid int) bool { return os.Getpid() != pid }
