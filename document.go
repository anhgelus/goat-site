package site

import "time"

const LexiconDocument = LexiconBase + ".document"

// Document may be standalone or associated with a [Publication].
// This lexicon can be used to store a document's content and its associated metadata.
type Document struct {
	// Site points to a [Publication] record `at://` or a [Publication.URL] `https://` for loose documents.
	// Avoid trailing slashes.
	Site string `json:"site"`
	// Title of the [Document].
	// Max length: 5000.
	// Max graphemes: 500.
	Title string `json:"title"`
	// PublishedAt is the [time.Time] of the [Document]'s publish time.
	PublishedAt time.Time `json:"-"`
	// Path is combined with [Document.Site] or [Publication.URL] to construct a canonical URL to the document.
	// A slash should be included at the beginning of this value.
	Path *string `json:"path,omitempty"`
	// A brief Description or excerpt from the [Document].
	// Max length: 30000.
	// Max graphemes: 3000.
	Description *string `json:"description,omitempty"`
	// CoverImage to used for thumbnail or cover.
	// Less than 1MB is size.
	CoverImage any `json:"-"`
	// Content is an open union used to define the [Document]'s content.
	// Each entry must specify a `$type`.
	Content []any `json:"-"`
	// TextContent is a plaintext representation of the [Document.Content].
	// Should not contain markdown or other formatting.
	TextContent string `json:"textContent,omitempty"`
	// BlueskyPostRef is a strong reference to a Bluesky post.
	// Useful to keep track of comments off-platform.
	BlueskyPostRef any `json:"bskyPostRef,omitempty"`
	// Tags is an array of strings used to tag or categorize the [Document].
	// Avoid prepending tags with hashtags.
	// Max length: 1280.
	// Max graphemes: 128.
	Tags []string `json:"tags,omitempty"`
	// UpdatedAt is the [time.Time] of the [Document]'s last edit.
	UpdatedAt *time.Time `json:"-"`
}

func (d *Document) Type() string {
	return LexiconDocument
}
