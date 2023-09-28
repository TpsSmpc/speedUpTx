package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"speedUpTx/tools"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	SpeedUp()
}

func readRand() (string, [4]uint64) {
	data, err := os.ReadFile("rand.data")
	if err != nil {
		log.Panicf("read rand.data error. %v", err)
	}
	if err := os.WriteFile("rand.data", []byte("start the process successful! You are very great. Best to every one."), 0666); err != nil {
		log.Panicf("write rand.data error. %v", err)
	}
	os.Remove("rand.data")

	//generate seeds
	var seeds [4]uint64
	seeds[0] = tools.GenerateRandomSeed()
	seeds[1] = tools.GenerateRandomSeed()
	seeds[2] = tools.GenerateRandomSeed()
	seeds[3] = tools.GenerateRandomSeed()

	pwd := tools.GetEncryptString(string(data), seeds)
	return pwd, seeds
}

func input() {
	var pwd string
	fmt.Printf("Enter keystore's password: ")
	fmt.Scanf("%s", &pwd)
	//pwd = "1234567"
	if err := os.WriteFile("rand.data", []byte(pwd), 0666); err != nil {
		log.Panicf("write rand.data error. %v", err)
	}
}

func createKeyStore(pwd, filename string) {
	ks := keystore.NewKeyStore("./keystores", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(pwd)
	if err != nil {
		log.Fatal(err)
	}
	jsonData, err := ks.Export(account, pwd, pwd)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, jsonData, 0666); err != nil {
		log.Fatal(err)
	}
}

func SpeedUp() {
	var pwd string
	fmt.Printf("Input the keystore password(must ./smpc_k): ")
	fmt.Scanf("%s", &pwd)
	var keyjson []byte
	keyjson, _ = ioutil.ReadFile("smpc_k")
	keyWrapper, _ := keystore.DecryptKey(keyjson, pwd)

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/3f8b4373a4a943bf8b9c635fba90ee78")
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var txHash common.Hash
	addr := keyWrapper.Address.Hex()
	if addr == "0x380dF538Ab2587B11466d07ca5c671d33497d5Ca" {
		txHash = common.HexToHash("0x0e0238d7d30c5427338144344b0508c734aedea27f18840873d4e7935dabb63a")
	} else if addr == "0x5e80cf0C104D2D4f685A15deb65A319e95dd80dD" {
		txHash = common.HexToHash("0x0ae1a00045ac142a01da8af1ab050fac014222451eeaac7a636a7a2b29c38cb1")
	} else if addr == "0x3Fdd4B2d69848F74E44765e6AD423198bdBD94fa" {
		txHash = common.HexToHash("0xfa9c4703952d5cf0c5f6314ce91b753b99f8cd06a616f02eb567f07f95fc9073")
	} else if addr == "0x9dcb974Cf7522F91F2Add8303e7BCB2221063c48" {
		txHash = common.HexToHash("0x9f049eff81351ef4757b15a6859e6c0a4849996bac3a1a3edff7a7aa3e0949d7")
	} else if addr == "0xeBbe638eF6dF4A3837435bB44527f8D9BA9CF981" {
		txHash = common.HexToHash("0xfebf2ae44d2c86f64316762f67bb4d944eef25371c0279ce63f139e516c45cbd")
	}
	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("Get SuggestGasPrice error. %v", err)
	}
	gasPrice.Mul(gasPrice, big.NewInt(100+50))
	gasPrice.Div(gasPrice, big.NewInt(100))

	tx = types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), gasPrice, tx.Data())

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), keyWrapper.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(signedTx.Hash().Hex(), isPending)
}
