package controller

import (
	"encoding/json"
	"errors"
	"golanger.com/log"
	"golanger.com/net/http/client"
	. "haiyi/model"
	"haiyi/util"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Processor struct {
	*Base
}

func NewProcessor(b *Base) *Processor {
	return &Processor{
		b,
	}
}

func (p *Processor) GetLastLivenews() {
	if data, err := p.getLivenews(""); err != nil {
		log.Error("<GetLastLivenews> err: ", err)
	} else {
		for _, v := range data {
			arr := strings.Split(v, "#")
			m := ModelLivenews{}
			if arr[0] == "0" {
				obj, err := p.transferNews(arr, v)
				if err != nil {
					log.Error("<GetLastLivenews> err: ", err)
					continue
				}
				m = obj
			} else if arr[0] == "1" {
				obj, err := p.transferData(arr, v)
				if err != nil {
					log.Error("<GetLastLivenews> err: ", err)
					continue
				}
				m = obj
			}
			if msg, err := AddLivenews(m, p.ColLivenews.C(), true, p.Config.SaveImgPath); err != nil {
				log.Error("<GetLastLivenews> err: ", err)
			} else if msg == "exist" {
				log.Debug("<GetLastLivenews> livenews exist")
				return
			} else {
				log.Debug("<GetLastLivenews> add a livenews")
			}
		}
	}
}

func (p *Processor) AutoGetLivenews() {
	if count, err := GetLivenewsCount(p.ColLivenews.C()); err != nil {
		log.Error("<AutoGetLivenews> err: ", err)
		return
	} else {
		log.Debug("<AutoGetLivenews> count: ", count)
	} /*else if count > 20000 {
		return
	}*/
	maxId, err := GetOldestLivenewsOriginalId(p.ColLivenews.C())
	if err != nil {
		log.Error("<AutoGetLivenews> err: ", err)
		return
	}
	log.Debug("<AutoGetLivenews> maxId: ", maxId)
	if data, err := p.getLivenews(maxId); err != nil {
		log.Error("<AutoGetLivenews> err: ", err)
	} else {
		for _, v := range data {
			arr := strings.Split(v, "#")
			m := ModelLivenews{}
			if arr[0] == "0" {
				obj, err := p.transferNews(arr, v)
				if err != nil {
					log.Error("<AutoGetLivenews> err: ", err)
					continue
				}
				m = obj
			} else if arr[0] == "1" {
				obj, err := p.transferData(arr, v)
				if err != nil {
					log.Error("<AutoGetLivenews> err: ", err)
					continue
				}
				m = obj
			}
			if msg, err := AddLivenews(m, p.ColLivenews.C(), true, p.Config.SaveImgPath); err != nil {
				log.Error("<AutoGetLivenews> err: ", err)
			} else if msg == "exist" {
				log.Debug("<AutoGetLivenews> livenews exist")
				continue
			} else {
				log.Debug("<AutoGetLivenews> add a livenews")
			}
		}
	}

}

func (p *Processor) getLivenews(maxId string) (data []string, err error) {
	url := "http://m.jin10.com/flash?maxId=" + maxId
	if resp, err1 := client.NewClient().Get(url); err1 != nil {
		err = err1
		return
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			err = errors.New("error status: " + resp.Status)
			log.Error("<getLivenews> status code:", resp.StatusCode)
			return
		}
		jsonData, err2 := ioutil.ReadAll(resp.Body)
		if err2 != nil {
			err = err2
			return
		}
		err = json.Unmarshal(jsonData, &data)
	}
	return
}

func (p *Processor) transferNews(arr []string, originalContent string) (obj ModelLivenews, err error) {
	if len(arr) < 12 {
		err = errors.New("format error.")
		return
	}
	obj.Type = 0
	importanceStr := strings.TrimSpace(arr[1])
	importanceInt, _ := strconv.Atoi(importanceStr)
	obj.Importance = importanceInt
	obj.PublishTime = strings.TrimSpace(arr[2])
	obj.Content = util.ClearHtmlTags(strings.TrimSpace(arr[3]))
	obj.Img = strings.TrimSpace(arr[6])
	obj.OriginalId = strings.TrimSpace(arr[11])
	obj.OriginalContent = originalContent
	return
}

func (p *Processor) transferData(arr []string, originalContent string) (obj ModelLivenews, err error) {
	if len(arr) < 14 {
		err = errors.New("format error.")
		return
	}
	starStr := strings.TrimSpace(arr[6])
	star, _ := strconv.Atoi(starStr)
	effectStr := strings.TrimSpace(arr[7])
	effect := 0
	switch effectStr {
	case "利多":
		effect = 1
	case "利空":
		effect = 2
	default:
		break
	}
	obj.Type = 1
	obj.Time = strings.TrimSpace(arr[1])
	obj.Content = strings.TrimSpace(arr[2])
	obj.Prefix = strings.TrimSpace(arr[3])
	obj.Predicted = strings.TrimSpace(arr[4])
	obj.Actual = strings.TrimSpace(arr[5])
	obj.Star = star
	obj.Effect = effect
	obj.PublishTime = strings.TrimSpace(arr[8])
	obj.Country = strings.TrimSpace(arr[9])
	obj.OriginalId = strings.TrimSpace(arr[12])
	obj.OriginalContent = originalContent
	return
}

func (p *Processor) AutoGetCalendar() {
	y, m, err := GetLastSavedCalendar(p.ColSavedCalendar.C())
	if y <= 2010 {
		return
	}
	if err != nil {
		log.Error("<AutoGetCalendar> err: ", err)
		return
	}
	yStart, mStart := util.GetPreMonth(y, m)
	start := util.FormatDate(yStart, mStart, 1)
	end := util.FormatDate(y, m, 1)
	if err := p.GetCalendar(start, end); err != nil {
		log.Error("<AutoGetCalendar> err: ", err)
		return
	}
	if err := AddSavedCalendar(yStart, mStart, p.ColSavedCalendar.C()); err != nil {
		log.Error("<AutoGetCalendar> err: ", err)
	}
}

func (p *Processor) UpdateNewCalendar() {
	t := time.Now()
	yStart := int(t.Year())
	mStart := int(t.Month())
	day := int(t.Day())
	yEnd, mEnd := util.GetNextMonth(yStart, mStart)
	yEnd, mEnd = util.GetNextMonth(yEnd, mEnd)
	start := util.FormatDate(yStart, mStart, 1)
	end := util.FormatDate(yEnd, mEnd, day)
	if err := p.GetCalendar(start, end); err != nil {
		log.Error("<AutoGetCalendar> err: ", err)
	}
}

func (p *Processor) GetCalendar(start, end string) error {
	url := "https://api-markets.wallstreetcn.com//v1/calendar.json?start=" + start + "&end=" + end
	if resp, err := client.NewClient().Get(url); err != nil {
		log.Error("<GetCalendar> err: ", err)
		return err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Error("<GetCalendar> status code:", resp.StatusCode, " error:", err)
			return err
		}
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Error("<GetCalendar> err:", err)
			return err
		} else {
			calendarData := CalendarData{}
			if err := json.Unmarshal(body, &calendarData); err != nil {
				log.Error("<GetCalendar> err: ", err)
				return err
			}
			for _, calendar := range calendarData.Results {
				if err := SaveCalendar(calendar, p.ColCalendar.C()); err != nil {
					log.Error("<GetCalendar> err: ", err)
					return err
				}
			}
		}
	}
	return nil
}
