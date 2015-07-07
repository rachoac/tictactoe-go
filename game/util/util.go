package util

import (
)

func JsonBytesToString(jsonBytes []byte) string {
	if ( jsonBytes == nil ) {
		return ""
	}

	n := len(jsonBytes)

	if ( n < 1 ) {
		return ""
	}

	s := string(jsonBytes[:n])
	return s
}