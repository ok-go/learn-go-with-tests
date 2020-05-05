package poker

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := CreateTempFile(t, "12345")
	defer clean()

	tape := &tape{file}
	_, err := tape.Write([]byte("abc"))
	AssertNoError(t, err)

	_, err = file.Seek(0, 0)
	AssertNoError(t, err)
	newFileContents, _ := ioutil.ReadAll(file)

	got := string(newFileContents)
	want := "abc"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
