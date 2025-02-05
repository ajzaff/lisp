package nats

import "strconv"

// ParseUint parses s as an unsigned int.
//
// See strconv.ParseInt for a description of bitSize.
func ParseUint(s string, bitSize int) (i uint64, err error) { return strconv.ParseUint(s, 10, bitSize) }
