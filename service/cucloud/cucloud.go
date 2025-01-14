// Copyright 2025 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"

	"golang.org/x/exp/maps"
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

	signVal, err := c.generateRequestSign(reqHeader, reqBody)
	if err != nil {
		return err
	}
	reqHeader["sign"] = signVal
	reqHeader["Content-Type"] = "application/json"
	reqHeader["Account-Id"] = c.AccountId
	reqHeader["User-Id"] = c.AccountId
	reqHeader["Region-Code"] = c.CloudRegionCode

	bodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
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

func (c *CuCloud) generateRequestSign(header map[string]string, body map[string]string) (string, error) {
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
			return "", err
		}
		signRawString += k + "=" + string(vJson) + "&"
	}

	signRawString = signRawString[:len(signRawString)-1]
	mac.Write([]byte(signRawString))
	return hex.EncodeToString(mac.Sum(nil)), nil
}
