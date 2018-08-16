package main

import (
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"time"
	"strings"
	"strconv"
)

//	const IP string = "localhost"		//本机运行
//	const IP string = "192.168.1.66"	// Aoch1
//	const IP string = "192.168.1.65"	// Aoch2
	var  IP string = "192.168.1.37" // Pih1
	//const IP string = "192.168.1.42"	// Pih2

var accountstate = "user/accountstate"
var sign = "admin/sign"
var rawtransaction = "user/rawtransaction"
var strgetTransactionReceipt = "user/getTransactionReceipt"
var unlock = "admin/account/unlock"
var unlockTx = "admin/transaction"

type Wallet struct {
	Address   string
	Balance   string
	Nonce     string
	Method    string
	strResult string // what Result?
	strRaw    string
	hash      string
	txhash    string
	status    string
	result    string // unlock result(true or flase)
}

func NewWallet(strAddr string) *Wallet {
	var newWallet = new(Wallet)
	newWallet.Address = strAddr
	return newWallet
}

func getPostUrl(cmd string) string {
	return fmt.Sprintf("http://%s:8685/v1/%s", IP, cmd)
}

//根据post返回的数据，提取出有价值的字段
func (w *Wallet) Json2Data() {
	strResult := w.strResult
	if strings.Contains(strResult, "error") {
		return
	}
	strResult = strResult[11 : len(strResult)-2]

	strRange := strings.Split(strResult, ",")
	for _, strSub := range strRange {
		strNew := strings.Split(strSub, ":")

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
		case "result":
			w.result = str1
		}
	}
}

//获取并返回post数据
func getPostData(url string, post string) string {
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
	if isPrint {
		println(strPostResult)
	} else if strings.Contains(strPostResult, "error") {
		//		println()
	}
	return strPostResult
}

//查看钱包信息
func (this *Wallet) getaccountstate() {
	print("\n查询余额, ")
	postUrl := getPostUrl(accountstate)
	postJson := fmt.Sprintf("{\"address\":\"%s\"}", this.Address)
	println(postUrl, postJson)
	this.strResult = getPostData(postUrl, postJson)
	this.Json2Data()
}

func GetString(strSouce string) string {
	if strings.Contains(strSouce, "\"") {
		return strSouce[1 : len(strSouce)-1]
	}
	return strSouce
}

//签名发送交易，返回hash值
func (this *Wallet) TxSign(Wto *Wallet, value int64) string {
	println("\nSendTransactionWithSign, ")
	if value == 0 {
		value = 1
	}
	strValue := strconv.FormatInt(value, 10)
	postUrl := getPostUrl(sign)
	CurrentiNonce, _ := strconv.Atoi(this.Nonce)
	CurrentiNonce += 1
	CurrentNonce := strconv.Itoa(CurrentiNonce)

	strValue = "200000000000000000000000"
	from2to := fmt.Sprintf("\"from\":\"%s\", \"to\":\"%s\", \"value\":\"%s\"", this.Address, Wto.Address, strValue)
	gas := "\"gasPrice\":\"1000000\",\"gasLimit\":\"2000000\""
	transaction := fmt.Sprintf("{%s, \"nonce\":%s,%s}", from2to, CurrentNonce, gas)
	postJson := fmt.Sprintf("{\"transaction\":%s, \"passphrase\":\"passphrase\"}", transaction)
	println(postUrl, postJson)
	this.strResult = getPostData(postUrl, postJson)
	this.Json2Data()
	this.RawTransaction()
	return this.hash
}

//解锁账户，返回解锁结果
func Unlock(Wfrom *Wallet) string {
	println("\nUnlock.........")

	postUrl := getPostUrl(unlock)

	//curl -i -H 'Content-Type: application/json' -X POST http://localhost:8685/v1/admin/account/unlock -d
	// '{"address":"n1FF1nz6tarkDVwWQkMnnwFPuPKUaQTdptE","passphrase":"passphrase","duration":"300000000000"}'

	postJson := fmt.Sprintf("{\"address\":%s, \"passphrase\":\"passphrase\",\"duration\":\"300000000000\"}", Wfrom.Address)
	println(postUrl, postJson)
	Wfrom.strResult = getPostData(postUrl, postJson)
	Wfrom.Json2Data()
	Wfrom.RawTransaction()
	return Wfrom.hash
}

//解锁发送交易，返回hash值
func UnlockTx(Wfrom, Wto *Wallet, value int64) string {
	return "不能解锁！！！要在节点操作。。"
	strResult := Unlock(Wto)
	if strings.Contains(strResult, "false") {
		println("\nUnlock Error!!!, ")
		return ""
	}
	println("\nSendTransactionWithUnlock, ")
	if value == 0 {
		value = 1
	}
	strValue := strconv.FormatInt(value, 10)
	postUrl := getPostUrl(sign)
	CurrentiNonce, _ := strconv.Atoi(Wfrom.Nonce)
	CurrentiNonce += 1
	CurrentNonce := strconv.Itoa(CurrentiNonce)

	from2to := fmt.Sprintf("\"from\":\"%s\", \"to\":\"%s\", \"value\":\"%s\"", Wfrom.Address, Wto.Address, strValue)
	gas := "\"gasPrice\":\"1000000\",\"gasLimit\":\"2000000\""
	postJson := fmt.Sprintf("{%s, \"nonce\":%s,%s}", from2to, CurrentNonce, gas)

	println(postUrl, postJson)
	Wfrom.strResult = getPostData(postUrl, postJson)
	Wfrom.Json2Data()
	Wfrom.RawTransaction()
	return Wfrom.hash
}

//交易上链，返回上链后的 hash
func (w *Wallet) RawTransaction() string {
	//print("\n交易上链, ")
	print("\n交易收据, ")
	postUrl := getPostUrl(rawtransaction)
	postJson := fmt.Sprintf("{\"data\":\"%s\"}", w.strRaw)
	//println(postUrl, postJson)
	w.strResult = getPostData(postUrl, postJson)
	w.Json2Data() //get unlock result(true or false)
	return w.result
}

//根据hash值，返回交易收据 状态值
func getTransactionReceipt(hash string) string {
	//print("\n交易收据, ")
	if hash == "" {
		println("hash is null！")
		return ""
	}
	strStatus := "0"
	postUrl := getPostUrl(strgetTransactionReceipt)
	postJson := fmt.Sprintf("{\"hash\":\"%s\"}", hash)
	//println(postUrl, postJson)

	strResult := getPostData(postUrl, postJson)
	if strings.Contains(strResult, "error") {
		return strStatus
	}

	strResult = strResult[11 : len(strResult)-2]

	strRange := strings.Split(strResult, ",")
	for _, strSub := range strRange {
		strNew := strings.Split(strSub, ":")
		str0 := GetString(strNew[0])
		str1 := GetString(strNew[1])
		if str0 == "status" {
			strStatus = str1
			//println(str0, str1)
			return strStatus
		}
	}
	return strStatus
}

func main() {
	PCuser1 := "n1XkoVVjswb5Gek3rRufqjKNpwrDdsnQ7Hq"
	PCuser2 := "n1FF1nz6tarkDVwWQkMnnwFPuPKUaQTdptE"
	//user3 := "n1UM7z6MqnGyKEPvUpwrfxZpM1eB7UpzmLJ"
	//user4 := "n1UnCsJZjQiKyQiPBr7qG27exqCLuWUf1d7"
	//newh1:n1bH1wF6wVcAk9szAzYga1ZqbSqjduiZoYb
	//newh2:n1Pn53VRExyR19hogeLqLJJb12ZfxsmT4gP
	Pih1 := "n1d8F3u7FjPARwpbiQdsaExaFWdAfjjRDrg"
	Pih2 := "n1QVhi6j72fc2vtS7CX9tSxqQfHbdBsNnmp"


	user1 := NewWallet(PCuser1)
	user2 := NewWallet(PCuser2)
	newh1 := NewWallet("n1FF1nz6tarkDVwWQkMnnwFPuPKUaQTdptE")
	newh2 := NewWallet(Pih1)
	newh3 := NewWallet(Pih2)

	user1.getaccountstate()
	user2.getaccountstate()
	newh1.getaccountstate()
	newh2.getaccountstate()
	newh3.getaccountstate()
	//return
	//newh1.hash = UnlockTx(newh1, newh2, 1)

	newh1.TxSign(newh2, 1)
	//newh2.hash = TxSign(newh1, newh2, 1)
	iSecond := 0

	for {
		//if iSecond % 60 == 0{
		println(fmt.Sprintf("Current wait time : %3d s", iSecond))
		//}
		if newh1.txhash == "" {
			println("hash is null！")
			return
		}
		newh1.status = getTransactionReceipt(newh1.txhash)
		if newh1.status == "1" || newh1.status == "0" {
			println()
			break
		}
		iSecond += 1
		if iSecond > 60*6 {
			println()
			break
		}
		time.Sleep(1 * time.Second)
	}
	newh1.getaccountstate()
	newh2.getaccountstate()
	//}
	println("End...")
}
