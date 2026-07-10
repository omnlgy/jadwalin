package domain

import "context"

type WhatsAppProvider interface {
	SendMessage(ctx context.Context, to, message string) error
}

type EmailProvider interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type NotificationService interface {
	SendOTPLoginWhatsApp(ctx context.Context, to, otpCode string) error
	SendEmail(ctx context.Context, to, subject, body string) error
	SendEmailInVoice(ctx context.Context, booking *Booking) error
}
