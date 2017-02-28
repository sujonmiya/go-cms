package models

import (
	"strings"

	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
)

const (
	PatternAlpha             string = "^[a-zA-Z]+$"
	PatternAlphaSpace        string = "^[a-zA-Z ]+$"
	PatternAlphaNumeric      string = "^[a-zA-Z0-9]+$"
	PatternAlphaSpaceNumeric string = "^[a-zA-Z0-9 ]+$"
	PatternAlphaNumericPun   string = "^[a-zA-Z0-9,:!?)( .-]+$"
)

func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func IsNotEmail(s string) bool {
	return !govalidator.IsEmail(s)
}

func LengthNotBetween(s string, minimum, maximum uint8) bool {
	length := uint8(len(s))
	return length < minimum || length > maximum
}

func LengthExceedsMax(s string, maximum uint16) bool {
	return len(s) > int(maximum)
}

func IsValidObjectId(ids ...string) bool {
	for _, id := range ids {
		if !bson.IsObjectIdHex(id) {
			return false
		}
	}

	return true
}

func ContainsNotAllowed(s, pattern string) bool {
	return !govalidator.Matches(s, pattern)
}
