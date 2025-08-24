// taken from `https://github.com/nsf/gothic/:
// - https://github.com/nsf/gothic/blob/97dfcc195b9de36c911a69a6ec2b5b2659c05652/fmt.go

// Copyright (C) 2013 https://github.com/nsf
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tk

import (
	"bytes"
	"unicode"
	"unicode/utf8"
)

func quote_rune(buf *bytes.Buffer, r rune, size int) {
	const lowerhex = "0123456789abcdef"
	if size == 1 && r == utf8.RuneError {
		// invalid rune, write the byte as is
		buf.WriteString(`\x`)
		buf.WriteByte(lowerhex[r>>4])
		buf.WriteByte(lowerhex[r&0xF])
		return
	}

	// first check for special TCL escaping cases
	switch r {
	case '{', '}', '[', ']', '"', '$', '\\':
		buf.WriteString("\\")
		buf.WriteRune(r)
		return
	}

	// other printable characters
	if unicode.IsPrint(r) {
		buf.WriteRune(r)
		return
	}

	// non-printable characters
	switch r {
	case '\a':
		buf.WriteString(`\a`)
	case '\b':
		buf.WriteString(`\b`)
	case '\f':
		buf.WriteString(`\f`)
	case '\n':
		buf.WriteString(`\n`)
	case '\r':
		buf.WriteString(`\r`)
	case '\t':
		buf.WriteString(`\t`)
	case '\v':
		buf.WriteString(`\v`)
	default:
		switch {
		case r < ' ':
			buf.WriteString(`\x`)
			buf.WriteByte(lowerhex[r>>4])
			buf.WriteByte(lowerhex[r&0xF])
		case r >= 0x10000:
			r = 0xFFFD
			fallthrough
		case r < 0x10000:
			buf.WriteString(`\u`)
			for s := 12; s >= 0; s -= 4 {
				buf.WriteByte(lowerhex[r>>uint(s)&0xF])
			}
		}
	}
}

func quote(buf *bytes.Buffer, s string) {
	buf.WriteString(`"`)
	size := 0
	for offset := 0; offset < len(s); offset += size {
		r := rune(s[offset])
		size = 1
		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRuneInString(s[offset:])
		}

		quote_rune(buf, r, size)
	}
	buf.WriteString(`"`)
}

// Works exactly like Eval("%{%q}"), but instead of evaluating returns a quoted string.
func Quote(s string) string {
	var tmp bytes.Buffer
	quote(&tmp, s)
	return tmp.String()
}

// Quotes each string in given `string_list`.
func QuoteAll(string_list []string) []string {
	new_string_list := []string{}
	for _, str := range string_list {
		new_string_list = append(new_string_list, Quote(str))
	}
	return new_string_list
}

// Quotes the rune just like if it was passed through Quote, the result is the
// same as: Quote(string(r)).
func QuoteRune(r rune) string {
	var tmp bytes.Buffer
	size := utf8.RuneLen(r)
	tmp.WriteString(`"`)
	quote_rune(&tmp, r, size)
	tmp.WriteString(`"`)
	return tmp.String()
}
