package cast

import (
	"testing"
	"time"
)

func TestToTimeUnixNumericString(t *testing.T) {
	const sec int64 = 1609459200 // 2021-01-01T00:00:00Z
	want := time.Unix(sec, 0).UTC()

	got, err := ToTimeE("1609459200")
	if err != nil {
		t.Fatal(err)
	}
	if !got.UTC().Equal(want) {
		t.Fatalf("string unix: got %v want %v", got.UTC(), want)
	}

	got, err = ToTimeE(float64(sec) + 0.5)
	if err != nil {
		t.Fatal(err)
	}
	if got.Unix() != sec {
		t.Fatalf("float unix sec: got %d", got.Unix())
	}
	// named formats still work
	got, err = ToTimeE("2021-01-01T00:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	if got.UTC().Format(time.RFC3339) != "2021-01-01T00:00:00Z" {
		t.Fatalf("rfc3339: got %v", got)
	}
}
