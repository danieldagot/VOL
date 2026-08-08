package lang

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunFileLoadsImportsAndExports(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "vol.config.json"), []byte(`{"name":"t","root":".","paths":{"@lib":"lib"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "lib"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "lib", "math.vol"), []byte("export add\nfn add(a, b) {\n    return a + b\n}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(root, "main.vol")
	if err := os.WriteFile(entry, []byte("import \"@lib/math\"\nprint add(20, 22)\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if d := RunFile(entry, &output, ExecuteOptions{}); d != nil {
		t.Fatal(d.Human(""))
	}
	if output.String() != "42\n" {
		t.Fatalf("got %q", output.String())
	}
}

func TestRunFileResolvesModVol(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "pkg", "mod.vol"), []byte("export value\nvalue := 9\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(root, "main.vol")
	if err := os.WriteFile(entry, []byte("import \"pkg\"\nprint value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if d := RunFile(entry, &output, ExecuteOptions{}); d != nil {
		t.Fatal(d.Human(""))
	}
	if output.String() != "9\n" {
		t.Fatalf("got %q", output.String())
	}
}

func TestRunFileDetectsImportCycle(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.vol"), []byte("import \"b\"\nexport x\nx := 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.vol"), []byte("import \"a\"\nexport y\ny := 2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	d := RunFile(filepath.Join(root, "a.vol"), &output, ExecuteOptions{})
	if d == nil || d.Code != "S033" || !strings.Contains(d.Message, "Import cycle") {
		t.Fatalf("got %#v", d)
	}
}

func TestRunFileMissingModule(t *testing.T) {
	root := t.TempDir()
	entry := filepath.Join(root, "main.vol")
	if err := os.WriteFile(entry, []byte("import \"missing\"\nprint 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	d := RunFile(entry, &bytes.Buffer{}, ExecuteOptions{})
	if d == nil || d.Code != "S031" {
		t.Fatalf("got %#v", d)
	}
}

func TestExecuteRejectsImportsWithoutLoader(t *testing.T) {
	program, d := Parse("t.vol", "import \"x\"\nprint 1\n")
	if d != nil {
		t.Fatal(d)
	}
	d = Execute(program, &bytes.Buffer{})
	if d == nil || d.Code != "S031" {
		t.Fatalf("got %#v", d)
	}
}

func TestRunFileImportCollisionPointsAtLaterImport(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "a.vol"), []byte("export shared\nshared := 1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "b.vol"), []byte("export shared\nshared := 2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	entry := filepath.Join(root, "main.vol")
	source := "import \"a\"\nimport \"b\"\nprint shared\n"
	if err := os.WriteFile(entry, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	d := RunFile(entry, &bytes.Buffer{}, ExecuteOptions{})
	if d == nil || d.Code != "S034" {
		t.Fatalf("got %#v", d)
	}
	if d.Pos.Line != 2 {
		t.Fatalf("S034 should point at later import on line 2, got line %d col %d", d.Pos.Line, d.Pos.Column)
	}
}
