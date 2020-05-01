package main

import (
	"learn-go-with-tests/maths/clockface/svg"
	"os"
	"time"
)

func main() {
	svg.Write(os.Stdout, time.Now())
}
