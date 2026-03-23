package site

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bluesky-social/indigo/atproto/syntax"
	lexutil "github.com/bluesky-social/indigo/lex/util"
)

const CollectionPublication = CollectionBase + ".publication"

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
	Icon *Blob `json:"icon,omitempty"`
	// Description of the [Publication].
	// Max length: 30000.
	// Max graphemes: 3000.
	Description *string `json:"description,omitempty"`
	// Simplified theme for tools and apps to utilize when displaying content.
	BasicTheme *Theme `json:"basicTheme,omitempty"`
	// Platform-specific [Preferences] for the [Publication], including discovery and visibility settings.
	Preferences *Preferences `json:"preferences,omitempty"`
}

func (p *Publication) Type() string {
	return CollectionPublication
}

func (p *Publication) MarshalMap() (map[string]any, error) {
	type t Publication
	pp := struct {
		t
		URL string `json:"url"`
	}{t(*p), strings.TrimSuffix(p.URL.String(), "/")}
	return MarshalToMap(pp)
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
func (p *Publication) Verify(ctx context.Context, client *http.Client, repo syntax.AtIdentifier, rkey syntax.RecordKey) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, p.URL.String()+GetPublicationVerificationURI(p.URL.Path), nil)
	if err != nil {
		return false, err
	}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return false, err
	}
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

// GetPublication returns the [Publication] in the repo associated with the rkey.
// Automatically uses the latest CID.
func GetPublication(ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey syntax.RecordKey) (*Publication, error) {
	return get[*Publication](ctx, client, CollectionPublication, repo, rkey)
}

// ListPublications returns all the [Publication]s stored in the repo and the cursor.
//
// See [MaxItemsPerList].
func ListPublications(ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, cursor string, reverse bool) ([]*Publication, *string, error) {
	return listRecord[*Publication](ctx, client, CollectionPublication, repo, cursor, reverse)
}

// CreatePublication in a repo with the given rkey.
// Always tries to validate the [Publication] against the lexicon saved.
//
// Rkey can be nil.
func CreatePublication(ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey *syntax.RecordKey, pub *Publication) (*Result, error) {
	return createRecord(ctx, client, CollectionPublication, repo, rkey, pub)
}

// UpdatePublication in a repo with the given rkey.
// Always tries to validate the [Publication] against the lexicon saved.
func UpdatePublication(ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey syntax.RecordKey, pub *Publication) (*Result, error) {
	return updateRecord(ctx, client, CollectionPublication, repo, rkey, pub)
}

// DeletePublication in a repo with the given rkey.
func DeletePublication(ctx context.Context, client lexutil.LexClient, repo syntax.AtIdentifier, rkey syntax.RecordKey) error {
	return deleteRecord(ctx, client, CollectionPublication, repo, rkey)
}

// getPublicationVerification returns the string used during the verification of the [Publication].
func getPublicationVerification(repo syntax.AtIdentifier, rkey syntax.RecordKey) string {
	return createAtURI(repo, CollectionPublication, rkey)
}

// HandlePublicationVerification returns an [http.Handler] used during the verification of the [Publication].
//
// See [GetPublicationVerificationURI].
func HandlePublicationVerification(repo syntax.AtIdentifier, rkey syntax.RecordKey) http.Handler {
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
	return "/.well-known/" + CollectionPublication + path
}
