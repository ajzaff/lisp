package builder

import (
	"strconv"

	"github.com/ajzaff/lisp"
)

// Builder provides a convenient type to build Cons.
//
// The zero Builder is useable.
type Builder struct {
	*builderFrame
}

// Reset clears the Builder to an initial state and drops all frames.
func (b *Builder) Reset() {
	b.builderFrame = nil
	b.BeginFrame()
}

// BeginFrame adds a new frame to the Builder, used to create nested Cons.
// Call EndFrame to append the Cons to the parent frame.
func (b *Builder) BeginFrame() {
	e := new(builderFrame)
	e.Reset(b)
	b.builderFrame = e
}

// DropFrame ends the current frame, if any, while discarding the value.
// If no previous builder frame exists, it's equivalent to calling Reset.
func (b *Builder) DropFrame() {
	if b.builderFrame != nil {
		b.builderFrame = b.builderFrame.prev
	}
}

// EndFrame unrolls the current frame into the previous frame.
// If no previous builder frame exists, EndFrame has no effect.
func (b *Builder) EndFrame() {
	prev := b.prev
	head := b.head
	if prev != nil {
		*b.builderFrame = builderFrame{} // GC hint.
		b.builderFrame = prev            // Revert frame in Builder.
		prev.appendVal(head)
	}
}

// AppendId appends the Id Lit to the Cons.
func (b *Builder) AppendId(text string) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(lisp.Lit{Token: lisp.Id, Text: text})
}

// AppendNat appends the unsigned integer n to the Cons.
func (b *Builder) AppendNat(n uint64) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(lisp.Lit{Token: lisp.Nat, Text: strconv.FormatUint(n, 10)})
}

// AppendRaw appends an raw text Lit with Token Invalid to the Cons.
func (b *Builder) AppendText(text string) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(lisp.Lit{Text: text})
}

// AppendVal appends a Val v to the Cons.
func (b *Builder) AppendVal(v lisp.Val) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(v)
}

// Build constructs and returns the Cons.
// All active frames are consolidated after calling Build.
func (b *Builder) Build() *lisp.Cons {
	if b == nil || b.builderFrame == nil {
		return &lisp.Cons{}
	}
	v := b.build(b)
	if v == nil {
		return &lisp.Cons{}
	}
	return v
}

type builderFrame struct {
	prev *builderFrame
	head *lisp.Cons
	last *lisp.Cons
}

func (b *builderFrame) Reset(e *Builder) {
	b.prev = e.builderFrame
	b.head = new(lisp.Cons)
	b.last = b.head
}

// precondition: b != nil.
// precondition: b.head != nil.
func (b *builderFrame) appendVal(v lisp.Val) {
	if b.last.Val != nil {
		b.last.Cons = new(lisp.Cons)
		b.last = b.last.Cons
	}
	b.last.Val = v
}

// precondition: b != nil.
func (b *builderFrame) build(e *Builder) *lisp.Cons {
	prev := b.prev
	if prev == nil {
		return b.head
	}
	e.EndFrame()
	return prev.build(e)
}
