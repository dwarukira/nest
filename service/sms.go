package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/solabsafrica/afrikanest/config"
	"github.com/solabsafrica/afrikanest/logger"
)

type SmsServiceWithContext func(ctx context.Context) SmsService

type SmsService interface {
	Send(to string, message string) error
}

type smsServiceImpl struct {
	ctx    context.Context
	config *config.Config
}

func (smsService *smsServiceImpl) Send(to string, message string) error {
	// service := sms.NewService(smsService.config.SmsConfig.Username, smsService.config.SmsConfig.ApiKey, smsService.config.SmsConfig.Env)

	// recipients, err := service.Send(to, message, "AFRIKANEST")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(recipients, to)
	postBody, _ := json.Marshal(map[string]string{
		"username": smsService.config.SmsConfig.Username,
		"to":       to,
		"message":  message,
		"from":     "AFRIKANEST",
	})
	logger.Errorf("Hello %s", postBody)

	// responseBody := bytes.NewBuffer(postBody)
	values := url.Values{}
	values.Set("username", smsService.config.SmsConfig.Username)
	values.Set("to", to)
	values.Set("message", message)
	values.Set("from", "AFRIKANEST")

	reader := strings.NewReader(values.Encode())

	req, err := http.NewRequest("POST", "https://api.africastalking.com/version1/messaging", reader)
	req.Header.Set("apiKey", smsService.config.SmsConfig.ApiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(reader.Len()))
	req.Header.Set("Accept", "application/json")
	if err != nil {
		logger.Errorf("failed to send sms %v", err)
		return err
	}
	client := &http.Client{}
	response, err := client.Do(req)
	logger.Info(response)
	if err != nil {
		logger.Errorf("failed to send sms %v", err)
	}
	defer response.Body.Close()
	return nil
}

func NewSmsServiceWithContext() SmsServiceWithContext {
	config := config.Get()
	return func(ctx context.Context) SmsService {
		return &smsServiceImpl{
			config: config,
			ctx:    ctx,
		}
	}
}
