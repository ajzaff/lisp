package delim

type Byte struct {
	Val byte
	Ok  bool
}

func (b0 *Byte) Set(v byte) { b0.Val = v; b0.Ok = true }

// Delim returns whether b is a delimiter byte.
func (b0 Byte) Delim() bool { return !b0.Ok || Delim(b0.Val) }

// Between returns whether a delimiter is needed between b0 and b.
func (b0 Byte) Between(b Byte) bool { return !b0.Delim() || !b.Delim() }

// Delim returns whether b is a delimiter byte.
//
// Use Byte.Delim to handle OOB bytes.
func Delim(b byte) bool {
	switch b {
	case '(', ')', ' ', '\t', '\r', '\n':
		return true
	default:
		return false
	}
}
