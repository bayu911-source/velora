// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name     string
		greeting string
		path     string
		want     string
	}{
		{
			name:     "world",
			greeting: "Hello",
			path:     "/",
			want:     "Hello, world!",
		},
		{
			name:     "explicit name",
			greeting: "Hi",
			path:     "/Gopher",
			want:     "Hi, Gopher!",
		},
		{
			name:     "escaped name",
			greeting: "Hey",
			path:     "/<script>alert('oops')</script>",
			want:     "Hey, &lt;script&gt;alert(&#39;oops&#39;)&lt;/script&gt;!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()
			helloHandler(tt.greeting)(w, req)
			res := w.Result()
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response body: %v", err)
			}
			if string(body) != tt.want {
				t.Errorf("got %q, want %q", string(body), tt.want)
			}
		})
	}
}

func TestVersionHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()
	const version = "v1.2.3"
	versionHandler(version)(w, req)
	res := w.Result()
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
			t.Fatalf("could not read response body: %v", err)
	}
	if string(body) != version+"\n" {
			t.Errorf("got %q, want %q", string(body), version+"\n")
	}
}
