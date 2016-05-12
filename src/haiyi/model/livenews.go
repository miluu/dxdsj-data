package model

import (
	"golanger.com/log"
	"haiyi/util"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"strconv"
	"strings"
)

var (
	ColLivenews = bson.M{
		"name": "livenews",
		"index": map[string][]string{
			"0": []string{"title"},
		},
		"background.index": map[string]bool{
			"0": true,
		},
	}
)

/*
livenews表
livenews {
	"id"         : <id>, 			  	//string 获取的id
	"status"     : <status>, 			//string 状态 "published"
	"title"      : <title>, 			//string 标题
	"type"       : <type>, 		  	//string 类型 "news", "data", "ad"
	"importance" : <importance>,	//string 重要 "1", "2", "3"
	"createdAt"  : <createdAt>, 	//string 创建时间 "1462543606"
	"updatedAt"  : <updatedAt>, 	//string 更新时间 "1462543606"
	"contentHtml": <contentHtml>,	//string 输出内容
	"channelSet" : <channelSet>,	//string 所属频道 "1,3"
	"imageUrls"  : [<imageUrls>],	//[]string 图片url
}
*/

type ModelLivenews struct {
	Id                bson.ObjectId `bson:"_id" json:"_id"`
	OriginalId        int64         `bson:"id" json:"id"`
	Title             string        `bson:"title" json:"title"`
	Type              string        `bson:"type" json:"type"`
	Importance        int           `bson:"importance" json:"importance"`
	CreatedAt         int64         `bson:"createdAt" json:"createdAt"`
	UpdatedAt         int64         `bson:"updatedAt" json:"updatedAt"`
	ContentHtml       string        `bson:"contentHtml" json:"contentHtml"`
	ChannelSet        []int         `bson:"channelSet" json:"channelSet"`
	CategorySet       []int         `bson:"categorySet" json:"categorySet"`
	ImageUrls         []string      `bson:"imageUrls" json:"imageUrls"`
	OriginalImageUrls []string      `bson:"orginalImageUrls" json:"orginalImageUrls"`
	Del               bool          `bson:"del" json:"del"`
}

type LivenewsOriginal struct {
	OriginalId  string   `bson:"id" json:"id"`
	Title       string   `bson:"title" json:"title"`
	Type        string   `bson:"type" json:"type"`
	Importance  string   `bson:"importance" json:"importance"`
	CreatedAt   string   `bson:"createdAt" json:"createdAt"`
	UpdatedAt   string   `bson:"updatedAt" json:"updatedAt"`
	ContentHtml string   `bson:"contentHtml" json:"contentHtml"`
	ChannelSet  string   `bson:"channelSet" json:"channelSet"`
	CategorySet string   `bson:"categorySet" json:"categorySet"`
	ImageUrls   []string `bson:"imageUrls" json:"imageUrls"`
}

type LivenewsOriginalData struct {
	Results []LivenewsOriginal `bson:"results" json:"results"`
}

func AddLivenews(ori LivenewsOriginal, col *mgo.Collection, path string) (msg string, err error) {
	m := ModelLivenews{}
	m.Id = bson.NewObjectId()
	if id, err1 := strconv.ParseInt(ori.OriginalId, 10, 0); err1 != nil {
		err = err1
		return
	} else {
		m.OriginalId = id
		if n, errI := col.Find(bson.M{"id": id}).Count(); errI != nil {
			err = errI
			return
		} else if n > 0 {
			msg = "exist"
			return
		}
	}
	if importance, err2 := strconv.Atoi(ori.Importance); err2 != nil {
		err = err2
		return
	} else {
		m.Importance = importance
	}
	if createdAt, err3 := strconv.ParseInt(ori.CreatedAt, 10, 0); err3 != nil {
		err = err3
		return
	} else {
		m.CreatedAt = createdAt
	}
	if updatedAt, err3 := strconv.ParseInt(ori.UpdatedAt, 10, 0); err3 != nil {
		err = err3
		return
	} else {
		m.UpdatedAt = updatedAt
	}
	channelSetStrArr := strings.Split(ori.ChannelSet, ",")
	channelSetIntArr := []int{}
	for _, v := range channelSetStrArr {
		if i, err4 := strconv.Atoi(v); err4 != nil {
			err = err4
			return
		} else {
			channelSetIntArr = append(channelSetIntArr, i)
		}
	}
	categorySetStrArr := strings.Split(ori.CategorySet, ",")
	categorySetIntArr := []int{}
	for _, v := range categorySetStrArr {
		if i, err4 := strconv.Atoi(v); err4 != nil {
			err = err4
			return
		} else {
			categorySetIntArr = append(categorySetIntArr, i)
		}
	}
	imageUrls := []string{}
	for _, imgUrl := range ori.ImageUrls {
		name := bson.NewObjectId().Hex()
		go func() {
			if err := util.SaveImg(imgUrl, name, path); err != nil {
				log.Error("err: ", err)
			} /* else {
				log.Debug("success.")
			}*/
		}()
		imageUrls = append(imageUrls, name)
	}
	m.ChannelSet = channelSetIntArr
	m.CategorySet = categorySetIntArr
	m.Title = util.ClearHtmlTags(ori.Title)
	m.Type = ori.Type
	m.ContentHtml = util.ClearHtmlTags(ori.ContentHtml)
	m.ImageUrls = imageUrls
	m.OriginalImageUrls = ori.ImageUrls
	err = col.Insert(m)
	return
}
