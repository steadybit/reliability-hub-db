package icon

import (
	"encoding/base64"
	"testing"
)

func TestDecode_Base64(t *testing.T) {
	svg := `<svg viewBox="0 0 24 24"><path d="M1 1"/></svg>`
	in := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
	got, err := Decode(in)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got != svg {
		t.Errorf("got %q, want %q", got, svg)
	}
}

func TestDecode_UTF8URLEncoded(t *testing.T) {
	svg := `<svg><path d="M1 1"/></svg>`
	in := "data:image/svg+xml;utf8,%3Csvg%3E%3Cpath%20d%3D%22M1%201%22%2F%3E%3C%2Fsvg%3E"
	got, err := Decode(in)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got != svg {
		t.Errorf("got %q, want %q", got, svg)
	}
}

func TestDecode_RawSVG(t *testing.T) {
	svg := `<svg><path d="M1 1"/></svg>`
	got, err := Decode(svg)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got != svg {
		t.Errorf("got %q, want %q", got, svg)
	}
}

func TestDecode_TrimsWhitespace(t *testing.T) {
	svg := `<svg/>`
	got, _ := Decode("\n  " + svg + "  \n")
	if got != svg {
		t.Errorf("got %q, want %q", got, svg)
	}
}

func TestDecode_Empty(t *testing.T) {
	got, _ := Decode("")
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
