package utils

import (
	"sync/atomic"
	"time"
)

var lastID int64

func Id() int64 {

	now :=
		time.Now().UnixNano()

	for {

		old :=
			atomic.LoadInt64(
				&lastID,
			)

		if now <= old {
			now = old + 1
		}

		if atomic.CompareAndSwapInt64(
			&lastID,
			old,
			now,
		) {
			return now
		}
	}
}
