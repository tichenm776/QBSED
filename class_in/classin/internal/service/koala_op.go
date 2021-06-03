package service

import (
	"github.com/alecthomas/log4go"
	"time"
	"zhiyuan/koalamate_statistics_server/configs"

	"zhiyuan/device_manager/server_koala/koala_middleware/interface/koala"
)

func(s *Service) Koala_Action()(bool){
	configs.Init("./conf.yaml")
	koala_ip := configs.Gconf.KoalaHost
	koala_port := configs.Gconf.KoalaPort
	koala.G_koala.Init(koala_ip,koala_port)
	flag,err := koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	count := 0
	log4go.Info("koala login")
	if err!=nil || flag != true{
		time.Sleep(60*time.Second)
		count +=1
		if count >10{
			return false
		}
		return s.Koala_Action()
	}
	go func() {
		for{
			time.Sleep(30*time.Minute)
			koala_ip := configs.Gconf.KoalaHost
			koala_port := configs.Gconf.KoalaPort
			koala.G_koala.Init(koala_ip,koala_port)
			flag,err := koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
			if err!=nil || flag != true{
				log4go.Error("koala login keep live err",err)
			}
			//s.Koala_Get_subjects()
		}
	}()
	//s.Koala_Get_subjects()
	return flag
}
