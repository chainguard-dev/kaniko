/*
Copyright 2026 Google LLC

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
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/go-containerregistry/pkg/authn"
)

// acrHostRE matches Azure Container Registry hostnames. Unlike the regex in
// github.com/chrismellard/docker-credential-acr-env (GO-2026-6225), it is
// anchored at both ends, so lookalike hosts such as
// "evil.azurecr.io.attacker.com" do not match.
var acrHostRE = regexp.MustCompile(`^.+\.azurecr\.(io|cn|de|us)$`)

// mcrHostname is accepted because the wrapped helper accepts it too; keeping
// it means the wrapper only ever narrows what the helper sees.
const mcrHostname = "mcr.microsoft.com"

// anchoredACRHelper guards an ACR credential helper so that it is only called
// for hosts that are genuinely ACR (or MCR). The upstream helper's own check is
// unanchored and would fetch an Azure AD token and send it to any host that
// merely contains ".azurecr.io". Upstream is unmaintained, so the check is
// done here, before the helper is ever called.
type anchoredACRHelper struct {
	inner authn.Helper
}

func newAnchoredACRHelper(inner authn.Helper) authn.Helper {
	return anchoredACRHelper{inner: inner}
}

// Get returns an error for anything that is not a plain ACR/MCR host, which
// authn.NewKeychainFromHelper treats as "no credentials" (anonymous). That is
// the same outcome the wrapped helper gives for non-ACR hosts today.
func (h anchoredACRHelper) Get(serverURL string) (string, string, error) {
	if !isACRHost(serverURL) {
		return "", "", fmt.Errorf("serverURL %q does not refer to Azure Container Registry", serverURL)
	}
	return h.inner.Get(serverURL)
}

// isACRHost reports whether serverURL is a bare registry host (optionally with
// a port) on an ACR domain, or mcr.microsoft.com.
//
// It is deliberately strict: the wrapped helper builds its token exchange URL
// as "https://" + serverURL, so anything that is not exactly host[:port]
// (userinfo, paths, queries, fragments, trailing dots) is rejected rather than
// interpreted.
func isACRHost(serverURL string) bool {
	u, err := url.Parse("https://" + serverURL)
	if err != nil {
		return false
	}
	// Reject anything other than a bare host[:port].
	if u.Host != serverURL || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.Opaque != "" {
		return false
	}
	host := u.Hostname()
	if host == "" {
		return false
	}
	return host == mcrHostname || acrHostRE.MatchString(host)
}
