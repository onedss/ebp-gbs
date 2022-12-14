package app

import (
	"fmt"
	"github.com/onedss/ebp-gbs/models"
	"github.com/onedss/ebp-gbs/mylog"
	"github.com/onedss/ebp-gbs/routers"
	"github.com/onedss/ebp-gbs/service"
	"github.com/onedss/ebp-gbs/utils"
	"log"
)

type application struct {
	servers []OneServer
}

func (p *application) Start(s service.Service) (err error) {
	log.Println("********** START **********")
	for _, server := range p.servers {
		port := server.GetPort()
		if utils.IsPortInUse(port) {
			err = fmt.Errorf("TCP port[%d] In Use", port)
			return
		}
	}
	err = models.Init()
	if err != nil {
		return
	}
	err = routers.Init()
	if err != nil {
		return
	}
	for _, server := range p.servers {
		err := server.Start()
		if err != nil {
			return err
		}
	}
	go func() {
		for range routers.API.RestartChan {
			log.Println("********** STOP **********")
			for _, server := range p.servers {
				server.Stop()
			}
			utils.ReloadConf()
			log.Println("********** START **********")
			for _, server := range p.servers {
				err := server.Start()
				if err != nil {
					return
				}
			}
		}
	}()
	return nil
}

func (p *application) Stop(s service.Service) (err error) {
	defer log.Println("********** STOP **********")
	defer mylog.CloseLogWriter()
	for _, server := range p.servers {
		server.Stop()
	}
	models.Close()
	return
}

func (p *application) AddServer(server OneServer) {
	p.servers = append(p.servers, server)
}
