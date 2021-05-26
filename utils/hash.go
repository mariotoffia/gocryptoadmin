package utils

import (
	"fmt"
	"hash/fnv"
	"time"
)

func HashFromString(s string) uint32 {

	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()

}

func HashFromTime(t time.Time) uint32 {

	return HashFromString(t.Format(time.RFC3339))

}

func ToString(v interface{}) string {

	return fmt.Sprintf("%v", v)

}
