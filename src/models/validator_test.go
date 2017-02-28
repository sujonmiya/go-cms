package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

const (
	s = "asdfgfdsgfdgd"
)

func TestIsEmpty(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsEmpty(""))
	assert.False(IsEmpty(s))
}

func TestLengthNotBetween(t *testing.T) {
	assert := assert.New(t)

	assert.False(LengthNotBetween(s, uint8(3), uint8(len(s))))
	assert.True(LengthNotBetween(s, uint8(3), uint8(5)))
}

func TestLengthExceedsMax(t *testing.T) {
	assert := assert.New(t)

	assert.False(LengthExceedsMax(s, uint16(len(s))))
	assert.True(LengthExceedsMax(s, 5))
}

func TestIsValidObjectId(t *testing.T) {
	assert := assert.New(t)

	assert.True(IsValidObjectId(bson.NewObjectId().Hex()))
	assert.False(IsValidObjectId(s))
}

func TestContainsNotAllowed(t *testing.T) {
	assert := assert.New(t)

	assert.False(ContainsNotAllowed(s, PatternAlphaSpace))
	assert.True(ContainsNotAllowed("sffHGHJ Hhj3532 dfs88", PatternAlphaSpace))
	assert.False(ContainsNotAllowed("sffHGHJ Hhj3532 dfs88", PatternAlphaSpaceNumeric))
	assert.True(ContainsNotAllowed("dsfsdf 43543, dsf: (df) - dfds! dsfds? fdgd.", PatternAlphaSpace))
	assert.False(ContainsNotAllowed("dsfsdf 43543, dsf: (df) - dfds! dsfds? fdgd.", PatternAlphaNumericPun))
}
