package model

import (
	// "golanger.com/log"
	"errors"
	"haiyi/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strings"
	// "strconv"
)

var (
	ColLivenews = bson.M{
		"name": "livenews",
		"index": map[string][]string{
			"0": []string{"content"},
		},
		"background.index": map[string]bool{
			"0": true,
		},
	}
)

type ModelLivenews struct {
	Id          bson.ObjectId `bson:"_id" json:"_id"`
	Type        int           `bson:"type" json:"type"`                 //类型 0.新闻 1.数据
	Importance  int           `bson:"importance" json:"importance"`     //重要性 0.重要 1.普通
	PublishTime string        `bson:"publish_time" json:"publish_time"` //发布时间 (2016-05-14 02:00:52)
	Content     string        `bson:"content" json:"content"`           //内容 html/text
	Img         string        `bson:"img" json:"img"`                   //图片

	Prefix    string `bson:"prefix" json:"prefix"`       //前值
	Predicted string `bson:"predicted" json:"predicted"` //预期
	Actual    string `bson:"actual" json:"actual"`       //实际
	Star      int    `bson:"star" json:"star"`           //Star 1-5
	Effect    int    `bson:"effect" json:"effect"`       //影响 0.无影响 1.利多 2.利空
	Country   string `bson:"country" json:"country"`     //国家
	Time      string `bson:"time" json:"time"`           //时间 (01:00)

	OriginalId      string `bson:"original_id" json:"original_id"`           //原始ID
	OriginalContent string `bson:"original_content" json:"original_content"` //原始内容
	Del             bool   `bson:"del" json:"del"`
}

func AddLivenews(m ModelLivenews, col *mgo.Collection, saveImg bool, saveImgPath string) (msg string, err error) {
	if m.OriginalId == "" {
		err = errors.New("No originalId")
		return
	}
	query := bson.M{
		"original_id": m.OriginalId,
	}
	if n, errF := col.Find(query).Count(); errF != nil {
		err = errF
		return
	} else if n > 0 {
		msg = "exist"
		return
	}
	if saveImg && m.Img != "" {
		imgUrl := strings.Replace(m.Img, "_lite", "", -1)
		imgUrl = "http://image.jin10x.com/" + imgUrl
		imgName := bson.NewObjectId().Hex()
		if saveImg {
			go func() {
				util.SaveImg(imgUrl, imgName, saveImgPath)
			}()
		}
		m.Img = imgName
	}

	m.Id = bson.NewObjectId()
	err = col.Insert(m)
	return
}

func GetLivenewsCount(col *mgo.Collection) (count int, err error) {
	query := bson.M{
		"del": false,
	}
	count, err = col.Find(query).Count()
	return
}

func GetOldestLivenewsOriginalId(col *mgo.Collection) (originalId string, err error) {
	query := bson.M{
		"del": false,
	}
	sorter := "original_id"
	oldestLivenews := ModelLivenews{}
	if err = col.Find(query).Sort(sorter).One(&oldestLivenews); err != nil {
		return
	}
	originalId = oldestLivenews.OriginalId
	return
}
