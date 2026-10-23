/*
Copyright 2022 Google LLC

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

package creds

import (
	"io"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	gitlab "github.com/ePirat/docker-credential-gitlabci/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

func GetKeychain() authn.Keychain {
	return authn.NewMultiKeychain(
		authn.DefaultKeychain,
		google.Keychain,
		authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithLogger(io.Discard))),
		// The ACR helper is wrapped so that only genuine ACR/MCR hosts reach it.
		// See acr.go and GO-2026-6225.
		authn.NewKeychainFromHelper(newAnchoredACRHelper(credhelper.NewACRCredentialsHelper())),
		authn.NewKeychainFromHelper(gitlab.NewGitLabCredentialsHelper()),
	)
}
