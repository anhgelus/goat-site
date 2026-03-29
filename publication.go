package site

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"tangled.org/anhgelus.world/xrpc"
	"tangled.org/anhgelus.world/xrpc/atproto"
)

var CollectionPublication = CollectionBase.Name("publication").Build()

// Publication represents a collection of [Document]s published to the web.
// It includes important information about a publication including its location on the web, theming information, user
// [Preferences], and more.
//
// The [Publication] [Record] is not a requirement, but is recommended when publishing collections of related
// [Document]s.
type Publication struct {
	// Base URL of the [Publication].
	// This value will be combined with the [Document.Path] to construct a full URL for the document.
	// Avoid trailing slashes.
	URL *url.URL `json:"-"`
	// Name of the [Publication].
	// Max length: 5000.
	// Max graphemes: 500.
	Name string `json:"name"`
	// Icon to identify the [Publication].
	// Must be a square image and should be at least 256x256.
	Icon *xrpc.Blob `json:"icon,omitempty"`
	// Description of the [Publication].
	// Max length: 30000.
	// Max graphemes: 3000.
	Description *string `json:"description,omitempty"`
	// Simplified theme for tools and apps to utilize when displaying content.
	BasicTheme *Theme `json:"basicTheme,omitempty"`
	// Platform-specific [Preferences] for the [Publication], including discovery and visibility settings.
	Preferences *Preferences `json:"preferences,omitempty"`
}

func (p *Publication) Collection() *atproto.NSID {
	return CollectionPublication
}

func (p *Publication) MarshalMap() (any, error) {
	type t Publication
	pp := struct {
		t
		URL string `json:"url"`
	}{t(*p), strings.TrimSuffix(p.URL.String(), "/")}
	return xrpc.MarshalToMap(pp)
}

func (p *Publication) UnmarshalJSON(b []byte) error {
	type t Publication
	var pp struct {
		t
		URL string `json:"url"`
	}
	err := json.Unmarshal(b, &pp)
	if err != nil {
		return err
	}
	*p = Publication(pp.t)
	p.URL, err = url.Parse(pp.URL)
	if err != nil {
		return err
	}
	p.URL.Path = strings.TrimSuffix(p.URL.Path, "/")
	return nil
}

// Verify the [Publication].
func (p *Publication) Verify(ctx context.Context, client *http.Client, repo *atproto.DID, rkey atproto.RecordKey) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, p.URL.String()+GetPublicationVerificationURI(p.URL.Path), nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	return string(b) == getPublicationVerification(repo, rkey), nil
}

// Preferences of the [Publication].
type Preferences struct {
	// ShowInDiscover decides whether the [Publication] should appear in discovery feeds.
	ShowInDiscover bool `json:"showInDiscover"`
}

// getPublicationVerification returns the string used during the verification of the [Publication].
func getPublicationVerification(repo *atproto.DID, rkey atproto.RecordKey) string {
	return atproto.NewURI(repo, CollectionPublication, rkey).String()
}

// HandlePublicationVerification returns an [http.Handler] used during the verification of the [Publication].
//
// See [GetPublicationVerificationURI].
func HandlePublicationVerification(repo *atproto.DID, rkey atproto.RecordKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, getPublicationVerification(repo, rkey))
	})
}

// GetPublicationVerificationURI returns the URI called during the verification of the [Publication].
//
// path must be empty if the [Publication] is located at the domain root.
// path must start with a slash.
//
// See [HandlePublicationVerification].
func GetPublicationVerificationURI(path string) string {
	if len(path) > 0 && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "/.well-known/" + CollectionPublication.String() + path
}
