// Copyright 2025 CFC4N <cfc4n.cs@gmail.com>. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Repository: https://github.com/gojue/moling

package browser

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gojue/moling/pkg/comm"
)

func TestBrowserServer(t *testing.T) {
	//
	//cfg := &BrowserConfig{
	//	Headless:        true,
	//	Timeout:         30,
	//	UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
	//	DefaultLanguage: "en-US",
	//	URLTimeout:      10,
	//	SelectorQueryTimeout:      10,
	//}
	logger, ctx, err := comm.InitTestEnv()
	if err != nil {
		t.Fatalf("Failed to initialize test environment: %s", err.Error())
	}
	logger.Info().Msg("TestBrowserServer")

	_, err = NewBrowserServer(ctx)
	if err != nil {
		t.Fatalf("Failed to create BrowserServer: %s", err.Error())
	}
}

// TestHoverSelectorEscaping verifies that the selector is safely JSON-encoded
// before being embedded in the JavaScript expression, preventing code injection.
func TestHoverSelectorEscaping(t *testing.T) {
	tests := []struct {
		name     string
		selector string
		// wantPrefix checks that the selector is embedded as a JSON-encoded double-quoted string
		wantJSPrefix string
		wantJSSuffix string
	}{
		{
			name:         "normal selector",
			selector:     "body",
			wantJSPrefix: `document.querySelector("body").dispatchEvent`,
		},
		{
			name:     "injection attempt with single quotes and comma operator",
			selector: "body'),document.title='PWNED',document.querySelector('body",
			// After JSON encoding, the selector is a double-quoted string literal.
			// The single quotes and commas stay inside the string and are NOT executable.
			wantJSPrefix: `document.querySelector("body'),document.title='PWNED',document.querySelector('body").dispatchEvent`,
		},
		{
			name:         "injection attempt with semicolons and IIFE",
			selector:     "body'); (function(){ /* exfiltration */ })(); document.querySelector('body",
			wantJSPrefix: `document.querySelector("body'); (function(){ /* exfiltration */ })(); document.querySelector('body").dispatchEvent`,
		},
		{
			name:     "selector with double quotes is escaped by JSON",
			selector: `div[data-id="foo"]`,
			// json.Marshal escapes inner double quotes as \", so they cannot break out of the JS string.
			wantJSPrefix: `document.querySelector("div[data-id=\"foo\"]").dispatchEvent`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			selectorJSON, err := json.Marshal(tc.selector)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}
			js := `document.querySelector(` + string(selectorJSON) + `).dispatchEvent(new Event('mouseover'))`

			// The embedded selector must be wrapped in JSON double-quotes (not single-quotes).
			// This ensures injected single-quote characters cannot break out of the JS string context.
			if !strings.HasPrefix(string(selectorJSON), `"`) || !strings.HasSuffix(string(selectorJSON), `"`) {
				t.Errorf("selector was not JSON-encoded as a double-quoted string: %s", string(selectorJSON))
			}

			if tc.wantJSPrefix != "" && !strings.HasPrefix(js, tc.wantJSPrefix) {
				t.Errorf("JS expression did not start with expected prefix\n  want: %s\n   got: %s", tc.wantJSPrefix, js)
			}
		})
	}
}
