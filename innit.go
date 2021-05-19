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
	X      NodeList
	RParen Pos
}

func (x *BasicLit) Pos() Pos { return x.ValuePos }
func (x *BasicLit) End() Pos { return x.ValuePos + Pos(len(x.Value)) }
func (x *Expr) Pos() Pos     { return x.LParen }
func (x *Expr) End() Pos     { return x.RParen + 1 }

type NodeList []Node

func (x NodeList) Pos() Pos {
	if len(x) > 0 {
		return x[0].Pos()
	}
	return NoPos
}

func (x NodeList) End() Pos {
	if n := len(x); n > 0 {
		return x[n-1].End()
	}
	return NoPos
}
