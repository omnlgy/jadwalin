package service

import (
	"context"
	"fmt"

	"github.com/omnlgy/jadwalin/internal/domain"
)

type NotificationService struct {
	waProvider    domain.WhatsAppProvider
	emailProvider domain.EmailProvider
}

func NewNotificationService(waProvider domain.WhatsAppProvider, emailProvider domain.EmailProvider) domain.NotificationService {
	return &NotificationService{
		waProvider:    waProvider,
		emailProvider: emailProvider,
	}
}

func (s *NotificationService) SendOTPLoginWhatsApp(ctx context.Context, to, otpCode string) error {
	message := fmt.Sprintf("⏱️ *Login Request*\nYour OTP is: *%s*\n\n⏳ This code expires in _5 minutes_.\n_If you did not request this, please ignore this message._", otpCode)
	return s.waProvider.SendMessage(ctx, to, message)
}
