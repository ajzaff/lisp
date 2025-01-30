package builder

import (
	"strconv"

	"github.com/ajzaff/lisp"
)

// Builder provides a convenient type to build a Group.
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

// BeginFrame adds a new frame to the Builder, used to create nested Group.
// Call EndFrame to append the Group to the parent frame.
func (b *Builder) BeginFrame() {
	e := new(builderFrame)
	e.Reset(b)
	e.group = lisp.Group{}
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
	head := b.group
	if prev != nil {
		*b.builderFrame = builderFrame{} // GC hint.
		b.builderFrame = prev            // Revert frame in Builder.
		prev.appendVal(head)
	}
}

// AppendId appends the Id Lit to the Group.
func (b *Builder) AppendId(text string) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(lisp.Lit{Token: lisp.Id, Text: text})
}

// AppendNat appends the unsigned integer n to the Group.
func (b *Builder) AppendNat(n uint64) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(lisp.Lit{Token: lisp.Nat, Text: strconv.FormatUint(n, 10)})
}

// AppendRaw appends an raw text Lit with Token Invalid to the Group.
func (b *Builder) AppendText(text string) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(lisp.Lit{Text: text})
}

// AppendVal appends a Val v to the Group.
func (b *Builder) AppendVal(v lisp.Val) {
	if b.builderFrame == nil {
		b.BeginFrame()
	}
	b.builderFrame.appendVal(v)
}

// Build constructs and returns the Group.
// All active frames are consolidated after calling Build.
func (b *Builder) Build() lisp.Group {
	if b == nil || b.builderFrame == nil {
		return lisp.Group{}
	}
	v := b.build(b)
	if v == nil {
		return lisp.Group{}
	}
	return v
}

type builderFrame struct {
	prev  *builderFrame
	group lisp.Group
}

func (b *builderFrame) Reset(e *Builder) {
	b.prev = e.builderFrame
}

// precondition: b != nil.
// precondition: b.head != nil.
func (b *builderFrame) appendVal(v lisp.Val) {
	b.group = append(b.group, v)
}

// precondition: b != nil.
func (b *builderFrame) build(e *Builder) lisp.Group {
	prev := b.prev
	if prev == nil {
		return b.group
	}
	e.EndFrame()
	return prev.build(e)
}
