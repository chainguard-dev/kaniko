/*
Copyright 2018 Google LLC

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

package integration

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestWithContext(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	dir := filepath.Join(cwd, "dockerfiles-with-context")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	testDirs := make([]fs.FileInfo, 0, len(entries))

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			t.Fatal(err)
		}
		testDirs = append(testDirs, info)
	}

	builder := NewDockerFileBuilder()

	for _, tdInfo := range testDirs {
		name := tdInfo.Name()
		testDir := filepath.Join(dir, name)

		t.Run("test_with_context_"+name, func(t *testing.T) {
			t.Parallel()

			if err := builder.BuildImageWithContext(
				t, config, "", name, testDir,
			); err != nil {
				t.Fatal(err)
			}

			dockerImage := GetDockerImage(config.imageRepo, name)
			kanikoImage := GetKanikoImage(config.imageRepo, name)

			containerDiff(t, daemonPrefix+dockerImage, kanikoImage, "--semantic", "--extra-ignore-file-permissions", "--extra-ignore-file-content", "--extra-ignore-layer-length-mismatch")
		})
	}

	if err := logBenchmarks("benchmark"); err != nil {
		t.Logf("Failed to create benchmark file: %v", err)
	}
}
