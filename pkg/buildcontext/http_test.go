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
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildWithHttpsTar(t *testing.T) {

	tests := []struct {
		name          string
		serverHandler http.HandlerFunc
	}{
		{
			name: "test http bad status",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("corrupted message"))
				if err != nil {
					t.Fatalf("Error sending response: %v", err)
				}
			}),
		},
		{
			name: "test http bad data",
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := w.Write([]byte("corrupted message"))
				if err != nil {
					t.Fatalf("Error sending response: %v", err)
				}
			}),
		},
	}

	for _, tcase := range tests {
		t.Run(tcase.name, func(t *testing.T) {
			server := httptest.NewServer(tcase.serverHandler)
			defer server.Close()

			context := &HTTPContext{
				context: server.URL + "/data.tar.gz",
			}

			_, err := context.UnpackTarFromBuildContext()
			if err == nil {
				t.Fatalf("Error expected but not returned: %s", err)
			}
		})
	}
}

func TestGetBuildContextHttpPrefix(t *testing.T) {
	tests := []struct {
		name        string
		srcContext  string
		expectedType string
	}{
		{
			name:        "http prefix",
			srcContext:  "http://example.com/context.tar.gz",
			expectedType: "*buildcontext.HTTPContext",
		},
		{
			name:        "https prefix",
			srcContext:  "https://example.com/context.tar.gz",
			expectedType: "*buildcontext.HTTPContext",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bc, err := GetBuildContext(test.srcContext, BuildOptions{})
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			
			actualType := fmt.Sprintf("%T", bc)
			if actualType != test.expectedType {
				t.Errorf("Expected type %s, got %s", test.expectedType, actualType)
			}
		})
	}
}
