package site

const LexiconPublication = LexiconBase + ".publication"

// Publication represents a collection of documents published to the web.
// It includes important information about a publication including its location on the web, theming information, user
// preferences, and more.
//
// The publication lexicon is not a requirement, but is recommended when publishing collections of related documents.
type Publication struct {
	// Base URL of the [Publication].
	// This value will be combined with the [Document.Path] to construct a full URL for the document.
	// Avoid trailing slashes.
	URL string `json:"string"`
	// Name of the [Publication].
	// Max length: 5000.
	// Max graphemes: 500.
	Name string `json:"name"`
	// Icon to identify the [Publication].
	// Must be a square image and should be at least 256x256.
	Icon any `json:"-"`
	// Description of the [Publication].
	// Max length: 30000.
	// Max graphemes: 3000.
	Description *string `json:"description,omitempty"`
	// Simplified theme for tools and apps to utilize when displaying content.
	// Ref to `site.standard.theme.basic`.
	BasicTheme any `json:"basicTheme,omitempty"`
	// Platform-specific [Preferences] for the [Publication], including discovery and visibility settings.
	Preferences *Preferences `json:"preferences,omitempty"`
}

// Preferences of the [Publication].
type Preferences struct {
	// ShowInDiscover decides whether the [Publication] should appear in discovery feeds.
	ShowInDiscover bool `json:"showInDiscover"`
}

func (p *Publication) Type() string {
	return LexiconPublication
}
