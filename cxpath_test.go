package cxpath

import (
	"os"
	"strings"
	"testing"
)

const testXML = `<a:root xmlns:a="anamespace">
  <a:sub>text</a:sub>
  <a:num>42</a:num>
  <a:flag>true</a:flag>
</a:root>`

func newTestContext(t *testing.T) *Context {
	t.Helper()
	ctx, err := NewFromReader(strings.NewReader(testXML))
	if err != nil {
		t.Fatalf("NewFromReader: %v", err)
	}
	ctx.SetNamespace("a", "anamespace")
	return ctx
}

func TestNewFromReader(t *testing.T) {
	ctx, err := NewFromReader(strings.NewReader(testXML))
	if err != nil {
		t.Fatalf("NewFromReader: %v", err)
	}
	if ctx.P == nil {
		t.Fatal("expected parser to be non-nil")
	}
}

func TestNewFromFile(t *testing.T) {
	f, err := os.CreateTemp("", "cxpath-test-*.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(testXML); err != nil {
		t.Fatal(err)
	}
	f.Close()

	ctx, err := NewFromFile(f.Name())
	if err != nil {
		t.Fatalf("NewFromFile: %v", err)
	}
	if ctx.P == nil {
		t.Fatal("expected parser to be non-nil")
	}
}

func TestNewFromFileNotFound(t *testing.T) {
	_, err := NewFromFile("nonexistent.xml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestRootLocalName(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	if root.Error != nil {
		t.Fatalf("Root: %v", root.Error)
	}
	name := root.Eval("local-name()")
	if name.Error != nil {
		t.Fatalf("Eval local-name: %v", name.Error)
	}
	if got := name.String(); got != "root" {
		t.Errorf("local-name() = %q, want %q", got, "root")
	}
}

func TestEvalSubElement(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	sub := root.Eval("a:sub")
	if sub.Error != nil {
		t.Fatalf("Eval a:sub: %v", sub.Error)
	}
	if got := sub.String(); got != "text" {
		t.Errorf("a:sub string = %q, want %q", got, "text")
	}
}

func TestEvalPreservesOriginalContext(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	_ = root.Eval("a:sub")
	// root should still work
	name := root.Eval("local-name()")
	if name.Error != nil {
		t.Fatalf("Eval after previous Eval: %v", name.Error)
	}
	if got := name.String(); got != "root" {
		t.Errorf("original context changed: local-name() = %q, want %q", got, "root")
	}
}

func TestInt(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	num := root.Eval("a:num")
	if num.Error != nil {
		t.Fatalf("Eval a:num: %v", num.Error)
	}
	if got := num.Int(); got != 42 {
		t.Errorf("Int() = %d, want %d", got, 42)
	}
	if num.Error != nil {
		t.Errorf("unexpected error from Int(): %v", num.Error)
	}
}

func TestBool(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	flag := root.Eval("a:flag")
	if flag.Error != nil {
		t.Fatalf("Eval a:flag: %v", flag.Error)
	}
	if got := flag.Bool(); got != true {
		t.Errorf("Bool() = %v, want true", got)
	}
}

func TestBoolFalse(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	// An empty sequence should evaluate to false
	result := root.Eval("a:nonexistent")
	if got := result.Bool(); got != false {
		t.Errorf("Bool() for empty sequence = %v, want false", got)
	}
}

func TestEach(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	sub := root.Eval("a:sub")

	var codepoints []int
	for cp := range sub.Each("string-to-codepoints(.)") {
		if cp.Error != nil {
			t.Fatalf("Each iteration error: %v", cp.Error)
		}
		codepoints = append(codepoints, cp.Int())
	}
	// "text" = 116, 101, 120, 116
	expected := []int{116, 101, 120, 116}
	if len(codepoints) != len(expected) {
		t.Fatalf("got %d codepoints, want %d", len(codepoints), len(expected))
	}
	for i, v := range expected {
		if codepoints[i] != v {
			t.Errorf("codepoint[%d] = %d, want %d", i, codepoints[i], v)
		}
	}
}

func TestEachInvalidXPath(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	for item := range root.Each("///invalid") {
		if item.Error == nil {
			t.Error("expected error for invalid XPath in Each()")
		}
		break
	}
}

func TestSetNamespace(t *testing.T) {
	ctx := newTestContext(t)
	root := ctx.Root()
	ns := root.Eval("namespace-uri(a:sub)")
	if ns.Error != nil {
		t.Fatalf("Eval namespace-uri: %v", ns.Error)
	}
	if got := ns.String(); got != "anamespace" {
		t.Errorf("namespace-uri = %q, want %q", got, "anamespace")
	}
}
