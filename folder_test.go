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
		"From joe Mon Jan 2 3:04:05 2006 \n",
		"From - Wed Feb 21 16:15:07 2007\r\n",
		"From MARS@cincom.com Tue Sep 16 10:32 2008\n",
		"From - Tue May  8 16:48:12 2007\n",
		"From - Tue May  8  4:48:12 2007\n",
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
