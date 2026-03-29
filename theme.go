package site

import (
	"fmt"
	"image/color"
	"math"

	"tangled.org/anhgelus.world/xrpc/atproto"
)

var (
	CollectionTheme          = CollectionBase.SubAuthority("theme")
	CollectionThemeBasic     = CollectionTheme.Name("basic").Build()
	CollectionThemeColor     = CollectionTheme.Name("color")
	CollectionThemeColorRGB  = CollectionThemeColor.Fragment("rgb").Build()
	CollectionThemeColorRGBA = CollectionThemeColor.Fragment("rgba").Build()
)

// Theme ensures [Publication]s maintain their visual identity across different reading applications and platforms by
// defining core colors for content display.
type Theme struct {
	// Background is the color used for content background.
	Background *RGB `json:"background"`
	// Foreground is the color used for content text.
	Foreground *RGB `json:"foreground"`
	// Accent is the color used for links and button backgrounds.
	Accent *RGB `json:"accent"`
	// AccentForeground is the color used for button text.
	AccentForeground *RGB `json:"accentForeground"`
}

func (t *Theme) Collection() *atproto.NSID {
	return CollectionThemeBasic
}

// RGB represents a RGB color.
//
// See also [RGBA].
type RGB struct {
	Red   uint8 `json:"r"`
	Green uint8 `json:"g"`
	Blue  uint8 `json:"b"`
}

func NewRGB(r, g, b uint8) *RGB {
	return &RGB{r, g, b}
}

func (r *RGB) Collection() *atproto.NSID {
	return CollectionThemeColorRGB
}

func (r *RGB) RGBA() *RGBA {
	return &RGBA{*r, 100}
}

func (r *RGB) String() string {
	return fmt.Sprintf("RGB(%d %d %d)", r.Red, r.Green, r.Blue)
}

// RGBA represents a [color.RGBA].
//
// See also [RGB].
type RGBA struct {
	RGB
	// Alpha is the alpha channel where 0 is transparent and 100 is opaque.
	Alpha uint8 `json:"a"`
}

func NewRawRGBA(r, g, b, a uint8) *RGBA {
	if a > 100 {
		panic("invalid alpha: must be <= 100")
	}
	return &RGBA{*NewRGB(r, g, b), a}
}

func NewRGBA(r *color.RGBA) *RGBA {
	red, g, b, a := r.RGBA()
	return &RGBA{
		RGB:   *NewRGB(uint8(red), uint8(g), uint8(b)),
		Alpha: uint8(math.Floor(float64(a)/float64(255)) * 100),
	}
}

func (r *RGBA) Collection() *atproto.NSID {
	return CollectionThemeColorRGBA
}

func (r *RGBA) String() string {
	return fmt.Sprintf("RGBA(%d %d %d %d)", r.Red, r.Green, r.Blue, r.Alpha)
}
