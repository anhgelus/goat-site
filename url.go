package site

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"

	"tangled.org/anhgelus.world/xrpc"
	"tangled.org/anhgelus.world/xrpc/atproto"
)

var (
	ErrIncompleteURL = errors.New("incomplete url")
)

// URL represents an [url.URL] that may be an [ATURL].
type URL struct {
	url *url.URL
	at  *atproto.RawURI
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
func (u *URL) AT() *atproto.RawURI {
	if !u.IsAT() {
		panic("not an AT URL")
	}
	return u.at
}

func (u *URL) MarshalMap() (any, error) {
	if u.IsAT() {
		return xrpc.MarshalToMap(u.AT())
	}
	return xrpc.MarshalToMap(u.URL().String())
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
		u, err := atproto.ParseRawURI(raw)
		if err != nil {
			return nil, err
		}
		return &URL{at: &u}, nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, err
	}
	u.Path = strings.TrimPrefix(u.Path, "/")
	return &URL{url: u}, nil
}
