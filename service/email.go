package service

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"net/url"
	"os"
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/database"
)

type emailService struct {
    smtpHost string
    smtpPort string
    from     string
    password string
    appName  string
    appURL   string
}

func ImplEmailService() *emailService {
    return &emailService{
        smtpHost: getEnv("SMTP_HOST", "smtp.gmail.com"),
        smtpPort: getEnv("SMTP_PORT", "587"),
        from:     getEnv("SMTP_FROM", "noreply@studiobooking.com"),
        password: getEnv("SMTP_PASSWORD", ""),
        appName:  getEnv("APP_NAME", "Studio Booking System"),
        appURL:   getEnv("APP_URL", "http://localhost:3000"),
    }
}

// SendBookingCreated - Notify customer booking created (pending payment via WhatsApp)
func (s *emailService) SendBookingCreated(booking *database.Booking) error {
    if booking.User == nil || booking.Studio == nil {
        return fmt.Errorf("booking missing user or studio relation")
    }

    subject := "Booking Created - Please Contact Admin for Payment"

    // Get admin WhatsApp from environment
    adminWhatsApp := getEnv("ADMIN_WHATSAPP", "+6289570608111")
    adminWhatsAppDisplay := getEnv("ADMIN_WHATSAPP_DISPLAY", "0895-7060-8111")
    adminName := getEnv("ADMIN_NAME", "Admin Studio Booking")

    // Generate WhatsApp message
    whatsappMessage := fmt.Sprintf(
        "Halo %s,\n\n"+
            "Saya ingin konfirmasi pembayaran untuk booking:\n\n"+
            "📋 *Booking ID:* #%d\n"+
            "🎵 *Studio:* %s\n"+
            "📅 *Tanggal:* %s\n"+
            "⏰ *Waktu:* %s - %s\n"+
            "💰 *Total Pembayaran:* Rp %s\n\n"+
            "Mohon informasi cara pembayarannya. Terima kasih!",
        adminName,
        booking.ID,
        booking.Studio.Name,
        booking.BookingDate.Format("02 January 2006"),
        booking.StartTime.Format("15:04"),
        booking.EndTime.Format("15:04"),
        formatNumber(booking.TotalPrice),
    )

    // Create WhatsApp link
    whatsappLink := fmt.Sprintf(
        "https://wa.me/%s?text=%s",
        adminWhatsApp,
        url.QueryEscape(whatsappMessage),
    )

    data := map[string]interface{}{
        "CustomerName":          booking.User.Name,
        "BookingID":             booking.ID,
        "StudioName":            booking.Studio.Name,
        "BookingDate":           booking.BookingDate.Format("Monday, 02 January 2006"),
        "StartTime":             booking.StartTime.Format("15:04"),
        "EndTime":               booking.EndTime.Format("15:04"),
        "TotalPrice":            formatCurrency(booking.TotalPrice),
        "AdminName":             adminName,
        "AdminWhatsApp":         adminWhatsAppDisplay,
        "AdminWhatsAppFull":     adminWhatsApp,
        "WhatsAppLink":          whatsappLink,
        "AppName":               s.appName,
        "AppURL":                s.appURL,
        "Year":                  time.Now().Year(),
    }

    body, err := s.renderTemplate("booking_created", data)
    if err != nil {
        return err
    }

    return s.sendEmail(booking.User.Email, subject, body)
}

// SendBookingConfirmed - Notify customer booking confirmed by admin
func (s *emailService) SendBookingConfirmed(booking *database.Booking) error {
    if booking.User == nil || booking.Studio == nil {
        return fmt.Errorf("booking missing user or studio relation")
    }

    subject := "Booking Confirmed! ✅"

    adminNotes := booking.AdminNotes
    if adminNotes == "" {
        adminNotes = "Payment has been received and verified."
    }

    data := map[string]interface{}{
        "CustomerName": booking.User.Name,
        "BookingID":    booking.ID,
        "StudioName":   booking.Studio.Name,
        "BookingDate":  booking.BookingDate.Format("Monday, 02 January 2006"),
        "StartTime":    booking.StartTime.Format("15:04"),
        "EndTime":      booking.EndTime.Format("15:04"),
        "TotalPrice":   formatCurrency(booking.TotalPrice),
        "AdminNotes":   adminNotes,
        "AppName":      s.appName,
        "AppURL":       s.appURL,
        "Year":         time.Now().Year(),
    }

    body, err := s.renderTemplate("booking_confirmed", data)
    if err != nil {
        return err
    }

    return s.sendEmail(booking.User.Email, subject, body)
}

// SendBookingCancelled - Notify customer booking cancelled
func (s *emailService) SendBookingCancelled(booking *database.Booking, reason string) error {
    if booking.User == nil || booking.Studio == nil {
        return fmt.Errorf("booking missing user or studio relation")
    }

    subject := "Booking Cancelled"

    if reason == "" {
        reason = "Your booking has been cancelled by admin."
    }

    data := map[string]interface{}{
        "CustomerName": booking.User.Name,
        "BookingID":    booking.ID,
        "StudioName":   booking.Studio.Name,
        "BookingDate":  booking.BookingDate.Format("Monday, 02 January 2006"),
        "StartTime":    booking.StartTime.Format("15:04"),
        "EndTime":      booking.EndTime.Format("15:04"),
        "Reason":       reason,
        "AppName":      s.appName,
        "AppURL":       s.appURL,
        "Year":         time.Now().Year(),
    }

    body, err := s.renderTemplate("booking_cancelled", data)
    if err != nil {
        return err
    }

    return s.sendEmail(booking.User.Email, subject, body)
}

// sendEmail - Send email via SMTP
func (s *emailService) sendEmail(to, subject, body string) error {
    if s.smtpHost == "" || s.smtpPort == "" || s.from == "" {
        log.Printf("⚠️  SMTP not configured, skipping email to %s", to)
        return nil
    }

    if s.password == "" {
        log.Printf("⚠️  SMTP password not set, skipping email to %s", to)
        return nil
    }

    auth := smtp.PlainAuth("", s.from, s.password, s.smtpHost)

    headers := make(map[string]string)
    headers["From"] = s.from
    headers["To"] = to
    headers["Subject"] = subject
    headers["MIME-Version"] = "1.0"
    headers["Content-Type"] = "text/html; charset=UTF-8"

    message := ""
    for k, v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
    }
    message += "\r\n" + body

    addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)
    err := smtp.SendMail(addr, auth, s.from, []string{to}, []byte(message))

    if err != nil {
        return fmt.Errorf("failed to send email to %s: %w", to, err)
    }

    log.Printf("✅ Email sent to %s: %s", to, subject)
    return nil
}

// renderTemplate - Render HTML email template
func (s *emailService) renderTemplate(templateName string, data map[string]interface{}) (string, error) {
    tmpl, err := template.New(templateName).Parse(getEmailTemplate(templateName))
    if err != nil {
        return "", err
    }

    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }

    return buf.String(), nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}

func formatCurrency(amount int) string {
    return fmt.Sprintf("Rp %s", formatNumber(amount))
}

func formatNumber(n int) string {
    s := fmt.Sprintf("%d", n)
    result := ""
    for i, c := range s {
        if i > 0 && (len(s)-i)%3 == 0 {
            result += "."
        }
        result += string(c)
    }
    return result
}

// getEmailTemplate - Get HTML email template
func getEmailTemplate(name string) string {
    templates := map[string]string{
        "booking_created": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; background: #f4f4f4; }
        .container { max-width: 600px; margin: 20px auto; background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; }
        .content { padding: 30px; }
        .booking-card { background: #f8f9fa; border-left: 4px solid #667eea; padding: 20px; margin: 20px 0; border-radius: 5px; }
        .detail-row { display: flex; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #e9ecef; }
        .detail-row:last-child { border-bottom: none; }
        .label { font-weight: 600; color: #495057; }
        .value { color: #212529; }
        .total-price { font-size: 24px; color: #667eea; font-weight: bold; }
        .whatsapp-card { background: linear-gradient(135deg, #25D366 0%, #128C7E 100%); color: white; padding: 25px; margin: 20px 0; border-radius: 10px; text-align: center; }
        .whatsapp-btn { display: inline-block; background: white; color: #25D366; padding: 14px 30px; text-decoration: none; border-radius: 25px; margin: 15px 0; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .whatsapp-btn:hover { background: #f0f0f0; }
        .alert-box { background: #fff3cd; border-left: 4px solid #ffc107; padding: 15px; margin: 20px 0; border-radius: 5px; }
        .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎵 Booking Created!</h1>
            <p style="margin: 10px 0 0 0; opacity: 0.9;">Your studio booking has been created</p>
        </div>
        
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            <p>Thank you for booking with us! Your booking has been successfully created.</p>
            
            <div class="booking-card">
                <h3 style="margin-top: 0; color: #667eea;">📋 Booking Details</h3>
                <div class="detail-row">
                    <span class="label">Booking ID</span>
                    <span class="value"><strong>#{{.BookingID}}</strong></span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Time</span>
                    <span class="value">{{.StartTime}} - {{.EndTime}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Total Price</span>
                    <span class="total-price">{{.TotalPrice}}</span>
                </div>
            </div>

            <div class="whatsapp-card">
                <h3 style="margin: 0 0 10px 0; font-size: 22px;">💬 Contact Admin for Payment</h3>
                <p style="margin: 10px 0; font-size: 18px; font-weight: 600;">{{.AdminName}}</p>
                <p style="margin: 5px 0; font-size: 20px;">📱 {{.AdminWhatsApp}}</p>
                <a href="{{.WhatsAppLink}}" class="whatsapp-btn">
                    💬 Chat on WhatsApp
                </a>
                <p style="margin: 15px 0 0 0; font-size: 14px; opacity: 0.9;">
                    Click the button above to send payment confirmation directly
                </p>
            </div>

            <div class="alert-box">
                <strong>⚠️ Important:</strong><br>
                Your booking is currently in <strong>PENDING</strong> status. Please contact admin via WhatsApp for payment instructions.
                After payment is verified, your status will be updated to <strong>CONFIRMED</strong>.
            </div>

            <p style="margin-top: 30px;">Thank you for choosing <strong>{{.AppName}}</strong>!</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

        "booking_confirmed": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; background: #f4f4f4; }
        .container { max-width: 600px; margin: 20px auto; background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #10b981 0%, #059669 100%); color: white; padding: 30px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; }
        .content { padding: 30px; }
        .success-badge { background: #d1fae5; color: #065f46; padding: 15px 25px; border-radius: 25px; display: inline-block; margin: 20px 0; font-weight: 600; font-size: 16px; }
        .booking-card { background: #f8f9fa; border-left: 4px solid #10b981; padding: 20px; margin: 20px 0; border-radius: 5px; }
        .detail-row { display: flex; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #e9ecef; }
        .detail-row:last-child { border-bottom: none; }
        .label { font-weight: 600; color: #495057; }
        .value { color: #212529; }
        .info-box { background: #dbeafe; border-left: 4px solid #3b82f6; padding: 15px; margin: 20px 0; border-radius: 5px; }
        .btn { display: inline-block; background: #10b981; color: white; padding: 14px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; font-weight: 600; }
        .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Booking Confirmed!</h1>
            <p style="margin: 10px 0 0 0; opacity: 0.9;">Your payment has been verified</p>
        </div>
        
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            
            <center>
                <div class="success-badge">
                    🎉 Your booking has been confirmed!
                </div>
            </center>

            <p>Great news! Your studio booking has been confirmed by our admin. You're all set!</p>
            
            <div class="booking-card">
                <h3 style="margin-top: 0; color: #10b981;">📋 Booking Details</h3>
                <div class="detail-row">
                    <span class="label">Booking ID</span>
                    <span class="value"><strong>#{{.BookingID}}</strong></span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Time</span>
                    <span class="value">{{.StartTime}} - {{.EndTime}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Total Price</span>
                    <span class="value">{{.TotalPrice}}</span>
                </div>
            </div>

            <div class="info-box">
                <strong>💼 Admin Notes:</strong><br>
                {{.AdminNotes}}
            </div>

            <p style="margin-top: 30px;"><strong>Important Reminders:</strong></p>
            <ul>
                <li>Please arrive 10 minutes before your scheduled time</li>
                <li>Bring a valid ID for verification</li>
                <li>Contact us if you need to reschedule</li>
            </ul>

            <p>See you at the studio! 🎵</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

        "booking_cancelled": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; background: #f4f4f4; }
        .container { max-width: 600px; margin: 20px auto; background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); color: white; padding: 30px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; }
        .content { padding: 30px; }
        .booking-card { background: #f8f9fa; border-left: 4px solid #ef4444; padding: 20px; margin: 20px 0; border-radius: 5px; }
        .detail-row { display: flex; justify-content: space-between; padding: 12px 0; border-bottom: 1px solid #e9ecef; }
        .detail-row:last-child { border-bottom: none; }
        .label { font-weight: 600; color: #495057; }
        .value { color: #212529; }
        .reason-box { background: #fee2e2; border-left: 4px solid #ef4444; padding: 15px; margin: 20px 0; border-radius: 5px; color: #991b1b; }
        .btn { display: inline-block; background: #667eea; color: white; padding: 14px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; font-weight: 600; }
        .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Booking Cancelled</h1>
            <p style="margin: 10px 0 0 0; opacity: 0.9;">Your booking has been cancelled</p>
        </div>
        
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            <p>We're writing to inform you that your booking has been cancelled.</p>
            
            <div class="booking-card">
                <h3 style="margin-top: 0; color: #ef4444;">📋 Cancelled Booking</h3>
                <div class="detail-row">
                    <span class="label">Booking ID</span>
                    <span class="value"><strong>#{{.BookingID}}</strong></span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Time</span>
                    <span class="value">{{.StartTime}} - {{.EndTime}}</span>
                </div>
            </div>

            <div class="reason-box">
                <strong>Cancellation Reason:</strong><br>
                {{.Reason}}
            </div>

            <center>
                <a href="{{.AppURL}}/studios" class="btn">Browse Other Studios</a>
            </center>

            <p style="margin-top: 30px;">If you have any questions or concerns, please don't hesitate to contact our support team.</p>
        </div>
        
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,
    }

    return templates[name]
}