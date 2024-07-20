package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Notify struct {
	JobName      string
	DingtalkFlag bool
}

func (a *Notify) Dingtalk(msg string) {
	reqStr := fmt.Sprintf(`	{
    	"version": "4",
    	"status": "firing",
    	"alerts": [
    	    {
    	        "status": "firing",
    	        "labels": {
    	            "告警名称": "%s"
    	        },
    	        "annotations": {
    	            "告警描述": "%s",
    	            "时间":"%s"
    	        }
    	    }
    	]
	}`, a.JobName, msg, time.Now().Format("2006-01-02 15:04:05"))
	log.Println(reqStr)
	_, err := http.Post("/dingtalk/webhook1/send", "application/json", bytes.NewBuffer([]byte(reqStr)))
	if err != nil {
		log.Println("请求告警接口失败:" + err.Error())
	}
}

func (a *Notify) Println(msg string) {
	log.Println(msg)
}

func (a *Notify) PrintlnWithDingtalk(msg string) {
	if a.DingtalkFlag {
		a.Dingtalk(msg)
	}
	log.Println(msg)
}
