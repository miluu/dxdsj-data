package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	ColCalendar = bson.M{
		"name": "calendar",
		"index": map[string][]string{
			"0": []string{"title"},
		},
		"background.index": map[string]bool{
			"0": true,
		},
	}
)

type ModelCalendar struct {
	Id            bson.ObjectId `bson:"_id" json:"_id"`
	OriginalId    int64         `bson:"id" json:"id"`
	Timestamp     int64         `bson:"timestamp" json:"timestamp"`
	LocalDateTime string        `bson:"localDateTime" json:"localDateTime"`
	Importance    int           `bson:"importance" json:"importance"`
	Title         string        `bson:"title" json:"title"`
	Forecast      string        `bson:"forecast" json:"forecast"`
	Actual        string        `bson:"actual" json:"actual"`
	Previous      string        `bson:"previous" json:"previous"`
	Revised       string        `bson:"revised" json:"revised"`
	CategoryId    int64         `bson:"category_id" json:"category_id"`
	RelatedAssets string        `bson:"relatedAssets" json:"relatedAssets"`
	Remark        string        `bson:"remark" json:"remark"`
	Mark          string        `bson:"mark" json:"mark"`
	Underline     int           `bson:"underline" json:"underline"`
	AccurateFlag  int           `bson:"accurateFlag" json:"accurateFlag"`
	Level         int           `bson:"level" json:"level"`
	Country       string        `bson:"country" json:"country"`
	Currency      string        `bson:"currency" json:"currency"`
	CalendarType  string        `bson:"calendarType" json:"calendarType"`
	Description   string        `bson:"description" json:"description"`
	Del           bool          `bson:"del" json:"del"`
	// eventRowId string        `bson:"eventRowId" json:"eventRowId"`
}

type CalendarData struct {
	Results []ModelCalendar `bson:"results" json:"results"`
}

func SaveCalendar(m ModelCalendar, col *mgo.Collection) error {
	id := m.OriginalId
	query := bson.M{"id": id}
	if n, err := col.Find(query).Count(); err != nil {
		return err
	} else if n > 0 {
		change := bson.M{
			"timestamp":     m.Timestamp,
			"localDateTime": m.LocalDateTime,
			"importance":    m.Importance,
			"title":         m.Title,
			"forecast":      m.Forecast,
			"actual":        m.Actual,
			"previous":      m.Previous,
			"revised":       m.Revised,
			"category_id":   m.CategoryId,
			"relatedAssets": m.RelatedAssets,
			"remark":        m.Remark,
			"mark":          m.Mark,
			"underline":     m.Underline,
			"accurateFlag":  m.AccurateFlag,
			"level":         m.Level,
			"country":       m.Country,
			"currency":      m.Currency,
			"calendarType":  m.CalendarType,
			"description":   m.Description,
		}
		err := col.Update(query, bson.M{"$set": change})
		return err
	}
	m.Id = bson.NewObjectId()
	err := col.Insert(m)
	return err
}
