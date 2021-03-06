package util

import (
	"time"
	"code.google.com/p/go-uuid/uuid"
)

func Elapsed( start, duration int64 ) bool {
	now := Now()
	return (now - start) > duration;
}

func Now() int64 {
	now := time.Now()
	nanos := now.UnixNano()
	millis := nanos / 1000000
	return millis
}

func UUID() string {
	return uuid.New();
}
