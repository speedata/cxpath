package cxpath

import (
	"strings"
	"testing"
)

const benchXML = `<root xmlns:ram="urn:ram">
  <ram:Item><ram:Name>A</ram:Name><ram:ID>1</ram:ID><ram:Qty>10</ram:Qty></ram:Item>
  <ram:Item><ram:Name>B</ram:Name><ram:ID>2</ram:ID><ram:Qty>20</ram:Qty></ram:Item>
  <ram:Item><ram:Name>C</ram:Name><ram:ID>3</ram:ID><ram:Qty>30</ram:Qty></ram:Item>
  <ram:Item><ram:Name>D</ram:Name><ram:ID>4</ram:ID><ram:Qty>40</ram:Qty></ram:Item>
  <ram:Item><ram:Name>E</ram:Name><ram:ID>5</ram:ID><ram:Qty>50</ram:Qty></ram:Item>
  <ram:Item><ram:Name>F</ram:Name><ram:ID>6</ram:ID><ram:Qty>60</ram:Qty></ram:Item>
  <ram:Item><ram:Name>G</ram:Name><ram:ID>7</ram:ID><ram:Qty>70</ram:Qty></ram:Item>
  <ram:Item><ram:Name>H</ram:Name><ram:ID>8</ram:ID><ram:Qty>80</ram:Qty></ram:Item>
  <ram:Item><ram:Name>I</ram:Name><ram:ID>9</ram:ID><ram:Qty>90</ram:Qty></ram:Item>
  <ram:Item><ram:Name>J</ram:Name><ram:ID>10</ram:ID><ram:Qty>100</ram:Qty></ram:Item>
</root>`

func newBenchContext(b *testing.B) *Context {
	b.Helper()
	ctx, err := NewFromReader(strings.NewReader(benchXML))
	if err != nil {
		b.Fatal(err)
	}
	ctx.SetNamespace("ram", "urn:ram")
	return ctx
}

// BenchmarkEachWithEval simulates the typical einvoice pattern:
// iterate items, extract multiple fields per item via Eval().
func BenchmarkEachWithEval(b *testing.B) {
	ctx := newBenchContext(b)
	root := ctx.Root()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for item := range root.Each("ram:Item") {
			_ = item.Eval("ram:Name").String()
			_ = item.Eval("ram:ID").String()
			_ = item.Eval("ram:Qty").String()
		}
	}
}

// BenchmarkEval measures a single Eval call.
func BenchmarkEval(b *testing.B) {
	ctx := newBenchContext(b)
	root := ctx.Root()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root.Eval("ram:Item")
	}
}
