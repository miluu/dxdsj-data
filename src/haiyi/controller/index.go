package controller

import (
	"encoding/json"
	"golanger.com/log"
	"golanger.com/net/http/client"
	. "haiyi/model"
	"haiyi/util"
	"io/ioutil"
	"strconv"
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

func (p *Processor) GetLivenews() {
	lastSavedPage, err := GetLastSavedPage(p.ColSavedPage.C())
	if err != nil {
		log.Error("<GetLivenews> err: ", err)
		return
	}
	page := lastSavedPage + 1
	if page > 1000 {
		return
	}
	url := "http://api.wallstreetcn.com/v2/livenews?order=-created_at&limit=100&channelId=0&extractImg=1&extractText=1&page=" + strconv.Itoa(page)
	if resp, err := client.NewClient().Get(url); err != nil {
		log.Error("<GetLivenews> err: ", err)
		return
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Error("<GetLivenews> status code:", resp.StatusCode, " error:", err)
			return
		}
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Error("<GetLivenews> err:", err)
			return
		} else {
			originalNews := LivenewsOriginalData{}
			if err := json.Unmarshal(body, &originalNews); err != nil {
				log.Error("<GetLivenews> err: ", err)
				return
			}
			for _, livenews := range originalNews.Results {
				if _, err := AddLivenews(livenews, p.ColLivenews.C(), p.Config.SaveImgPath); err != nil {
					log.Error("<GetLivenews> err: ", err)
				} /* else if msg == "exist" {
					log.Debug("<GetLivenews> exist.")
				} else {
					log.Debug("<GetLivenews> saved.")
				}*/
			}
			AddSavedPage(page, p.ColSavedPage.C())
		}
	}
}

func (p *Processor) GetLastLivenews() {
	url := "http://api.wallstreetcn.com/v2/livenews?order=-created_at&limit=20&channelId=0&extractImg=1&extractText=1"
	if resp, err := client.NewClient().Get(url); err != nil {
		log.Error("<GetLastLivenews> err: ", err)
		return
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			log.Error("<LivenewsApi> status code:", resp.StatusCode, " error:", err)
			return
		}
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Error("<LivenewsApi> err:", err)
			return
		} else {
			originalNews := LivenewsOriginalData{}
			if err := json.Unmarshal(body, &originalNews); err != nil {
				log.Error("<LivenewsApi> err: ", err)
				return
			}
			for _, livenews := range originalNews.Results {
				if _, err := AddLivenews(livenews, p.ColLivenews.C(), p.Config.SaveImgPath); err != nil {
					log.Error("<LivenewsApi> err: ", err)
				} /* else if msg == "exist" {
					log.Debug("<LivenewsApi> exist.")
				} else {
					log.Debug("<LivenewsApi> saved.")
				}*/
			}
		}
	}
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
