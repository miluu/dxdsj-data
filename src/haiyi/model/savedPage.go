package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	ColSavedPage = bson.M{
		"name": "saved-page",
		"index": map[string][]string{
			"0": []string{"page"},
		},
		"background.index": map[string]bool{
			"0": true,
		},
	}
)

type ModelSavedPage struct {
	Id   bson.ObjectId `bson:"_id" json:"_id"`
	Page int           `bson:"page" json:"page"`
}

func AddSavedPage(page int, col *mgo.Collection) error {
	m := ModelSavedPage{}
	m.Id = bson.NewObjectId()
	m.Page = page
	err := col.Insert(m)
	return err
}

func GetLastSavedPage(col *mgo.Collection) (p int, err error) {
	mLast := ModelSavedPage{}
	find := col.Find(bson.M{})
	n := 0
	if n, err = find.Count(); err != nil {
		return
	} else if n == 0 {
		p = n
		return
	}
	if err = find.Sort("-page").One(&mLast); err != nil {
		return
	}
	p = mLast.Page
	return
}
