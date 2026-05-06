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

// rmglob is a tiny static helper binary that expands one or more glob
// patterns and removes the matching paths. It exists so that distroless
// gpu-operator container images can run path cleanup from a Kubernetes
// `lifecycle.preStop` hook without needing a shell on the image.
//
// It is the path-cleanup analog of k8s-cc-manager's vendored static `/bin/rm`.
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: rmglob <glob>...")
		os.Exit(2)
	}

	var failed bool
	for _, pattern := range os.Args[1:] {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			//#nosec G705 -- stderr diagnostic, not a network-reachable sink
			fmt.Fprintf(os.Stderr, "rmglob: invalid pattern %q: %v\n", pattern, err)
			failed = true
			continue
		}
		for _, m := range matches {
			// Path removal is the binary's sole purpose; the patterns come from
			// gpu-operator-rendered manifests, not external user input.
			//#nosec G703 -- intentional path removal
			if err := os.RemoveAll(m); err != nil {
				//#nosec G705 -- stderr diagnostic, not a network-reachable sink
				fmt.Fprintf(os.Stderr, "rmglob: remove %q: %v\n", m, err)
				failed = true
			}
		}
	}
	if failed {
		os.Exit(1)
	}
}
