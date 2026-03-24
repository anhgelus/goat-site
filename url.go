package site

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"github.com/bluesky-social/indigo/atproto/syntax"
)

var (
	ErrIncompleteURL = errors.New("incomplete url")
)

// ATURL represents a full AT [url.URL].
type ATURL struct {
	syntax.ATURI
}

func (at *ATURL) String() string {
	// not using [fmt.Sprintf] because it may be slower
	return "at://" + at.Authority().String() +
		"/" + at.Collection().String() +
		"/" + at.RecordKey().String()
}

// ParseATURL returns an [ATURL] from a raw string.
func ParseATURL(raw string) (*ATURL, error) {
	u, err := syntax.ParseATURI(raw)
	if err != nil {
		return nil, err
	}
	if u.Collection() == "" {
		return nil, ErrIncompleteURL
	}
	if u.RecordKey() == "" {
		return nil, ErrIncompleteURL
	}
	return &ATURL{u}, nil
}

// URL represents an [url.URL] that may be an [ATURL].
type URL struct {
	url *url.URL
	at  *ATURL
}

// IsAT returns true if the [URL] is an [ATURL].
//
// See [URL.AT] and [URL.URL].
func (u *URL) IsAT() bool {
	return u.at != nil
}

// URL returns the [url.URL].
// Panics if it is an [ATURL].
//
// See [URL.IsAT].
func (u *URL) URL() *url.URL {
	if u.IsAT() {
		panic("not an URL")
	}
	return u.url
}

// AT returns the [ATURL].
// Panics if it isn't an [ATURL].
//
// See [URL.IsAT].
func (u *URL) AT() *ATURL {
	if !u.IsAT() {
		panic("not an AT URL")
	}
	return u.at
}

func (u *URL) String() string {
	if u.IsAT() {
		return u.at.String()
	}
	return u.url.String()
}

func (u *URL) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	url, err := ParseURL(s)
	if err != nil {
		return err
	}
	*u = *url
	return nil
}

// ParseURL returns an [URL] from a raw string.
func ParseURL(raw string) (*URL, error) {
	if strings.HasPrefix(raw, "at://") {
		u, err := ParseATURL(raw)
		if err != nil {
			return nil, err
		}
		return &URL{at: u}, nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimPrefix(u.Path, "/")
	return &URL{url: u}, nil
}
