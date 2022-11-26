package app

import (
	"github.com/common-nighthawk/go-figure"
	"github.com/onedss/ebp-gbs/client"
	"github.com/onedss/ebp-gbs/service"
	"github.com/onedss/ebp-gbs/utils"
	"log"
	"os"
)

func StartApp() {
	log.Println("ConfigFile -->", utils.ConfFile())
	sec := utils.Conf().Section("service")
	svcConfig := &service.Config{
		Name:        sec.Key("name").MustString("EbpGBS_Service"),
		DisplayName: sec.Key("display_name").MustString("EbpGBS_Service"),
		Description: sec.Key("description").MustString("EbpGBS_Service"),
	}

	httpPort := utils.Conf().Section("http").Key("port").MustInt(51180)
	oneHttpServer := client.NewOneHttpServer(httpPort)
	p := &application{}
	p.AddServer(oneHttpServer)

	var s, err = service.New(p, svcConfig)
	if err != nil {
		log.Println(err)
		utils.PauseExit()
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "install" || os.Args[1] == "stop" {
			figure.NewFigure("Ebp-GBS", "", false).Print()
		}
		log.Println(svcConfig.Name, os.Args[1], "...")
		if err = service.Control(s, os.Args[1]); err != nil {
			log.Println(err)
			utils.PauseExit()
		}
		log.Println(svcConfig.Name, os.Args[1], "ok")
		return
	}
	figure.NewFigure("Ebp-GBS", "", false).Print()
	if err = s.Run(); err != nil {
		log.Println(err)
		utils.PauseExit()
	}
}
