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
	"context"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
)

// recordingHelper stands in for the real ACR helper. It records every call
// and hands back a fake token, so the tests never touch Azure.
type recordingHelper struct {
	calls []string
}

func (r *recordingHelper) Get(serverURL string) (string, string, error) {
	r.calls = append(r.calls, serverURL)
	return "<token>", "fake-acr-refresh-token", nil
}

// Hosts that must never reach the wrapped helper.
var rejectedHosts = []string{
	"evil.azurecr.io.attacker.com", // GO-2026-6225 lookalike
	"evil.azurecr.io.attacker.com:443",
	"foo.azurecr.io.",                // trailing dot
	"foo.azurecr.io@attacker.com",    // userinfo, real host is attacker.com
	"attacker.com@foo.azurecr.io",    // userinfo on an ACR host
	"user:pass@foo.azurecr.io",       // userinfo with password
	"attacker.com/foo.azurecr.io",    // path
	"attacker.com#.azurecr.io",       // fragment
	"attacker.com?x=.azurecr.io",     // query
	"foo.azurecr.io/v2/",             // path on a real ACR host
	"azurecr.io",                     // apex, no registry name
	".azurecr.io",                    // empty label
	"foo.azurecr.me",                 // not an ACR domain
	"fooazurecr.io",                  // missing dot
	"foo.azurecr.io.evil",            // suffix
	"mcr.microsoft.com.attacker.com", // MCR lookalike
	"index.docker.io",
	"gcr.io",
	"localhost:5000",
	"127.0.0.1:12345",
	"",
	"notaurl-)(*$@)(*@)(*",
}

// Hosts that must still be passed through to the wrapped helper.
var acceptedHosts = []string{
	"foo.azurecr.io",
	"foo.azurecr.us",
	"foo.azurecr.cn",
	"foo.azurecr.de",
	"foo.azurecr.io:443",
	"foo.westeurope.data.azurecr.io", // dedicated data endpoint
	"mcr.microsoft.com",
}

func TestIsACRHost(t *testing.T) {
	for _, h := range rejectedHosts {
		if isACRHost(h) {
			t.Errorf("isACRHost(%q) = true, want false", h)
		}
	}
	for _, h := range acceptedHosts {
		if !isACRHost(h) {
			t.Errorf("isACRHost(%q) = false, want true", h)
		}
	}
}

func TestAnchoredACRHelperDoesNotCallInnerForNonACRHosts(t *testing.T) {
	for _, h := range rejectedHosts {
		inner := &recordingHelper{}
		u, p, err := newAnchoredACRHelper(inner).Get(h)
		if err == nil {
			t.Errorf("Get(%q): expected an error, got none", h)
		}
		if u != "" || p != "" {
			t.Errorf("Get(%q) returned credentials %q/%q, want none", h, u, p)
		}
		if len(inner.calls) != 0 {
			t.Errorf("Get(%q): wrapped helper was called %d times, want 0", h, len(inner.calls))
		}
	}
}

func TestAnchoredACRHelperPassesThroughACRHosts(t *testing.T) {
	for _, h := range acceptedHosts {
		inner := &recordingHelper{}
		u, p, err := newAnchoredACRHelper(inner).Get(h)
		if err != nil {
			t.Errorf("Get(%q): unexpected error %v", h, err)
		}
		if u != "<token>" || p != "fake-acr-refresh-token" {
			t.Errorf("Get(%q) = %q/%q, want the wrapped helper's credentials", h, u, p)
		}
		if len(inner.calls) != 1 || inner.calls[0] != h {
			t.Errorf("Get(%q): wrapped helper calls = %v, want exactly [%q]", h, inner.calls, h)
		}
	}
}

// TestKeychainResolvesAnonymousForLookalike checks the end-to-end path used by
// GetKeychain: a lookalike registry resolves to Anonymous and the wrapped
// helper never runs, while a real ACR registry gets the token.
func TestKeychainResolvesAnonymousForLookalike(t *testing.T) {
	inner := &recordingHelper{}
	kc := authn.NewKeychainFromHelper(newAnchoredACRHelper(inner))

	evil, err := name.NewRegistry("evil.azurecr.io.attacker.com")
	if err != nil {
		t.Fatal(err)
	}
	auth, err := authn.Resolve(context.Background(), kc, evil)
	if err != nil {
		t.Fatal(err)
	}
	if auth != authn.Anonymous {
		t.Errorf("lookalike registry resolved to %v, want Anonymous", auth)
	}
	if len(inner.calls) != 0 {
		t.Errorf("wrapped helper called for lookalike registry: %v", inner.calls)
	}

	good, err := name.NewRegistry("foo.azurecr.io")
	if err != nil {
		t.Fatal(err)
	}
	auth, err = authn.Resolve(context.Background(), kc, good)
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := authn.Authorization(context.Background(), auth)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.IdentityToken != "fake-acr-refresh-token" {
		t.Errorf("ACR registry got identity token %q, want the wrapped helper's token", cfg.IdentityToken)
	}
}
