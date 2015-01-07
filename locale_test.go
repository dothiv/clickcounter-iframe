package clickcounteriframe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestThatItParsesLocale(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("en", GetClientLocales("")[0]) // default to en
	locales := GetClientLocales("de-DE,de;q=0.8,en;q=0.6")
	assert.Equal("de-DE", locales[0])
	locales2 := GetClientLocales("en-US,en;q=0.8,de;q=0.6")
	assert.Equal("en-US", locales2[0])
	locales3 := GetClientLocales("en-US;q=0.8,en;q=0.7,de;q=0.9")
	assert.Equal("de", locales3[0])
	assert.Equal("en-US", locales3[1])
	assert.Equal("en", locales3[2])

}
