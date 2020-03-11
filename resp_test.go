// Copyright 2020 Ross Light
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"strings"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		s   string
		tok string
		err bool
	}{
		{s: "", tok: ""},
		{s: "foo", err: true},
		{s: "+OK\r\n", tok: "+OK\r\n"},
		{s: "-Error message\r\n", tok: "-Error message\r\n"},
		{s: ":0\r\n", tok: ":0\r\n"},
		{s: ":1000\r\n", tok: ":1000\r\n"},
		{s: "$6\r\nfoobar\r\n", tok: "$6\r\nfoobar\r\n"},
		{s: "$0\r\n\r\n", tok: "$0\r\n\r\n"},
		{s: "$-1\r\n", tok: "$-1\r\n"},
		{s: "$-2\r\n", err: true},
		{s: "$-1\r\n+OK\r\n", tok: "$-1\r\n"},
		{s: "*0\r\n", tok: "*0\r\n"},
		{s: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n", tok: "*2\r\n"},
		{s: "*-1\r\n", tok: "*-1\r\n"},
		{s: "$-1\r\n+OK\r\n", tok: "$-1\r\n"},
	}
	for _, test := range tests {
		if !strings.HasPrefix(test.s, test.tok) {
			t.Errorf("Test for %q has invalid token %q. Skipping.", test.s, test.tok)
			continue
		}
		got, err := lex([]byte(test.s))
		if want := len(test.tok); got != want || (err != nil) != test.err {
			wantErr := "<nil>"
			if test.err {
				wantErr = "<some error>"
			}
			if got < len(test.s) {
				t.Errorf("lex(%q) = %d (%q), %v; want %d (%q), %s", test.s, got, test.s[:got], err, want, test.tok, wantErr)
			} else {
				t.Errorf("lex(%q) = %d, %v; want %d (%q), %s", test.s, got, err, want, test.tok, wantErr)
			}
		}
	}
}
