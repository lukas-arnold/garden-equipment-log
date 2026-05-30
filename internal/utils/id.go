package utils

import "time"

func Id() int64 {
	return time.Now().Unix()
}
