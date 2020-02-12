package utils

import (
	"testing"
	"time"
)

func TestUtils(t *testing.T) {
	tm := time.Now()
	t.Fatal(FormatTime(tm))
}