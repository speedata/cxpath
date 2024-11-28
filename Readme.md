[![Go Reference](https://pkg.go.dev/badge/github.com/speedata/cxpath.svg)](https://pkg.go.dev/github.com/speedata/cxpath)

# cxpath - a programmer's friendly interface to XPath

**Experimental, not ready for production use**

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


## Installation

    go get github.com/speedata/cxpath


