package functools

import (
	"strings"
)

func PathFromIdGenerator(id string) string {
	sb := strings.Builder{}
	idLen := len(id)
	twoLetterSize := idLen >> 1
	for i := 0; i < twoLetterSize; i++ {
		sb.WriteByte('/')
		sb.WriteByte(id[idLen-2*i-1])
		sb.WriteByte('/')
		sb.WriteByte(id[idLen-2*i-2])
	}

	if idLen != twoLetterSize<<1 {
		sb.WriteByte('/')
		sb.WriteByte(id[0])
	}
	return sb.String()
}
