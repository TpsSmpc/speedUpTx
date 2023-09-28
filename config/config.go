package config

import (
	"crypto/ecdsa"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"speedUpTx/tools"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

var (
	seeds            [4]uint64
	password         string
	SpeedUps         []SpeedUp
	CheckPendingTime time.Duration // minute
)

//Load load config file
func Load(pwd string, _seeds [4]uint64) {
	password, seeds = pwd, _seeds
	file, err := os.Open("config/config.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	type Config struct {
		SpeedUps         []SpeedUp     `json:"speed_up"`
		CheckPendingTime time.Duration `json:"check_pending_time"`
	}
	all := &Config{}
	if err = json.NewDecoder(file).Decode(all); err != nil {
		panic(err)
	}

	SpeedUps = all.SpeedUps
	CheckPendingTime = all.CheckPendingTime

	for i, t := range SpeedUps {
		if len(t.KeyStore) > 0 {
			var keyjson []byte
			keyjson, err := ioutil.ReadFile(t.KeyStore)
			if err != nil {
				log.Panicf("Read keystore file fail. %s : %v\n", t.KeyStore, err)
			}
			SpeedUps[i].KeyStore = string(keyjson)
		}
	}

	checkKeyStore()
}

func checkKeyStore() {
	pwd := tools.GetDecryptString(password, seeds)
	for i, t := range SpeedUps {
		if len(t.KeyStore) == 0 {
			continue
		}
		keyjson := []byte(t.KeyStore)
		keyWrapper, err := keystore.DecryptKey(keyjson, string(pwd))
		if err != nil {
			panic(t.KeyStore + " keystore error : " + err.Error())
		}
		SpeedUps[i].Address = keyWrapper.Address
	}
}

type SpeedUp struct {
	Address  common.Address
	KeyStore string    `json:"keystore"`
	Networks []NetWork `json:"networks"`
}

type NetWork struct {
	Rpc           string `json:"rpc"`
	GasPriceUpper int64  `json:"gas_price_upper"`
}

func GetPrivateKey(i int) (common.Address, *ecdsa.PrivateKey, error) {
	pwd := tools.GetDecryptString(password, seeds)
	keyWrapper, err := keystore.DecryptKey([]byte(SpeedUps[i].KeyStore), string(pwd))
	if err != nil {
		return common.Address{}, nil, err
	}
	return keyWrapper.Address, keyWrapper.PrivateKey, nil
}
