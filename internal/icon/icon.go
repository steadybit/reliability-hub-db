// Package icon normalizes the various icon encodings extensions emit
// (data:image/svg+xml;base64,…, data:image/svg+xml;utf8,…, or raw SVG) into a
// single canonical form: the raw SVG markup with no surrounding whitespace.
package icon

import (
	"encoding/base64"
	"net/url"
	"strings"
)

// Decode normalizes an icon value as emitted by the extension --describe JSON.
// Supported inputs:
//   - data:image/svg+xml;base64,<base64-svg>
//   - data:image/svg+xml;utf8,<url-encoded-svg>
//   - data:image/svg+xml,<url-encoded-svg>
//   - raw SVG markup (anything starting with "<svg" or "<?xml")
func Decode(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", nil
	}
	if !strings.HasPrefix(s, "data:image/svg+xml") {
		return s, nil
	}
	// Strip "data:image/svg+xml" and one optional ";<param>".
	rest := strings.TrimPrefix(s, "data:image/svg+xml")
	param := ""
	if strings.HasPrefix(rest, ";") {
		commaIdx := strings.Index(rest, ",")
		if commaIdx < 0 {
			return s, nil // malformed; leave alone
		}
		param = rest[1:commaIdx]
		rest = rest[commaIdx:]
	}
	if !strings.HasPrefix(rest, ",") {
		return s, nil
	}
	payload := rest[1:]
	switch param {
	case "base64":
		b, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	case "utf8", "":
		dec, err := url.QueryUnescape(payload)
		if err != nil {
			// Some encoders leave the SVG unescaped after ";utf8,". Fall back.
			return strings.TrimSpace(payload), nil
		}
		return strings.TrimSpace(dec), nil
	default:
		return s, nil
	}
}
