package main

import (
	"net/http"
	"bytes"
	"fmt"
	"strings"
	"io/ioutil" //io 工具包
	//"time"
	"./logWrite"
)

var IP string


var accountstate = "user/accountstate"
var sign = "admin/sign"
var rawtransaction = "user/rawtransaction"
var strgetTransactionReceipt ="user/getTransactionReceipt"


type Wallet struct{
	Address string
	Balance string
	Nonce string
	Method string
	strResult string
	strRaw string
	hash string
	txhash string
	status string
}

func NewWallet(strAddr string) *Wallet {
	var newWallet = new(Wallet)
	newWallet.Address = strAddr
	return newWallet
}

func getPostUrl(cmd string) string{
	return fmt.Sprintf("http://%s:8685/v1/%s", IP, cmd)
}

//根据post返回的数据，提取出有价值的字段
func (w *Wallet)Json2Data(){
	strResult := w.strResult
	if strings.Contains(strResult, "error"){
		return
	}
	strResult = strResult[11 : len(strResult) - 2]

	strRange := strings.Split(strResult, ",")
	for _, strSub := range strRange{
		strNew := strings.Split(strSub,":")

		str0 := GetString(strNew[0])

		str1 := GetString(strNew[1])
		switch str0 {
		case "balance":
			w.Balance = str1
		case "nonce":
			w.Nonce = str1
		case "data":
			w.strRaw = str1
		case "txhash":
			w.txhash = str1
		}
	}
}

//获取并返回post数据
func getPostData(url string, post string) string{
	strWrite := fmt.Sprintf("\ncurl -i -H Accept:application/json -X POST %s -d '%s'", url, post)
	logWrite.WriteLog(strWrite)

	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	PostResult, _ := ioutil.ReadAll(resp.Body)
	strPostResult := string(PostResult)

	isPrint := true
	isPrint = isPrint && (!strings.Contains(strPostResult, "status"))
	isPrint = isPrint && (!strings.Contains(strPostResult, "data"))
	isPrint = isPrint && (!strings.Contains(strPostResult, "error"))
	if isPrint{
		println(strPostResult)
	}
	logWrite.WriteLog(strPostResult)
	return strPostResult
}

//查看钱包信息
func (this *Wallet) getaccountstate(){
	print("查询余额, ")
	postUrl := getPostUrl(accountstate)
	postJson := fmt.Sprintf("{\"address\":\"%s\"}", this.Address)
	println(postUrl, postJson)
	this.strResult = getPostData(postUrl, postJson)
	this.Json2Data()
}

func GetString(strSouce string) string {
	if strings.Contains(strSouce, "\""){
		return strSouce[1 : len(strSouce)-1]
	}
	return strSouce
}

func main() {
	//user1:n1XkoVVjswb5Gek3rRufqjKNpwrDdsnQ7Hq
	//user2:n1FF1nz6tarkDVwWQkMnnwFPuPKUaQTdptE
	//user3:n1UM7z6MqnGyKEPvUpwrfxZpM1eB7UpzmLJ
	//user4:n1UnCsJZjQiKyQiPBr7qG27exqCLuWUf1d7
	//newh1:n1bH1wF6wVcAk9szAzYga1ZqbSqjduiZoYb
	//newh2:n1Pn53VRExyR19hogeLqLJJb12ZfxsmT4gP
	//Pih1:n1aXC9QBsvBhUgt4WkEBeLxsbJeJV8DUWDr
	//Pih2:n1QyRg1DGBp6kgV6qZL9M7V7f1YEAiUGcMD

	AOCIPs := [2]string {"192.168.1.65", "192.168.1.66"}
	PiIPs := [2]string  {"192.168.1.37", "192.168.1.42"}

	for _, ip := range AOCIPs{
		IP = ip
		newh1 := NewWallet("n1bH1wF6wVcAk9szAzYga1ZqbSqjduiZoYb")
		newh2 := NewWallet("n1Pn53VRExyR19hogeLqLJJb12ZfxsmT4gP")
		newh1.getaccountstate()
		newh2.getaccountstate()
		println()
	}

	for _, ip := range PiIPs{
		IP = ip
		newh1 := NewWallet("n1aXC9QBsvBhUgt4WkEBeLxsbJeJV8DUWDr")
		newh2 := NewWallet("n1QyRg1DGBp6kgV6qZL9M7V7f1YEAiUGcMD")
		newh1.getaccountstate()
		newh2.getaccountstate()
		println()
	}
	//2006-01-02 15:04:05
	//for{
	//	strP := (time.Now().Format("20060102_15_04"))
	//	println(strP)
	//
	//	logWrite.WriteLog(strP)
	//	logWrite.LastLog = strP
	//	time.Sleep(time.Second * 5)
	//}
	//
	println("End...")
}