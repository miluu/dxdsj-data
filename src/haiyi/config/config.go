package config

import (
	"time"
)

type Conf struct {
	MaxProcs               int           `json:"max_procs"`
	MgoDns                 string        `json:"mgo_dns"`
	MgoMode                string        `json:"mgo_mode"`
	MgoRefresh             bool          `json:"mgo_refresh"`
	LocateRelativeExecPath bool          `json:"locate_relative_exec_path"`
	LogDir                 string        `json:"log_dir"`
	LogFile                string        `json:"log_file"`
	LogDebugLevel          string        `json:"log_debug_level"`
	TimeInterval           time.Duration `json:"time_interval"`
	SaveImgPath            string        `json:"save_img_path"`
}
