package models

import (
	"image/color"

	"gopkg.in/mgo.v2/bson"
)

type ColorScheme struct {
	Id                 bson.ObjectId
	Name               string
	BackgroundColor    color.Color
	LinkColor          color.Color
	TextColor          color.Color
	SecondaryTextColor color.Color
}
