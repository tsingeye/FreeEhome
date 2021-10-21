package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func init() {
	rand.Seed(time.Now().Unix())
}

//海康GPS转换公式
var GpsCalcTransform = func(v int64) float64 {
	_du := int64(v / 3600 / 100)
	_fen := int64((v - _du*3600*100) / 100 / 60)
	_miao := int64((v - _du*3600*100 - _fen*60*100) / 100)
	_latlon := float64(_du) + float64(_fen)/60 + float64(_miao)/60/60

	return _latlon
}

//MD5加密，生成32位MD5字符串
func GetMD5String(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

//十六进制转十进制
func HexToBigInt(hex string) *big.Int {
	bigInt := new(big.Int)
	num, _ := bigInt.SetString(hex, 16)
	return num
}

//生成UUID
func GetUUID() string {
	return strings.ToUpper(uuid.Must(uuid.NewV4()).String())
}

func GetRangeNum(min, max int64) int64 {
	num := rand.Int63n(max-min) + min
	if num%2 == 0 {
		return num
	} else {
		return num - 1
	}
}

//将GB2312类型的XML数据转码成UTF-8类型
func GB2312ToUTF8(data []byte) (utf8Data []byte, err error) {
	utf8Data, err = simplifiedchinese.HZGB2312.NewDecoder().Bytes(data)
	if err != nil {
		return utf8Data, err
	}
	utf8Data = bytes.Replace(utf8Data, []byte(`encoding="GB2312"?`), []byte(`encoding="UTF-8"?`), 1)

	return utf8Data, nil
}

//将UTF-8类型的XML数据转码成GB2312类型
func UTF8ToGB2312(data []byte) (gb2312Data []byte, err error) {
	data = []byte(StringsJoin(`<?xml version="1.0" encoding="GB2312" ?>`, "\n", string(data), "\n"))
	gb2312Data, err = simplifiedchinese.HZGB2312.NewEncoder().Bytes(data)
	if err != nil {
		return gb2312Data, err
	}

	return gb2312Data, nil
}

//字符串拼接
func StringsJoin(str ...string) string {
	var b bytes.Buffer
	strLen := len(str)
	if strLen == 0 {
		return ""
	}
	for i := 0; i < strLen; i++ {
		b.WriteString(str[i])
	}

	return b.String()
}

//字符串逆序
func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

//返回当前工作目录的绝对路径
func GetAbsPath() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0])) //作为服务时使用
	//path, err := os.Getwd()
	if err != nil {
		log.Fatalln("GetAbsPath() error: ", err)
	}
	path = filepath.ToSlash(path)
	return StringsJoin(path, "/")
}

//根据文件路径创建对应的文件夹
func MkdirAllFile(path, appPath string) string {
	//filePath文件路径，fileName文件名
	filePath, fileName := filepath.Split(path)
	//日志文件的绝对路径
	filePath = StringsJoin(appPath, filePath)
	//panic日志文件的绝对路径+文件名
	path = StringsJoin(filePath, fileName)
	//根据文件路径创建对应的文件夹
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Fatalln("panicFile failed to MkdirAll(): ", err)
	}
	return path
}
