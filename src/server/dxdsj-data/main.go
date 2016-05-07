package main

import (
	"flag"
	"golanger.com/config"
	"golanger.com/db/mongo"
	"golanger.com/log"
	cfg "haiyi/config"
	. "haiyi/controller"
	. "haiyi/model"
	"os"
	"path"
	"runtime"
	"time"
)

var (
	confFile = flag.String("f", "./config/config.conf", "config file")
)

func main() {
	flag.Parse()
	var conf cfg.Conf
	config.Files(*confFile).Load(&conf)
	cdpath := ""
	if conf.LocateRelativeExecPath {
		cdpath = path.Dir(os.Args[0]) + "/"
	}

	log.SetLevel(conf.LogDebugLevel)
	if conf.LogDir != "" {
		if _, err := os.Stat(cdpath + conf.LogDir); err != nil {
			os.Mkdir(cdpath+conf.LogDir, 0700)
		}
	}

	if conf.LogFile != "" {
		logFi, err := os.OpenFile(cdpath+conf.LogDir+"/"+conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0700)
		if err != nil {
			log.Fatalln(err)
		}

		log.SetOutput(logFi)
	}

	runtime.GOMAXPROCS(conf.MaxProcs)

	//Init Begin
	mon := mongo.NewMongoPool("", conf.MgoDns, conf.MgoMode, conf.MgoRefresh, 0, 0, 100)
	monLivenews := mon.C(ColLivenews)
	monSavedPage := mon.C(ColSavedPage)
	monCalendar := mon.C(ColCalendar)
	monSavedCalendar := mon.C(ColSavedCalendar)
	defer mon.Close()
	base := &Base{
		Config:           conf,
		ColLivenews:      monLivenews,
		ColSavedPage:     monSavedPage,
		ColCalendar:      monCalendar,
		ColSavedCalendar: monSavedCalendar,
	}
	//Init End

	processor := NewProcessor(base)
	processor.GetLastLivenews()
	processor.GetLivenews()
	processor.AutoGetCalendar()
	processor.UpdateNewCalendar()

	timer := time.NewTicker(conf.TimeInterval * time.Second)
	for {
		select {
		case <-timer.C:
			processor.GetLastLivenews()
			processor.GetLivenews()
			processor.AutoGetCalendar()
			processor.UpdateNewCalendar()
		}
	}
}
