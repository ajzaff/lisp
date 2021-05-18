package innit

type Node interface {
	Pos() Pos
	End() Pos
}

type BasicLit struct {
	Tok      Token
	ValuePos Pos
	Value    string
}

type Expr struct {
	LParen Pos
	X      []Node
	RParen Pos
}

func (x *BasicLit) Pos() Pos { return x.ValuePos }
func (x *BasicLit) End() Pos { return x.ValuePos + Pos(len(x.Value)) }
func (x *Expr) Pos() Pos     { return x.LParen }
func (x *Expr) End() Pos     { return x.RParen + 1 }
