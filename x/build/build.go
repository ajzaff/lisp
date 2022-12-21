package builder

import (
	"strconv"

	"github.com/ajzaff/lisp"
)

// Builder provides a convenient type to build Cons.
//
// The zero Builder is usable.
// BeginFrame must be called before using Append* methods.
type Builder struct {
	*builderFrame
}

// Reset clears the Builder to an initial state and drops all frames.
func (b *Builder) Reset() {
	b.builderFrame = new(builderFrame)
}

// BeginFrame adds a new frame to the Builder, used to create nested Cons.
// Call EndFrame to append the Cons to the parent frame.
func (b *Builder) BeginFrame() {
	e := new(builderFrame)
	e.prev = b.builderFrame
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
// If no previous builder frame exists, it's equivalent to calling Reset.
func (b *Builder) EndFrame() {
	prev := b.prev
	head := b.head
	*b.builderFrame = builderFrame{} // GC hint.
	if prev != nil {
		b.builderFrame = prev // Revert frame in Builder.
		prev.appendVal(head)
	}
}

// AppendId appends the Id Lit to the Cons.
func (b *Builder) AppendId(text string) {
	if b.builderFrame == nil {
		b.Reset()
	}
	b.builderFrame.appendVal(lisp.Lit{Token: lisp.Id, Text: text})
}

// AppendNat appends the unsigned integer n to the Cons.
func (b *Builder) AppendNat(n uint64) {
	if b.builderFrame == nil {
		b.Reset()
	}
	b.builderFrame.appendVal(lisp.Lit{Token: lisp.Nat, Text: strconv.FormatUint(n, 10)})
}

// AppendRaw appends an raw text Lit with Token Invalid to the Cons.
func (b *Builder) AppendText(text string) {
	if b.builderFrame == nil {
		b.Reset()
	}
	b.builderFrame.appendVal(lisp.Lit{Text: text})
}

// AppendVal appends a Val v to the Cons.
func (b *Builder) AppendVal(v lisp.Val) {
	if b.builderFrame == nil {
		b.Reset()
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

// precondition: b != nil.
func (b *builderFrame) appendVal(v lisp.Val) {
	if b.head == nil {
		b.head = new(lisp.Cons)
		b.last = b.head
		b.last.Val = v
		return
	}
	b.last.Cons = new(lisp.Cons)
	b.last = b.last.Cons
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
