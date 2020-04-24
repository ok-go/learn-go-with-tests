package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Dima")
	want := "Hello, Dima!"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
