package innit

type Pos int

type Node interface {
	Pos() Pos
	End() Pos
}

type Ident struct {
	NamePos Pos
	Name    string
}

type BasicLit struct {
	Tok      Token
	ValuePos Pos
	Value    string
}

type Expr struct {
	LParen Pos
	Ident  Ident
	X      Node
	RParen Pos
}

func (x *Ident) Pos() Pos    { return x.NamePos }
func (x *Ident) End() Pos    { return x.NamePos + Pos(len(x.Name)) }
func (x *BasicLit) Pos() Pos { return x.ValuePos }
func (x *BasicLit) End() Pos { return x.ValuePos + Pos(len(x.Value)) }
func (x *Expr) Pos() Pos     { return x.LParen }
func (x *Expr) End() Pos     { return x.RParen + 1 }
