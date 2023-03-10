package cfg

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestParseMastodonVisibility(t *testing.T) {
	assert.Equal(t, parseMastodonVisibility("public"), "public")
	assert.Equal(t, parseMastodonVisibility("direct"), "direct")
	assert.Equal(t, parseMastodonVisibility("unlisted"), "unlisted")
	assert.Equal(t, parseMastodonVisibility("private"), "private")

	assert.Equal(t, parseMastodonVisibility("Public"), "public")
	assert.Equal(t, parseMastodonVisibility("diRect"), "direct")
	assert.Equal(t, parseMastodonVisibility("unlisTED"), "unlisted")
	assert.Equal(t, parseMastodonVisibility("PRIVATE"), "private")

	assert.Equal(t, parseMastodonVisibility("True"), "unlisted")
	assert.Equal(t, parseMastodonVisibility("eriol"), "unlisted")
	assert.Equal(t, parseMastodonVisibility(""), "unlisted")
	assert.Equal(t, parseMastodonVisibility(" "), "unlisted")
}

func TestParseMastodonMaxCharacters(t *testing.T) {
	assert.Equal(t, parseMastodonMaxCharacters("42"), 42)
	assert.Equal(t, parseMastodonMaxCharacters("-42"), 500)
	assert.Equal(t, parseMastodonMaxCharacters("hello"), 500)
}

func TestParseBoolOrFalse(t *testing.T) {
	assert.Equal(t, parseBoolOrFalse("True"), true)
	assert.Equal(t, parseBoolOrFalse("TRUE"), true)
	assert.Equal(t, parseBoolOrFalse("true"), true)
	assert.Equal(t, parseBoolOrFalse(""), false)
	assert.Equal(t, parseBoolOrFalse("False"), false)
	assert.Equal(t, parseBoolOrFalse("FALSE"), false)
	assert.Equal(t, parseBoolOrFalse("false"), false)
}
