//go:build windows

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

import "os"

var errLockBusy = os.ErrExist

func lockFD(int) error { return nil }

func unlockFD(int) error { return nil }

func closeFD(int) error { return nil }

func accessWritable(string) bool { return true }

func setCloseOnExec(int, bool) error { return nil }

func isCloseOnExec(int) (bool, error) { return false, nil }

func fileOwnerUID(os.FileInfo) (int, bool) { return 0, true }

func fileIdentity(os.FileInfo) (uint64, uint64, int, int, bool) { return 0, 0, 0, 0, false }

func fileNLink(os.FileInfo) (uint32, bool) { return 1, true }

func processAlive(pid int) bool { return pid > 0 }

func syncFD(int) error { return nil }
