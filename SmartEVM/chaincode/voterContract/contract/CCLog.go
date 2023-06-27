package contract

import (
	"fmt"
)

type CCLog struct {
	PrintDebug bool
}

func (m *CCLog) Debug(msg string, args ...interface{}) {
	if m.PrintDebug {
		fmt.Printf(msg+"\n", args...)
	}
}
