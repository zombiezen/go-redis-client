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
	"bytes"
	"strconv"

	"golang.org/x/xerrors"
)

const terminator = "\r\n"

var terminatorBytes = []byte(terminator)

const maxBulkStringLength = 512 << 20

// lex parses the next RESP token from a stream. An error is returned if the
// input is unparseable. If the input is an incomplete token, then lex returns
// ("", s, nil).
//
// Array headers are treated as their own token and do not include the rest of
// the array.
func lex(s []byte) (tok string, rest []byte, err error) {
	if len(s) == 0 {
		return "", s, nil
	}
	i := bytes.Index(s, terminatorBytes)
	switch s[0] {
	case '+', '-', ':', '*':
		if i == -1 {
			return "", s, nil
		}
		return string(s[:i+len(terminator)]), s[i+len(terminator):], nil
	case '$': // Bulk string
		if i == -1 {
			return "", s, nil
		}
		bulkLength, err := strconv.ParseInt(string(s[1:i]), 10, 32)
		if err != nil {
			return "", s, xerrors.Errorf("parse RESP: bad bulk string length: %w", err)
		}
		if bulkLength < -1 || bulkLength > maxBulkStringLength {
			return "", s, xerrors.Errorf("parse RESP: invalid bulk string length %d", bulkLength)
		}
		if bulkLength == -1 {
			return string(s[:i+len(terminator)]), s[i+len(terminator):], nil
		}
		stringEnd := i + len(terminator) + int(bulkLength)
		if len(s) < stringEnd+len(terminator) {
			return "", s, nil
		}
		if !bytes.HasPrefix(s[stringEnd:], terminatorBytes) {
			return "", s, xerrors.New("parse RESP: unterminated bulk string")
		}
		return string(s[:stringEnd+len(terminator)]), s[stringEnd+len(terminator):], nil
	default:
		return "", s, xerrors.Errorf("parse RESP: invalid tag %q", s[0])
	}
}
