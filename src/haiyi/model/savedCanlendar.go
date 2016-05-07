package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

var (
	ColSavedCalendar = bson.M{
		"name": "saved-calendar",
		"index": map[string][]string{
			"0": []string{"year", "month"},
		},
		"background.index": map[string]bool{
			"0": true,
		},
	}
)

type ModelSavedCalendar struct {
	Id    bson.ObjectId `bson:"_id" json:"_id"`
	Year  int           `bson:"year" json:"year"`
	Month int           `bson:"month" json:"month"`
}

func AddSavedCalendar(year, month int, col *mgo.Collection) error {
	m := ModelSavedCalendar{}
	m.Id = bson.NewObjectId()
	m.Year = year
	m.Month = month
	err := col.Insert(m)
	return err
}

func GetLastSavedCalendar(col *mgo.Collection) (year int, month int, err error) {
	mLast := ModelSavedCalendar{}
	find := col.Find(bson.M{})
	now := time.Now()
	yearNow := int(now.Year())
	monthNow := int(now.Month())
	n := 0
	if n, err = find.Count(); err != nil {
		return
	} else if n == 0 {
		year = yearNow
		month = monthNow
		return
	}
	if err = find.Sort([]string{"year", "month"}...).One(&mLast); err != nil {
		return
	}
	year = mLast.Year
	month = mLast.Month
	return
}
