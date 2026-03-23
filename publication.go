package site

import (
	"context"
	"fmt"
	"net/http"
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
	URL string `json:"url"`
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
	pp := t(*p)
	pp.URL = strings.TrimSuffix(pp.URL, "/")
	return MarshalToMap(pp)
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

// HandlePublicationVerification returns an [http.Handler] used during the verification of the [Publication].
//
// See [GetPublicationVerificationURI].
func HandlePublicationVerification(repo syntax.AtIdentifier, rkey syntax.RecordKey) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "at://%s/%s/%s", repo, CollectionPublication, rkey)
	})
}

// GetPublicationVerificationURI returns the URI called during the verification of the [Publication].
// Path must be empty if the [Publication] is located at the domain root.
//
// See [HandlePublicationVerification].
func GetPublicationVerificationURI(path string) string {
	return "/.well-known/" + CollectionPublication + path
}
