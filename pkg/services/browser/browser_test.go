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
		t.Fatalf("Failed to initialize test environment: %w", err)
	}
	logger.Info().Msg("TestBrowserServer")

	_, err = NewBrowserServer(ctx)
	if err != nil {
		t.Fatalf("Failed to create BrowserServer: %w", err)
	}
}
