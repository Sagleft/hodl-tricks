package main

import "testing"

func TestCheckHostIP(t *testing.T) {
	isLocal, err := checkHostIP("127.0.0.1")
	if err != nil {
		t.Fatal("error: " + err.Error())
	}
	if !isLocal {
		t.Fatal("host should be local")
	}
}
