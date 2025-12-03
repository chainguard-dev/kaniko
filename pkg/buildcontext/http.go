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

package buildcontext

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	kConfig "github.com/chainguard-dev/kaniko/pkg/config"
	"github.com/chainguard-dev/kaniko/pkg/constants"
	"github.com/chainguard-dev/kaniko/pkg/util"
	"github.com/sirupsen/logrus"
)

// HTTPContext struct for http/https tar.gz files processing
type HTTPContext struct {
	context string
}

// UnpackTarFromBuildContext downloads context file from http/https server
func (h *HTTPContext) UnpackTarFromBuildContext() (directory string, err error) {

	logrus.Info("Retrieving tar file from URL")

	// Create directory and target file for downloading the context file
	directory = kConfig.BuildContextDir
	tarPath := filepath.Join(directory, constants.ContextTar)
	file, err := util.CreateTargetTarfile(tarPath)
	if err != nil {
		return
	}

	// Download tar file from remote https server
	// and save it into the target tar file
	resp, err := http.Get(h.context) //nolint:noctx
	if err != nil {
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return directory, fmt.Errorf("HTTP bad status from server: %s", resp.Status)
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		return tarPath, err
	}

	logrus.Info("Retrieved tar file from URL")

	if err = util.UnpackCompressedTar(tarPath, directory); err != nil {
		return
	}

	logrus.Info("Extracted tar file from URL")

	// Remove the tar so it doesn't interfere with subsequent commands
	return directory, os.Remove(tarPath)
}
