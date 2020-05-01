package di

import (
	"fmt"
	"io"
)

func Greet(writer io.Writer, name string) {
	if _, err := fmt.Fprintf(writer, "Hello, %s", name); err != nil {
	}
}
