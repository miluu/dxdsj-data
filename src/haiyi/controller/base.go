package controller

import (
	"golanger.com/db/mongo"
	"haiyi/config"
)

type Base struct {
	Config           config.Conf
	ColLivenews      *mongo.Collection
	ColCalendar      *mongo.Collection
	ColSavedCalendar *mongo.Collection
}
