package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetChaptcha() string {

	req, _ := http.NewRequest("GET", "http://localhost:8021/u/v1/base/captcha", nil)
	rsp, _ := http.DefaultClient.Do(req)
	if rsp.StatusCode != http.StatusOK {
		fmt.Println("[GetChaptcha] 获取验证码失败")
		return ""
	}
	body, _ := io.ReadAll(rsp.Body)
	rspJson := map[string]interface{}{}
	json.Unmarshal(body, &rspJson)
	return rspJson["answer"].(string)
}

func main() {
	// var k string
	// fmt.Println("输入y或Y发送手机短信")
	// fmt.Scanf("%s", &k)
	// if k == "y" || k == "Y" {
	// 	fmt.Println(1)
	// } else {
	// 	fmt.Println(k)
	// }
	// var sendSmsInfo struct {
	// 	Mobile string
	// 	Type   uint
	// }
	// sendSmsInfo.Mobile = "123"
	// sendSmsInfo.Type = uint(1)
	// fmt.Println(sendSmsInfo)
	// testMap := make(map[int]int)
	// testMap[2] = 2
	// testMap[1] = 1
	// if _, ok := testMap[1]; ok {
	// 	fmt.Println(ok)
	// }
	str := "TpROeKd5MzilW4u9$63b287466ca6b3469c0ec5b53a25c2a48f9d539e48171282af399a576f286646"
	fmt.Println(len(str))
}
