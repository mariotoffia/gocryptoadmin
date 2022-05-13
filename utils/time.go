package utils

import "time"

var Epoch = time.Unix(0, 0)

func ToUnixMillisFromTimeStamp(ts int64) time.Time {

	return Epoch.Add(time.Duration(ts) * time.Millisecond)

}
