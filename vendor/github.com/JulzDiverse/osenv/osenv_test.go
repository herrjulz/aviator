package osenv_test
  
  import (
  	. "osenv"
  	"reflect"
  	"strings"
  	"testing"
  )
  
  // testGetenv gives us a controlled set of variables for testing Expand.
  func testGetenv(s string) string {
  	switch s {
  	case "*":
  		return "all the args"
  	case "#":
  		return "NARGS"
  	case "$":
  		return "PID"
  	case "1":
  		return "ARGUMENT1"
  	case "HOME":
  		return "/usr/gopher"
  	case "H":
  		return "(Value of H)"
  	case "home_1":
  		return "/usr/foo"
  	case "_":
  		return "underscore"
  	}
  	return ""
  }
  
  var expandTests = []struct {
  	in, out string
  }{
  	{"", ""},
  	{"$*", "all the args"},
  	{"$$", "PID"},
  	{"${*}", "all the args"},
  	{"$1", "ARGUMENT1"},
  	{"${1}", "ARGUMENT1"},
  	{"now is the time", "now is the time"},
  	{"$HOME", "/usr/gopher"},
  	{"$home_1", "/usr/foo"},
  	{"${HOME}", "/usr/gopher"},
  	{"${H}OME", "(Value of H)OME"},
  	{"A$$$#$1$H$home_1*B", "APIDNARGSARGUMENT1(Value of H)/usr/foo*B"},
  }
  
  func TestExpand(t *testing.T) {
  	for _, test := range expandTests {
  		result := Expand(test.in, testGetenv)
  		if result != test.out {
  			t.Errorf("Expand(%q)=%q; expected %q", test.in, result, test.out)
  		}
  	}
  }
  
  func TestConsistentEnviron(t *testing.T) {
  	e0 := Environ()
  	for i := 0; i < 10; i++ {
  		e1 := Environ()
  		if !reflect.DeepEqual(e0, e1) {
  			t.Fatalf("environment changed")
  		}
  	}
  }
  
  func TestUnsetenv(t *testing.T) {
  	const testKey = "GO_TEST_UNSETENV"
  	set := func() bool {
  		prefix := testKey + "="
  		for _, key := range Environ() {
  			if strings.HasPrefix(key, prefix) {
  				return true
  			}
  		}
  		return false
  	}
  	if err := Setenv(testKey, "1"); err != nil {
  		t.Fatalf("Setenv: %v", err)
  	}
  	if !set() {
  		t.Error("Setenv didn't set TestUnsetenv")
  	}
  	if err := Unsetenv(testKey); err != nil {
  		t.Fatalf("Unsetenv: %v", err)
  	}
  	if set() {
  		t.Fatal("Unsetenv didn't clear TestUnsetenv")
  	}
  }
  
  func TestClearenv(t *testing.T) {
  	const testKey = "GO_TEST_CLEARENV"
  	const testValue = "1"
  
  	// reset env
  	defer func(origEnv []string) {
  		for _, pair := range origEnv {
  			// Environment variables on Windows can begin with =
  			// http://blogs.msdn.com/b/oldnewthing/archive/2010/05/06/10008132.aspx
  			i := strings.Index(pair[1:], "=") + 1
  			if err := Setenv(pair[:i], pair[i+1:]); err != nil {
  				t.Errorf("Setenv(%q, %q) failed during reset: %v", pair[:i], pair[i+1:], err)
  			}
  		}
  	}(Environ())
  
  	if err := Setenv(testKey, testValue); err != nil {
  		t.Fatalf("Setenv(%q, %q) failed: %v", testKey, testValue, err)
  	}
  	if _, ok := LookupEnv(testKey); !ok {
  		t.Errorf("Setenv(%q, %q) didn't set $%s", testKey, testValue, testKey)
  	}
  	Clearenv()
  	if val, ok := LookupEnv(testKey); ok {
  		t.Errorf("Clearenv() didn't clear $%s, remained with value %q", testKey, val)
  	}
  }
  
  func TestLookupEnv(t *testing.T) {
  	const smallpox = "SMALLPOX"      // No one has smallpox.
  	value, ok := LookupEnv(smallpox) // Should not exist.
  	if ok || value != "" {
  		t.Fatalf("%s=%q", smallpox, value)
  	}
  	defer Unsetenv(smallpox)
  	err := Setenv(smallpox, "virus")
  	if err != nil {
  		t.Fatalf("failed to release smallpox virus")
  	}
  	value, ok = LookupEnv(smallpox)
  	if !ok {
  		t.Errorf("smallpox release failed; world remains safe but LookupEnv is broken")
  	}
  }

// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// General environment variables.

package osenv

import "syscall"

// Expand replaces ${var} or $var in the string based on the mapping function.
// For example, os.ExpandEnv(s) is equivalent to os.Expand(s, os.Getenv).
func Expand(s string, mapping func(string) string) string {
	buf := make([]byte, 0, 2*len(s))
	// ${} is all ASCII, so bytes are fine for this operation.
	i := 0
	for j := 0; j < len(s); j++ {
		if s[j] == '$' && j+1 < len(s) {
			buf = append(buf, s[i:j]...)
			name, w := getShellName(s[j+1:])
			buf = append(buf, mapping(name)...)
			j += w
			i = j + 1
		}
	}
	return string(buf) + s[i:]
}

// ExpandEnv replaces ${var} or $var in the string according to the values
// of the current environment variables. References to undefined
// variables are replaced by the empty string.
func ExpandEnv(s string) string {
	return Expand(s, Getenv)
}

// isShellSpecialVar reports whether the character identifies a special
// shell variable such as $*.
func isShellSpecialVar(c uint8) bool {
	switch c {
	case '*', '#', '$', '@', '!', '?', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return true
	}
	return false
}

// isAlphaNum reports whether the byte is an ASCII letter, number, or underscore
func isAlphaNum(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
