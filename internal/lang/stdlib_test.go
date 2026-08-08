package lang

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func writeTempVol(t *testing.T, dir, name, source string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func runVolSource(t *testing.T, source string) (string, *Diagnostic) {
	t.Helper()
	dir := t.TempDir()
	path := writeTempVol(t, dir, "main.vol", source)
	var out bytes.Buffer
	d := RunFile(path, &out, ExecuteOptions{Input: strings.NewReader("")})
	return out.String(), d
}

func TestStdReservedCannotRemap(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "vol.config.json"), []byte(`{"name":"t","root":".","paths":{"@std":"fake"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_ = os.MkdirAll(filepath.Join(dir, "fake"), 0o755)
	_ = os.WriteFile(filepath.Join(dir, "fake", "math.vol"), []byte("export hijack\nhijack := 1\n"), 0o600)
	path := writeTempVol(t, dir, "main.vol", `
import "@std/math"
fn main() {
    return math.abs(-3)?
}
print main()
`)
	var out bytes.Buffer
	d := RunFile(path, &out, ExecuteOptions{})
	if d != nil {
		t.Fatalf("unexpected diagnostic: %v", d)
	}
	if out.String() != "3\n" {
		t.Fatalf("got %q", out.String())
	}
}

func TestDictRuntime(t *testing.T) {
	out, d := runVolSource(t, `
d := dict()
d["a"] = 1
d["b"] = 2
print d.len
print d["a"]
print d.keys()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	if out != "2\n1\n[a, b]\n" {
		t.Fatalf("got %q", out)
	}

	out, d = runVolSource(t, `
d := dict("a", 1, "b", 2)
print d.len
print d["a"]
print d.keys()
`)
	if d != nil {
		t.Fatalf("pairs diagnostic: %+v", d)
	}
	if out != "2\n1\n[a, b]\n" {
		t.Fatalf("pairs got %q", out)
	}

	_, d = runVolSource(t, `print dict("only")`)
	if d == nil || d.Code != "R018" {
		t.Fatalf("odd arity diagnostic = %#v", d)
	}
	if d.Fix == "" {
		t.Fatalf("odd arity Fix empty")
	}
	_, d = runVolSource(t, `print dict(1, "x")`)
	if d == nil || d.Code != "R045" {
		t.Fatalf("non-string key diagnostic = %#v", d)
	}
}

func TestStdMathAndStrings(t *testing.T) {
	out, d := runVolSource(t, `
import "@std/math"
import "@std/strings"
fn show() {
    print math.abs(-7)?
    print strings.trim("  hi  ")
    print strings.join(["a", "b"], "-")
    print strings.has("hello", "ell")
    print strings.prefix("hello", "he")
    print strings.suffix("hello", "lo")
    print strings.replace("a-a", "-", "+")
}
show()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	want := "7\nhi\na-b\ntrue\ntrue\ntrue\na+a\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestStdPath(t *testing.T) {
	out, d := runVolSource(t, `
import "@std/path"
print path.join("x", "y", "z")
print path.base("/a/b.txt")
print path.dir("/a/b.txt")
print path.ext("/a/b.txt")
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	want := "x/y/z\nb.txt\n/a\n.txt\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestStdFSEnvTime(t *testing.T) {
	dir := t.TempDir()
	file := filepath.ToSlash(filepath.Join(dir, "note.txt"))
	t.Setenv("VOL_SF3_TEST", "from-env")
	out, d := runVolSource(t, `
import "@std/fs"
import "@std/env"
import "@std/time"
fn demo() {
    path := "`+file+`"
    fs.write(path, "hello")?
    print fs.read(path)?
    print fs.exists(path)
    print env.get("VOL_SF3_TEST") ?? "missing"
    print env.get("VOL_SF3_MISSING") ?? "fallback"
    time.now()
    print time.format(0, "2006")
}
demo()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	if out != "hello\ntrue\nfrom-env\nfallback\n1970\n" {
		t.Fatalf("got %q", out)
	}
}

func TestStdURL(t *testing.T) {
	out, d := runVolSource(t, `
import "@std/url"
fn demo() {
    u := url.parse("https://example.com:8443/x?q=1")?
    print u.scheme
    print u.host
    print u.port
    print u.path
    print u.query["q"]
}
demo()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	if out != "https\nexample.com\n8443\n/x\n1\n" {
		t.Fatalf("got %q", out)
	}
}

func TestStdJSON(t *testing.T) {
	out, d := runVolSource(t, `
import "@std/json"
import "@std/strings"
fn demo() {
    v := json.parse("{\"a\":1,\"b\":null}")?
    print v["a"]
    print v["b"]
    print strings.has(json.dump(v)?, "\"a\":1")
}
demo()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	if out != "1\nnone\ntrue\n" {
		t.Fatalf("got %q", out)
	}
}

func TestStdYAML(t *testing.T) {
	out, d := runVolSource(t, `
import "@std/yaml"
fn demo() {
    y := yaml.parse("a: 2\nb: null\n")?
    print y["a"]
    print y["b"]
}
demo()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	if out != "2\nnone\n" {
		t.Fatalf("got %q", out)
	}
}

func TestStdProcessAndDB(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.ToSlash(filepath.Join(dir, "t.db"))
	out, d := runVolSource(t, `
import "@std/strings"
import "@std/process"
import "@std/db"
fn demo() {
    p := process.run(["echo", "hi"])?
    print p.status
    print strings.trim(p.stdout)
    conn := db.open("`+dbPath+`")?
    db.exec(conn, "create table hits (n int)")?
    db.exec(conn, "insert into hits(n) values (3)")?
    rows := db.query(conn, "select n from hits")?
    print rows[0]["n"]
    db.close(conn)?
}
demo()
`)
	if d != nil {
		t.Fatalf("diagnostic: %+v", d)
	}
	if out != "0\nhi\n3\n" {
		t.Fatalf("got %q", out)
	}
}

func TestStdHTTPFetchAndListen(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
	srv := &http.Server{Handler: mux}
	go func() { _ = srv.Serve(ln) }()
	defer func() { _ = srv.Close() }()
	addr := ln.Addr().String()

	out, d := runVolSource(t, `
import "@std/http"
fn demo() {
    r := http.fetch("http://`+addr+`/ping")?
    print r.status
    print r.body
}
demo()
`)
	if d != nil {
		t.Fatalf("fetch diagnostic: %+v", d)
	}
	if out != "200\npong\n" {
		t.Fatalf("fetch got %q", out)
	}

	ln2, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	listenAddr := ln2.Addr().String()
	_ = ln2.Close()

	dir := t.TempDir()
	path := writeTempVol(t, dir, "server.vol", `
import "@std/http"
fn handle(req) {
  return http.reply(200, dict())
}
http.listen("`+listenAddr+`", handle)
`)
	go func() { _ = RunFile(path, io.Discard, ExecuteOptions{}) }()
	deadline := time.Now().Add(3 * time.Second)
	var resp *http.Response
	for time.Now().Before(deadline) {
		resp, err = http.Get("http://" + listenAddr + "/")
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		t.Fatalf("listen not ready: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("content-type %q", ct)
	}
}

func TestUnknownStdModule(t *testing.T) {
	_, d := runVolSource(t, `import "@std/nope"`)
	if d == nil || d.Code != "S031" {
		t.Fatalf("want S031, got %+v", d)
	}
}

func TestStdJoinNamespacesNoCollision(t *testing.T) {
	out, d := runVolSource(t, `
import "@std/strings"
import "@std/path"
print strings.join(["a", "b"], "-")
print path.join("a", "b")
`)
	if d != nil {
		t.Fatalf("unexpected diagnostic: %+v", d)
	}
	want := "a-b\na/b\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestStdModuleNameCollision(t *testing.T) {
	// User modules with the same basename collide on the binding name.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "vol.config.json"), []byte(`{"root":"."}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_ = os.MkdirAll(filepath.Join(dir, "a"), 0o755)
	_ = os.MkdirAll(filepath.Join(dir, "b"), 0o755)
	writeTempVol(t, filepath.Join(dir, "a"), "util.vol", "export value\nvalue := 1\n")
	writeTempVol(t, filepath.Join(dir, "b"), "util.vol", "export value\nvalue := 2\n")
	path := writeTempVol(t, dir, "main.vol", `
import "a/util"
import "b/util"
print util.value
`)
	d := RunFile(path, io.Discard, ExecuteOptions{})
	if d == nil || d.Code != "S034" {
		t.Fatalf("want S034 module binding collision, got %+v", d)
	}
}
