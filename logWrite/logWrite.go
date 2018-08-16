package logWrite

import (
	"fmt"
	"os"
	"io"
	"time"
)

var LastLog string
func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * 判断文件夹是否存在  存在返回 true 不存在返回false
 */
func checkDirIsExist(dirname string) bool {
	var exist = true
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err2 := os.MkdirAll(dirname,0777)
		if err2 != nil{
			exist = false
		}else {
			exist = true
		}
	}
	return exist
}
/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		_, err2 := os.OpenFile(filename, os.O_CREATE, 0666)
		if err2 != nil{
			exist = false
		}else {
			exist = true
		}
	}
	return exist
}
func WriteLog(writeString string){
	if LastLog == writeString{
		return
	}
	strTime := time.Now().Format("20060102_15")
	strDay := strTime[:8]
	strHour := strTime[len(strTime) - 2 :]

	strDir := fmt.Sprintf("./log/%s/", strDay)
	if !checkDirIsExist(strDir) {
		return
	}

	var filename =fmt.Sprintf("%s/%s_log.log", strDir, strHour)
	var f *os.File
	var err1 error
	
	if checkFileIsExist(filename) { //如果文件存在
		f, err1 = os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
		defer f.Close()
	} else {
		fmt.Println("文件不存在")
		return 
	}
	check(err1)
	writeString += "\n"
	_, err1 = io.WriteString(f, writeString) //写入文件(字符串)
	check(err1)
	LastLog = writeString
}