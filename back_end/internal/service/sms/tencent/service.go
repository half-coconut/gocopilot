package tencent

import sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

type Service struct {
	sms      *sms.Client
	appId    *string
	signName *string
}

func NewService(sms *sms.Client, appId, signName *string) *Service {
	return &Service{
		sms:      sms,
		appId:    appId,
		signName: signName,
	}
}

func Send() {

}
