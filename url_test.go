package site_test

import (
	"fmt"
	"strings"
	"testing"

	"pgregory.net/rapid"
	site "tangled.org/anhgelus.world/goat-site"
)

func genDid(t *rapid.T, label string) string {
	return "did:plc:" + rapid.StringOfN(rapidLowerRunes, 24, -1, 24).Draw(t, label)
}

func genURL(t *rapid.T, label string) string {
	scheme := "http"
	if rapid.Bool().Draw(t, label+"_secure?") {
		scheme += "s"
	}
	base := rapid.StringOfN(rapidLowerRunes, 1, -1, 63).Draw(t, label+"_base")
	tld := rapid.StringOfN(rapidLowerRunes, 2, -1, 10).Draw(t, label+"_tld")
	sub := rapid.StringOfN(rapidLowerRunes, -1, -1, 32).Draw(t, label+"_sub")
	var sb strings.Builder
	sb.Grow(len(base) + len(tld) + len(sub) + len(scheme) + 5)
	sb.WriteString(scheme)
	sb.WriteString("://")
	if sub != "" {
		sb.WriteString(sub)
		sb.WriteRune('.')
	}
	sb.WriteString(base)
	sb.WriteRune('.')
	sb.WriteString(tld)
	path := genPath(t, label+"_path")
	if path != "/" {
		sb.Grow(len(path))
		sb.WriteString(path)
	}
	return sb.String()
}

func genPath(t *rapid.T, label string) string {
	valid := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	return "/" + rapid.StringOfN(valid, -1, -1, 64).Draw(t, label)
}

func genRecordKey(t *rapid.T, label string) string {
	valid := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyz0123456789"))
	return rapid.StringOfN(valid, 1, -1, 128).Draw(t, label)
}

func genNSID(t *rapid.T, label string) string {
	ascii := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	asciiNums := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	asciiNumsHyp := rapid.RuneFrom([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"))
	authority := func(l string, gen *rapid.Generator[rune]) string {
		a := rapid.StringOfN(ascii, 1, -1, 1).Draw(t, l+"_first")
		if rapid.Bool().Draw(t, l+"_bigger") {
			a += rapid.StringOfN(gen, 0, -1, 61).Draw(t, l+"_second") +
				rapid.StringOfN(asciiNums, 1, -1, 1).Draw(t, l+"_final")
		}
		return strings.ToLower(a)
	}
	var sb strings.Builder
	sb.WriteString(authority(label+"_tld", asciiNumsHyp))
	sb.WriteRune('.')
	sb.WriteString(authority(label+"_base", asciiNumsHyp))
	sb.WriteRune('.')
	ok := true
	for i := 0; ok; i++ {
		sb.WriteString(authority(fmt.Sprintf("%s_sub_%d", label, i), asciiNumsHyp))
		sb.WriteRune('.')
		ok = rapid.Bool().Draw(t, label+"_continue?")
	}
	sb.WriteString(authority(label+"_final", asciiNums))
	return sb.String()
}

func TestParseURL(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		raw := genURL(t, "url")
		url, err := site.ParseURL(raw)
		if err != nil {
			t.Fatal(err)
		}
		if url.IsAT() {
			t.Errorf("invalid url kind: %s is AT", url)
		}
		if url.String() != raw {
			t.Errorf("invalid string: %s, wanted %s", url.String(), raw)
		}
	})
	rapid.Check(t, func(t *rapid.T) {
		did := genDid(t, "did")
		collection := genNSID(t, "collection")
		rkey := genRecordKey(t, "record_key")
		raw := fmt.Sprintf("at://%s/%s/%s", did, collection, rkey)
		url, err := site.ParseURL(raw)
		if err != nil {
			t.Fatal(err)
		}
		if !url.IsAT() {
			t.Errorf("invalid url kind: %s is not AT", url)
		}
		if url.String() != raw {
			t.Errorf("invalid string: %s, wanted %s", url.String(), raw)
		}
	})
}
