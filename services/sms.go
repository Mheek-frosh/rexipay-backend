package services

import (
	"fmt"
)

type SMSService interface {
	SendOTP(phoneNumber, code string) error
}

type MockSMSService struct{}

func (s *MockSMSService) SendOTP(phoneNumber, code string) error {
	fmt.Printf(" [MOCK SMS] Sending OTP %s to %s\n", code, phoneNumber)
	return nil
}

var DefaultSMSService SMSService = &MockSMSService{}
