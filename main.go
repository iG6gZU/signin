package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"flag"
	"log"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"

	"github.com/robfig/cron"
)

var (
	// 账号
	username = ""
	// 用户id，从页面报文中获取
	member_id = 0
	// aes加密后的密码,从登录接口请求报文中获取
	password = ""
	// MFA
	secret = ""
	// AES加密key
	key = ""

	// 项目id，从签到记录接口返回报文中获取
	module_id = 0

	// 签到签退开关，用于调试，默认
	sign_flag = flag.Bool("sign_flag", true, "签到签退开关")

	// 定时任务
	cron_in  = flag.String("cron_in", "0 0 7 * * *", "签到定时")
	cron_out = flag.String("cron_out", "0 5 20 * * *", "签退定时")

	// 钉钉告警开关，用于调试
	dingtalk = flag.Bool("dingtalk", true, "告警开关")
)

type B2Json struct {
	Status      string `json:"status"`
	AccessToken string `json:"accessToken"`
}

func SignIn(signFlag bool, signType string) {
	log.Println("定时任务开始执行...")
	signTypeMap := map[string]string{"in": "签到", "out": "签退"}
	signTypeDesc := signTypeMap[signType]
	notify := Notify{JobName: "天眼" + signTypeDesc + "任务", DingtalkFlag: *dingtalk}

	n, _ := rand.Int(rand.Reader, big.NewInt(100))

	notify.Println("随机等待" + strconv.Itoa(int(n.Int64())) + "秒...")

	time.Sleep(time.Duration(n.Int64()) * time.Second)

	notify.Println(signTypeDesc + "任务开始执行...")

	code := strconv.Itoa(int(getCode(secret, 0)))
	notify.Println("MFA code:" + code)
	captcha := hex.EncodeToString(AesEncryptCBC([]byte(code), []byte(key)))
	init_day, _ := time.Parse("2006-01-02", "2020-07-04")
	today := time.Now()
	days := int(today.Sub(init_day).Hours() / 24)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	jar, _ := cookiejar.New(nil)

	c := Cli{
		Client: &http.Client{
			Jar:       jar,
			Transport: tr,
		},
	}

	err := c.Tickets(captcha)
	if err != nil {
		notify.PrintlnWithDingtalk("账号(" + username + ")" + signTypeDesc + "异常:" + err.Error())
		return
	}
	err = c.GetToken(captcha)
	if err != nil {
		notify.PrintlnWithDingtalk("账号(" + username + ")" + signTypeDesc + "异常:" + err.Error())
		return
	}
	err = c.GetUserInfo()
	if err != nil {
		notify.PrintlnWithDingtalk("账号(" + username + ")" + signTypeDesc + "异常:" + err.Error())
		return
	}
	if signFlag {
		err = c.AddDutySign(days, signType)
		if err != nil {
			notify.PrintlnWithDingtalk("账号(" + username + ")" + signTypeDesc + "异常:" + err.Error())
			return
		}
		err = c.InsertSignCheckLog(signType)
		if err != nil {
			notify.PrintlnWithDingtalk("账号(" + username + ")" + signTypeDesc + "异常:" + err.Error())
			return
		}
	}
	notify.PrintlnWithDingtalk("账号(" + username + ")" + signTypeDesc + "成功")
}

func main() {
	flag.Parse()
	c := cron.New()
	c.AddFunc(*cron_in, func() { SignIn(*sign_flag, "in") })
	log.Println("签到定时任务已添加:" + *cron_in)
	c.AddFunc(*cron_out, func() { SignIn(*sign_flag, "out") })
	log.Println("签退定时任务已添加:" + *cron_out)
	c.Start()
	select {}
}
