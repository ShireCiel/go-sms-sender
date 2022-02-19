package go_sms_sender

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type SubMailClient struct {
	appId    string
	appKey   string
	template string
}

const (
	API = "https://api.mysubmail.com/message/xsend"
)

func GetSubMailClient(accessId string, accessKey string, templateId string) (*SubMailClient, error) {
	if len(accessId) < 1 {
		return nil, fmt.Errorf("missing parameter: appId")
	}

	submailClient := &SubMailClient{
		appId:    accessId,
		appKey:   accessKey,
		template: templateId,
	}

	return submailClient, nil
}

func (c *SubMailClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	//所需参数
	vars := make(map[string]string)
	vars["code"] = "123456"
	postdata := make(map[string]string)
	postdata["appid"] = c.appId
	postdata["signature"] = ""
	postdata["project"] = c.template
	postdata["to"] = targetPhoneNumber[0]
	postdata["timestamp"] = ""
	postdata["sign_type"] = "md5"
	postdata["sign_version"] = "2"
	bs, _ := json.Marshal(vars)
	postdata["vars"] = string(bs)
	//获取服务器时间戳，该时间戳为 UNIX 时间戳，也可以自己生成
	q, _ := http.Get("https://api.mysubmail.com/service/timestamp")
	r, _ := ioutil.ReadAll(q.Body)
	m := make(map[string]float64)
	json.Unmarshal(r, &m)
	postdata["timestamp"] = strconv.FormatFloat(m["timestamp"], 'f', -1, 64)
	//签名加密
	sign := make(map[string]string)
	sign["appid"] = postdata["appid"]
	sign["to"] = postdata["to"]
	sign["project"] = postdata["project"]
	sign["timestamp"] = postdata["timestamp"]
	sign["sign_type"] = postdata["sign_type"]
	sign["sign_version"] = postdata["sign_version"]
	keys := make([]string, 0, 32)
	for key := range sign {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	sign_list := make([]string, 0, 32)
	for _, key := range keys {
		sign_list = append(sign_list, key+"="+sign[key])
	}
	sign_str := c.appId + c.appKey + strings.Join(sign_list, "&") + c.appId + c.appKey
	mymd5 := md5.New()
	io.WriteString(mymd5, sign_str)
	postdata["signature"] = hex.EncodeToString(mymd5.Sum(nil))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range postdata {
		_ = writer.WriteField(key, val)
	}
	contentType := writer.FormDataContentType()
	writer.Close()
	resp, _ := http.Post(API, contentType, body)
	result, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(result))
	return nil
}
