package cxpath

import (
	"io"
	"iter"
	"os"

	"github.com/speedata/goxpath"
)

// A Context stores the current state of the XPath parser.
type Context struct {
	P   *goxpath.Parser
	Seq goxpath.Sequence
	Err error
}

// NewFromFile returns a new context from a file name.
func NewFromFile(filename string) (*Context, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return NewFromReader(r)
}

// NewFromReader returns a new context from a reader.
func NewFromReader(r io.Reader) (*Context, error) {
	p, err := goxpath.NewParser(r)
	if err != nil {
		return nil, err
	}
	ctx := Context{
		P: p,
	}
	return &ctx, nil
}
func (ctx *Context) String() string {
	return ctx.Seq.Stringvalue()
}

// SetNamespace sets a prefix/URI pair for XPath queries.
func (ctx *Context) SetNamespace(prefix, uri string) *Context {
	ctx.P.Ctx.Namespaces[prefix] = uri
	return ctx
}

// Root returns the top most element of the XML file.
func (ctx *Context) Root() *Context {
	newContext := Context{
		P: &goxpath.Parser{
			Ctx: goxpath.CopyContext(ctx.P.Ctx),
		},
	}
	newContext.Seq, newContext.Err = newContext.P.Ctx.Root()
	return &newContext
}

// Each can be used as an iterator which returns a Context for each item in the
// resulting sequence.
func (ctx *Context) Each(eval string) iter.Seq[*Context] {
	p := &goxpath.Parser{
		Ctx: goxpath.CopyContext(ctx.P.Ctx),
	}
	seq, _ := p.Evaluate(eval)

	return func(yield func(*Context) bool) {
		for _, itm := range seq {
			newContext := Context{
				P: &goxpath.Parser{
					Ctx: goxpath.CopyContext(ctx.P.Ctx),
				},
				Seq: goxpath.Sequence{itm},
			}
			newContext.P.Ctx.SetContextSequence(newContext.Seq)
			if !yield(&newContext) {
				return
			}
		}
	}
}

// Eval executes the given XPath expression relative to the current context. It
// returns a new Context, so the old one is still available for further use.
func (ctx *Context) Eval(eval string) *Context {
	newContext := Context{
		P: &goxpath.Parser{
			Ctx: goxpath.CopyContext(ctx.P.Ctx),
		},
	}
	newContext.Seq, newContext.Err = newContext.P.Evaluate(eval)
	return &newContext
}
