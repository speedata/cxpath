[![Go Reference](https://pkg.go.dev/badge/github.com/speedata/cxpath.svg)](https://pkg.go.dev/github.com/speedata/cxpath)

# cxpath - a programmer's friendly interface to XPath

With cxpath you can access XML files via XPath 2.0 in a Go friendly matter.

```xml
<a:root xmlns:a="anamespace">
  <a:sub>text</a:sub>
</a:root>
```

```go
package main

import (
	"fmt"
	"log"

	"github.com/speedata/cxpath"
)

func dothings() error {
	ctx, err := cxpath.NewFromFile("myfile.xml")
	if err != nil {
		return err
	}
	// for XPath queries
	ctx.SetNamespace("a", "anamespace")
	root := ctx.Root()
	// prints 'root'
	fmt.Println(root.Eval("local-name()"))
	// prints sub
	fmt.Println(root.Eval("local-name(a:sub)"))
	// prints anamespace
	fmt.Println(root.Eval("namespace-uri(a:sub)"))
	sub := root.Eval("a:sub")
	for cp := range sub.Each("string-to-codepoints(.)") {
		// prints 116, 101, 120, 116 - the codepoints for 'text'
		fmt.Println(cp)
	}
	return nil
}

func main() {
	if err := dothings(); err != nil {
		log.Fatal(err)
	}
}
```


## Error handling

The constructors `NewFromFile` and `NewFromReader` return errors the usual Go way.
For all other methods, errors are stored in the `Error` field of the returned `Context` so that method chaining stays clean:

```go
root := ctx.Root()
result := root.Eval("some/xpath")
if result.Error != nil {
    log.Fatal(result.Error)
}
fmt.Println(result.String())
```

The same applies to the `Each` iterator. If the XPath expression is invalid, the iterator yields a single `Context` with the `Error` field set:

```go
for item := range root.Each("some/xpath") {
    if item.Error != nil {
        log.Fatal(item.Error)
    }
    fmt.Println(item.Int())
}
```

The value accessors `Int()` and `Bool()` also store conversion errors in `Context.Error` rather than returning them directly.

## Installation

    go get github.com/speedata/cxpath

