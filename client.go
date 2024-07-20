package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Cli struct {
	*http.Client
	Header http.Header
}

type Rsp struct {
	Code string `json:"code"`
}

func (c *Cli) Tickets(captcha string) (err error) {
	u := "/tickets"
	p := fmt.Sprintf(`{"type":"MFA","password":"%s","captcha":"%s"}`, password, captcha)
	d := url.Values{}
	d.Set("username", username)
	d.Set("password", p)
	resp, err := c.PostForm(u, d)
	if err != nil {
		return fmt.Errorf("请求tickets接口失败:%s", err)
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("创建tickets返回码异常:%d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取tickets接口返回报文失败:%s", err)
	}
	log.Println(string(b))
	return nil
}

func (c *Cli) GetToken(captcha string) (err error) {
	u := "/getToken"
	d := make(map[string]interface{})
	d["userName"] = username
	d["password"] = password
	d["captcha"] = captcha
	d["type"] = "MFA"
	d_json, _ := json.Marshal(d)
	resp, err := c.Post(u, "application/json", bytes.NewReader(d_json))
	if err != nil {
		return fmt.Errorf("请求getToken接口失败:%s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("调用getToken接口返回码异常:%d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取getToken接口返回报文失败:%s", err)
	}
	log.Println(string(b))
	var b2Json B2Json
	json.Unmarshal(b, &b2Json)
	if b2Json.Status != "ok" {
		return fmt.Errorf("解析getToken接口返回报文获取token失败:%s", b2Json)
	}
	c.Header = resp.Header
	c.Header.Set("Content-Type", "application/json")
	c.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	c.Header.Set("Accesstoken", b2Json.AccessToken)
	log.Println(c.Header.Get("Accesstoken"))
	return nil
}

func (c *Cli) GetUserInfo() (err error) {
	u := "/getUserInfo"
	req, _ := http.NewRequest("GET", u, bytes.NewReader(nil))
	req.Header = c.Header
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("请求getUserInfo接口失败:%s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("调用getUserInfo接口返回码异常:%d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取getUserInfo接口返回报文失败:%s", err)
	}
	log.Println(string(b))
	c.Header.Set("Accesstoken", resp.Header.Get("Accesstoken"))
	return nil
}

func (c *Cli) AddDutySign(days int, sign_type string) (err error) {
	u := "/addDutySign"
	d := make(map[string]interface{})
	d["dateId"] = days
	d["memberId"] = member_id
	d["moduleId"] = module_id
	d["signOutTime"] = time.Now().Format("2006-01-02 15:04:05")
	d["signType"] = sign_type
	d["userMail"] = username
	d_json, _ := json.Marshal(d)
	log.Println(string(d_json))
	req, _ := http.NewRequest("POST", u, bytes.NewReader(d_json))
	req.Header = c.Header
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("请求addDutySign接口失败:%s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("调用addDutySign接口返回码异常:%d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取addDutySign接口返回报文失败:%s", err)
	}
	log.Println(string(b))
	var rsp Rsp
	json.Unmarshal(b, &rsp)
	if rsp.Code != "0000" {
		return fmt.Errorf("addDutySign接口报文体返回码异常:%s", rsp.Code)
	}
	return nil
}

func (c *Cli) InsertSignCheckLog(sign_type string) (err error) {
	u := "/insertSignCheckLog"
	d := make(map[string]interface{})
	d["moduleId"] = module_id
	d["userMail"] = username
	d["signType"] = sign_type
	d["questionOneStatus"] = 1
	d["questionTwoStatus"] = 0
	d["questionThreeStatus"] = 0
	d["questionOneLog"] = ""
	d["questionTwoLog"] = ""
	d["questionThreeLog"] = ""
	d["remark"] = ""
	if sign_type == "out" {
		d["remark"] = "系统运行稳定"
	}
	d_json, _ := json.Marshal(d)
	req, _ := http.NewRequest("POST", u, bytes.NewReader(d_json))
	req.Header = c.Header
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("请求insertSignCheckLog接口失败:%s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("调用insertSignCheckLog接口返回码异常:%d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取insertSignCheckLog接口返回报文失败:%s", err)
	}
	log.Println(string(b))
	var rsp Rsp
	json.Unmarshal(b, &rsp)
	if rsp.Code != "0000" {
		return fmt.Errorf("insertSignCheckLog接口报文体返回码异常:%s", rsp.Code)
	}
	return nil
}
