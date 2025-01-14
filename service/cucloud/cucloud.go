package cucloud

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

type CuCloud struct {
	AccessKey       string
	SecretKey       string
	TopicName       string
	MessageTitle    string
	CloudRegionCode string
	AccountId       string
	NotifyType      string
	Client          http.Client
}

type CuCloudResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func New(accessKey, secretKey, topicName, messageTitle, cloudRegionCode, accountId, notifyType string) *CuCloud {
	return &CuCloud{
		accessKey,
		secretKey,
		topicName,
		messageTitle,
		cloudRegionCode,
		accountId,
		notifyType,
		http.Client{},
	}
}

func (c *CuCloud) Send(ctx context.Context, subject, content string) error {
	timeNow := time.Now().UnixMilli()

	reqHeader := map[string]string{}
	reqHeader["algorithm"] = "HmacSHA256"
	reqHeader["requestTime"] = strconv.FormatInt(timeNow, 10)
	reqHeader["accessKey"] = c.AccessKey

	reqBody := map[string]string{}

	reqBody["notifyType"] = c.NotifyType
	reqBody["messageTitle"] = url.QueryEscape(c.MessageTitle)
	reqBody["messageType"] = "text"
	reqBody["messageContent"] = url.QueryEscape(content)
	reqBody["topicName"] = url.QueryEscape(c.TopicName)
	reqBody["messageTag"] = ""
	reqBody["templateName"] = ""
	reqBody["cloudRegionCode"] = "cn-langfang-2"

	signVal := c.generateRequestSign(reqHeader, reqBody)
	reqHeader["sign"] = signVal
	reqHeader["Content-Type"] = "application/json"
	reqHeader["Account-Id"] = c.AccountId
	reqHeader["User-Id"] = c.AccountId
	reqHeader["Region-Code"] = c.CloudRegionCode

	bodyJson, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://gateway.cucloud.cn/smn/SMNService/api/message/notify", bytes.NewReader(bodyJson))
	if err != nil {
		return err
	}

	for k, v := range reqHeader {
		req.Header.Set(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respBody CuCloudResp

	err = json.Unmarshal(respBodyRaw, &respBody)
	if err != nil {
		return err
	}

	if respBody.Code != 200 {
		return fmt.Errorf(respBody.Message)
	}

	return nil
}

func (c *CuCloud) generateRequestSign(header map[string]string, body map[string]string) string {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	reqSignMap := make(map[string]string)
	maps.Copy(reqSignMap, header)
	maps.Copy(reqSignMap, body)
	signRawString := ""

	var keys []string
	for k, _ := range reqSignMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		vJson, err := json.Marshal(reqSignMap[k])
		if err != nil {
			panic(err)
		}
		signRawString += k + "=" + string(vJson) + "&"
	}

	signRawString = signRawString[:len(signRawString)-1]
	mac.Write([]byte(signRawString))
	return hex.EncodeToString(mac.Sum(nil))
}
