package lisputil

const (
	IdPattern   = `\p{L}[\p{L}\d]*`
	IntPattern  = `0|[1-9]\d*`
	ExprPattern = `[()]`
	ValPattern  = IdPattern + "|" + IntPattern + "|" + ExprPattern
)
