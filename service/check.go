package service

import (
	"fmt"
	"speedUpTx/config"
	"time"
)

func Start() {
	ticker := time.NewTicker(time.Minute * config.CheckPendingTime)
	for range ticker.C {
		for i := range config.SpeedUps {
			for _, net := range config.SpeedUps[i].Networks {
				fmt.Println(net.GasPriceUpper)
			}
			time.Sleep(time.Second * 5)
		}
	}
}
