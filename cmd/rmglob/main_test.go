/*
Copyright (c) NVIDIA CORPORATION.  All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var rmglobBin string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "rmglob-test-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "tempdir: %v\n", err)
		os.Exit(2)
	}

	rmglobBin = filepath.Join(dir, "rmglob")
	if out, err := exec.Command("go", "build", "-o", rmglobBin, ".").CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build: %v\n%s", err, out)
		os.RemoveAll(dir)
		os.Exit(2)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

func TestRmglob(t *testing.T) {
	tmpDir := t.TempDir()
	for _, name := range []string{"a-ready", "b-ready", "keep.txt"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("x"), 0600); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	//#nosec G204 -- test-only invocation of a binary built by TestMain
	if out, err := exec.Command(rmglobBin, filepath.Join(tmpDir, "*-ready")).CombinedOutput(); err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, "a-ready")); !os.IsNotExist(err) {
		t.Errorf("a-ready should be removed, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "b-ready")); !os.IsNotExist(err) {
		t.Errorf("b-ready should be removed, stat err=%v", err)
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "keep.txt")); err != nil {
		t.Errorf("keep.txt should remain, stat err=%v", err)
	}
}

func TestRmglobNoArgs(t *testing.T) {
	if err := exec.Command(rmglobBin).Run(); err == nil {
		t.Errorf("rmglob with no args expected non-zero exit, got 0")
	}
}
