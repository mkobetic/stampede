package main

import (
	"testing"
)

func TestStartsNewMessage(t *testing.T) {
	fromLines := []string{
		"From joe Mon Jan 2 15:04:05 2006\n",
		"From joe@acme.com Mon Jan 2 15:04:05 2006\n",
		"From - Mon Jan 2 15:04:05 2006\n",
		"From joe Mon Jan 02 15:04:05 2006\n",
		"From joe Mon Jan 12 15:04:05 2006\n",
		"From joe Mon Jan 2 15:04:05 2006 \n",
	}
	for _, line := range fromLines {
		if !startsNewMessage([]byte(line)) {
			t.Fatal("Failed parsing: ", line)
		}
	}
	fromLines = []string{
		"From: Mon Jan 2 15:04:05 2006\n",
		"From Mon Jan 2 15:04:05 2006\n",
		"From - Mon Jan 2 15:04:05 2006 -0700\n",
	}
	for _, line := range fromLines {
		if startsNewMessage([]byte(line)) {
			t.Fatal("Should fail parsing: ", line)
		}
	}
}
