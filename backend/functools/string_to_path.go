package functools

import (
	"errors"
	"strings"
)

func StringToPath(str string, div int)(string, error){
	strLen := len(str)
	if strLen % div != 0 {
		return "", errors.New("string should be divisible by div value")
	}

	sb:= strings.Builder{}

	for i:=0; i < strLen; i+=div {
		sb.WriteString(str[i : i + div])
		sb.WriteString("/")
	}
	return sb.String(), nil
}
