package lisp

// Regular patterns describing Lisp syntactic elements.
const (
	IdPattern    = `\p{L}+`
	IntPattern   = `0|[1-9]\d*`
	GroupPattern = `[()]`
	ValPattern   = IdPattern + "|" + IntPattern + "|" + GroupPattern
)
