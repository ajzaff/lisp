package innit

type Pos int

type Expr interface {
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

type Closure struct {
	LParen Pos
	Ident  Ident
	X      Expr
	RParen Pos
}
