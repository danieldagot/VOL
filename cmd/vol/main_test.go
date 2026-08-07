package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "program.vol")
	if err := os.WriteFile(path, []byte("print 42\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"run", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code %d, stderr %q", code, stderr.String())
	}
	if stdout.String() != "42\n" {
		t.Fatalf("stdout %q", stdout.String())
	}
}

func TestLegacyFileCommand(t *testing.T) {
	path := filepath.Join(t.TempDir(), "program.vol")
	if err := os.WriteFile(path, []byte("print 7\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{path}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code %d, stderr %q", code, stderr.String())
	}
	if stdout.String() != "7\n" {
		t.Fatalf("stdout %q", stdout.String())
	}
}

func TestRunUsage(t *testing.T) {
	for _, args := range [][]string{nil, {}, {"run"}, {"one.vol", "extra"}} {
		var stdout, stderr bytes.Buffer
		if code := run(args, &stdout, &stderr); code != 2 {
			t.Fatalf("args %#v: exit code %d, want 2", args, code)
		}
		if stdout.Len() != 0 || !strings.Contains(stderr.String(), "usage: vol run <file.vol>") {
			t.Fatalf("args %#v: stdout %q, stderr %q", args, stdout.String(), stderr.String())
		}
	}
}

func TestRunForwardsProgramArguments(t *testing.T) {
	path := filepath.Join(t.TempDir(), "arguments.vol")
	if err := os.WriteFile(path, []byte("print args\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "separator", args: []string{"run", path, "--", "first", "second"}, want: "[first, second]\n"},
		{name: "no separator", args: []string{"run", path, "first"}, want: "[first]\n"},
		{name: "empty after separator", args: []string{"run", path, "--"}, want: "[]\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(test.args, &stdout, &stderr); code != 0 || stderr.Len() != 0 || stdout.String() != test.want {
				t.Fatalf("code/output = %d, %q, %q", code, stdout.String(), stderr.String())
			}
		})
	}
}

func TestRunReportsFileParseAndRuntimeFailures(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		source     *string
		stderrText string
	}{
		{name: "missing file", path: filepath.Join(t.TempDir(), "missing.vol"), stderrText: "vol:"},
		{name: "parse", source: pointer("print }"), stderrText: "error[E101]"},
		{name: "resolver", source: pointer("print missing"), stderrText: "error[S002]"},
		{name: "runtime", source: pointer("print 1 / 0"), stderrText: "error[R014]"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := test.path
			if test.source != nil {
				path = filepath.Join(t.TempDir(), test.name+".vol")
				if err := os.WriteFile(path, []byte(*test.source), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			var stdout, stderr bytes.Buffer
			if code := run([]string{"run", path}, &stdout, &stderr); code != 1 {
				t.Fatalf("exit code = %d", code)
			}
			if stdout.Len() != 0 || !strings.Contains(stderr.String(), test.stderrText) {
				t.Fatalf("stdout %q, stderr %q", stdout.String(), stderr.String())
			}
		})
	}
}

func pointer(value string) *string { return &value }
