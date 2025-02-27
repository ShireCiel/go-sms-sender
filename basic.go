// Copyright 2021 The casbin Authors. All Rights Reserved.
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

package go_sms_sender

import "fmt"

const (
	Aliyun       = "Aliyun SMS"
	TencentCloud = "Tencent Cloud SMS"
	VolcEngine   = "Volc Engine SMS"
	Huyi         = "Huyi SMS"
	SubMail      = "SubMail SMS"
)

type SmsClient interface {
	SendMessage(param map[string]string, targetPhoneNumber ...string) error
}

func NewSmsClient(provider string, accessId string, accessKey string, sign string, template string, other ...string) (SmsClient, error) {
	switch provider {
	case Aliyun:
		return GetAliyunClient(accessId, accessKey, sign, template)
	case TencentCloud:
		return GetTencentClient(accessId, accessKey, sign, template, other)
	case VolcEngine:
		return GetVolcClient(accessId, accessKey, sign, template, other)
	case Huyi:
		return GetHuyiClient(accessId, accessKey, template)
	case SubMail:
		return GetSubMailClient(accessId, accessKey, template)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
