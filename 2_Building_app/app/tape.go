package poker

import (
	"fmt"
	"os"
)

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	if err := t.file.Truncate(0); err != nil {
		return 0, fmt.Errorf("problem with truncate, %v", err)
	}
	if _, err := t.file.Seek(0, 0); err != nil {
		return 0, fmt.Errorf("problem with seek, %v", err)
	}
	return t.file.Write(p)
}
