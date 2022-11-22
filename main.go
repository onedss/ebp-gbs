package main

import (
	"fmt"
	"github.com/onedss/ebp-gbs/buildtime"
	"github.com/onedss/ebp-gbs/mylog"
	"github.com/onedss/ebp-gbs/utils"
	"log"
)

var (
	gitCommitCode string
	buildDateTime string
)

func main() {
	log.SetPrefix("[Ebp-GBS] ")
	if utils.Debug {
		log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	} else {
		log.SetFlags(log.LstdFlags)
	}
	buildtime.BuildVersion = fmt.Sprintf("%s.%s", buildtime.BuildVersion, gitCommitCode)
	buildtime.BuildTimeStr = buildtime.BuildTime.Format(utils.DateTimeLayout)
	mylog.Info("BuildVersion:", buildtime.BuildVersion)
	mylog.Info("BuildTimeStr:", buildtime.BuildTimeStr)
}
