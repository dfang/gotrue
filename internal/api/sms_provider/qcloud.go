package sms_provider

import (
	"fmt"

	"github.com/supabase/auth/internal/conf"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

const (
	// defaultTextLocalApiBase    = "https://api.textlocal.in"
	defaultQcloudSmsApiBase = "https://sms.tencentcloudapi.com"
	// textLocalTemplateErrorCode = 80
)

type QcloudProvider struct {
	Config  *conf.QcloudSmsProviderConfiguration
	APIPath string
}

type QcloudError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type QcloudSmsResponse struct {
	// *tchttp.BaseResponse
	Response *sms.SendSmsResponseParams `json:"Response"`
	// Response struct {
	// SendStatusSet []SendStatus `json:"SendStatusSet"`
	// RequestID     string       `json:"RequestId"`
	// } `json:"Response"`
}

type SendStatus struct {
	SerialNo       string `json:"SerialNo"`
	PhoneNumber    string `json:"PhoneNumber"`
	Fee            int    `json:"Fee"`
	SessionContext string `json:"SessionContext"`
	Code           string `json:"Code"`
	Message        string `json:"Message"`
	IsoCode        string `json:"IsoCode"`
}

// Creates a SmsProvider with the qcloud sms Config
func NewQcloudSmsProvider(config conf.QcloudSmsProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	apiPath := defaultQcloudSmsApiBase + ""
	return &QcloudProvider{
		Config:  &config,
		APIPath: apiPath,
	}, nil
}

func (t *QcloudProvider) SendMessage(phone, message, channel, otp string) (string, error) {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, message, channel)
	default:
		return "", fmt.Errorf("channel type %q is not supported for qcloud SMS", channel)
	}
}

// Send an SMS containing the OTP with qcloud SMS's API
func (t *QcloudProvider) SendSms(phone string, message string, channel string) (string, error) {
	// https://console.cloud.tencent.com/api/explorer?Product=sms&Version=2021-01-11&Action=SendSms
	credential := common.NewCredential(
		t.Config.SecretID,
		t.Config.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.Debug = true
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, err := sms.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		panic(err)
	}

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := sms.NewSendSmsRequest()

	// request.PhoneNumberSet = common.StringPtrs([]string{"+8615618903080"})
	// request.SmsSdkAppId = common.StringPtr(t.Config.SmsSdkAppID)
	// request.SignName = common.StringPtr("搭个线")
	// request.TemplateId = common.StringPtr(t.Config.TemplateID)
	request.PhoneNumberSet = common.StringPtrs([]string{"15618903080"})
	request.SmsSdkAppId = common.StringPtr("1400874085")
	request.SignName = common.StringPtr("武汉搭个线网络科技")
	request.TemplateId = common.StringPtr("2017731")
	request.TemplateParamSet = common.StringPtrs([]string{"987654"})

	// 返回的resp是一个SendSmsResponse的实例，与请求对象对应
	response, err := client.SendSms(request)
	fmt.Printf("response\n")
	fmt.Printf("%+v\n", response)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		// return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
	// logger.DebugJSON("短信[腾讯云]", "请求内容", request)
	// logger.DebugJSON("短信[腾讯云]", "接口响应", response)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		// logger.ErrorString("短信[腾讯云]", "调用接口错误", err.Error())
		return "", nil
	}
	if err != nil {
		// logger.ErrorString("短信[腾讯云]", "服务商返回错误", err.Error())
		// return "", nil
		panic(err)
	}

	// {"PhoneNumberSet":["15618903080"],"SmsSdkAppId":"1400874085","TemplateId":"2017731","SignName":"武汉搭个线网络科技","TemplateParamSet":["987654"]}
	// {"Response":{"SendStatusSet":[{"SerialNo":"3369:15051191117040689090090308","PhoneNumber":"+8615618903080","Fee":1,"SessionContext":"","Code":"Ok","Message":"send success","IsoCode":"CN"}],"RequestId":"3ea30583-fa02-4ac1-9564-f47846856d18"}}
	return *response.Response.RequestId, nil
}
