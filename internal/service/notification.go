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
	message := fmt.Sprintf("Your OTP is: %s", otpCode)
	return s.waProvider.SendMessage(ctx, to, message)
}

func (s *NotificationService) SendEmail(ctx context.Context, to, subject, body string) error {
	return s.emailProvider.SendEmail(ctx, to, subject, body)
}

func (s *NotificationService) SendEmailInVoice(ctx context.Context, booking *domain.Booking) error {
	_, _ = fmt.Printf("DEBUG: booking.Client=%+v\n", booking.Client)
	_, _ = fmt.Printf("DEBUG: booking.Treatment=%+v\n", booking.Treatment)
	_, _ = fmt.Printf("DEBUG: booking.Staff=%+v\n", booking.Staff)
	emailBody := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
</head>
<body>
<h2>Konfirmasi & Invoice Booking Treatment</h2>
<p>Halo <strong>%s</strong>,</p>
<p>Terima kasih telah melakukan reservasi treatment di tempat kami. Berikut adalah detail booking Anda:</p>
<table>
<tr><td>Nama Treatment</td><td>%s</td></tr>
<tr><td>Jadwal Booking</td><td>%s</td></tr>
<tr><td>Staff Penangan</td><td>%s</td></tr>
<tr><td>Total Harga</td><td>Rp %d</td></tr>
</table>
<p>Terima kasih,<br><strong>Jadwalin</strong></p>
</body>
</html>`, booking.Client.FullName, booking.Treatment.Name, booking.StartTime.Format("2006-01-02 15:04"), booking.Staff.FullName, booking.Treatment.Price)
	return s.emailProvider.SendEmail(ctx, booking.Client.Email, "Konfirmasi & Invoice Booking Treatment", emailBody)
}
