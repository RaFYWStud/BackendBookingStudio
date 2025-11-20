package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"time"

	"github.com/RaFYWStud/BackendBookingStudio/database"
)

type EmailService interface {
    SendBookingCreated(booking *database.Booking) error
    SendBookingConfirmed(booking *database.Booking) error
    SendBookingCancelled(booking *database.Booking, reason string) error
    SendPaymentVerified(payment *database.Payment) error
    SendPaymentRejected(payment *database.Payment, reason string) error
}

type emailService struct {
    smtpHost string
    smtpPort string
    from     string
    password string
    appName  string
    appURL   string
}

func ImplEmailService() EmailService {
    return &emailService{
        smtpHost: os.Getenv("SMTP_HOST"),
        smtpPort: os.Getenv("SMTP_PORT"),
        from:     os.Getenv("SMTP_FROM"),
        password: os.Getenv("SMTP_PASSWORD"),
        appName:  getEnv("APP_NAME", "Studio Booking System"),
        appURL:   getEnv("APP_URL", "http://localhost:3000"),
    }
}

// SendBookingCreated - Notify customer booking created (pending payment)
func (s *emailService) SendBookingCreated(booking *database.Booking) error {
    if booking.User == nil || booking.Studio == nil {
        return fmt.Errorf("booking missing user or studio relation")
    }

    subject := "Booking Created - Please Complete Payment"

    data := map[string]interface{}{
        "CustomerName": booking.User.Name,
        "BookingID":    booking.ID,
        "StudioName":   booking.Studio.Name,
        "BookingDate":  booking.BookingDate.Format("Monday, 02 January 2006"),
        "StartTime":    booking.StartTime.Format("15:04"),
        "EndTime":      booking.EndTime.Format("15:04"),
        "TotalPrice":   formatCurrency(booking.TotalPrice),
        "DPAmount":     formatCurrency(booking.DPAmount),
        "DPDeadline":   booking.DPDeadline.Format("02 Jan 2006 15:04"),
        "AppName":      s.appName,
        "AppURL":       s.appURL,
        "Year":         time.Now().Year(),
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

    data := map[string]interface{}{
        "CustomerName":    booking.User.Name,
        "BookingID":       booking.ID,
        "StudioName":      booking.Studio.Name,
        "BookingDate":     booking.BookingDate.Format("Monday, 02 January 2006"),
        "StartTime":       booking.StartTime.Format("15:04"),
        "EndTime":         booking.EndTime.Format("15:04"),
        "RemainingAmount": formatCurrency(booking.RemainingAmount),
        "AppName":         s.appName,
        "AppURL":          s.appURL,
        "Year":            time.Now().Year(),
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

    data := map[string]interface{}{
        "CustomerName": booking.User.Name,
        "BookingID":    booking.ID,
        "StudioName":   booking.Studio.Name,
        "BookingDate":  booking.BookingDate.Format("Monday, 02 January 2006"),
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

// SendPaymentVerified - Notify customer payment verified
func (s *emailService) SendPaymentVerified(payment *database.Payment) error {
    if payment.Booking == nil || payment.Booking.User == nil || payment.Booking.Studio == nil {
        return fmt.Errorf("payment missing booking relations")
    }

    booking := payment.Booking
    subject := "Payment Verified ✅"

    paymentType := "Down Payment (DP)"
    if payment.PaymentType == "full_payment" {
        paymentType = "Full Payment"
    }

    data := map[string]interface{}{
        "CustomerName":  booking.User.Name,
        "PaymentType":   paymentType,
        "Amount":        formatCurrency(payment.Amount),
        "BookingID":     booking.ID,
        "StudioName":    booking.Studio.Name,
        "BookingDate":   booking.BookingDate.Format("Monday, 02 January 2006"),
        "BookingStatus": string(booking.Status),
        "AppName":       s.appName,
        "AppURL":        s.appURL,
        "Year":          time.Now().Year(),
    }

    body, err := s.renderTemplate("payment_verified", data)
    if err != nil {
        return err
    }

    return s.sendEmail(booking.User.Email, subject, body)
}

// SendPaymentRejected - Notify customer payment rejected
func (s *emailService) SendPaymentRejected(payment *database.Payment, reason string) error {
    if payment.Booking == nil || payment.Booking.User == nil || payment.Booking.Studio == nil {
        return fmt.Errorf("payment missing booking relations")
    }

    booking := payment.Booking
    subject := "Payment Rejected - Action Required"

    paymentType := "Down Payment (DP)"
    if payment.PaymentType == "full_payment" {
        paymentType = "Full Payment"
    }

    data := map[string]interface{}{
        "CustomerName": booking.User.Name,
        "PaymentType":  paymentType,
        "Amount":       formatCurrency(payment.Amount),
        "BookingID":    booking.ID,
        "StudioName":   booking.Studio.Name,
        "Reason":       reason,
        "AppName":      s.appName,
        "AppURL":       s.appURL,
        "Year":         time.Now().Year(),
    }

    body, err := s.renderTemplate("payment_rejected", data)
    if err != nil {
        return err
    }

    return s.sendEmail(booking.User.Email, subject, body)
}

// sendEmail - Send email via SMTP
func (s *emailService) sendEmail(to, subject, body string) error {
    if s.smtpHost == "" || s.smtpPort == "" || s.from == "" || s.password == "" {
        return fmt.Errorf("SMTP configuration incomplete. Please check .env file")
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

    fmt.Printf("✅ Email sent to %s: %s\n", to, subject)
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
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4F46E5; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .booking-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #e5e7eb; }
        .label { font-weight: bold; color: #6B7280; }
        .value { color: #111827; }
        .alert { background: #FEF3C7; border-left: 4px solid #F59E0B; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; background: #4F46E5; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; color: #6B7280; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎵 Booking Created!</h1>
        </div>
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            <p>Your studio booking has been created successfully!</p>
            
            <div class="booking-details">
                <h3>Booking Details</h3>
                <div class="detail-row">
                    <span class="label">Booking ID:</span>
                    <span class="value">#{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio:</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date:</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Time:</span>
                    <span class="value">{{.StartTime}} - {{.EndTime}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Total Price:</span>
                    <span class="value">{{.TotalPrice}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">DP Amount:</span>
                    <span class="value"><strong>{{.DPAmount}}</strong></span>
                </div>
            </div>

            <div class="alert">
                <strong>⏰ Important!</strong><br>
                Please complete your DP payment before <strong>{{.DPDeadline}}</strong> to secure your booking.
                Your booking will be automatically cancelled if payment is not received by the deadline.
            </div>

            <center>
                <a href="{{.AppURL}}/bookings/{{.BookingID}}" class="button">View Booking Details</a>
            </center>

            <p>Thank you for choosing {{.AppName}}!</p>
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
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #10B981; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .success-badge { background: #D1FAE5; color: #065F46; padding: 10px 20px; border-radius: 20px; display: inline-block; margin: 20px 0; }
        .booking-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #e5e7eb; }
        .label { font-weight: bold; color: #6B7280; }
        .value { color: #111827; }
        .info-box { background: #DBEAFE; border-left: 4px solid #3B82F6; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; background: #10B981; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; color: #6B7280; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>✅ Booking Confirmed!</h1>
        </div>
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            
            <center>
                <div class="success-badge">
                    <strong>🎉 Your booking has been confirmed!</strong>
                </div>
            </center>

            <p>Great news! Your studio booking has been approved by our admin.</p>
            
            <div class="booking-details">
                <h3>Booking Details</h3>
                <div class="detail-row">
                    <span class="label">Booking ID:</span>
                    <span class="value">#{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio:</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date:</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Time:</span>
                    <span class="value">{{.StartTime}} - {{.EndTime}}</span>
                </div>
            </div>

            <div class="info-box">
                <strong>💰 Remaining Payment</strong><br>
                Please complete the remaining payment of <strong>{{.RemainingAmount}}</strong> before your booking date.
            </div>

            <center>
                <a href="{{.AppURL}}/bookings/{{.BookingID}}" class="button">View Booking</a>
            </center>

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
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #EF4444; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .booking-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #e5e7eb; }
        .label { font-weight: bold; color: #6B7280; }
        .value { color: #111827; }
        .reason-box { background: #FEE2E2; border-left: 4px solid #EF4444; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; background: #4F46E5; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; color: #6B7280; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Booking Cancelled</h1>
        </div>
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            <p>Your booking has been cancelled.</p>
            
            <div class="booking-details">
                <h3>Cancelled Booking</h3>
                <div class="detail-row">
                    <span class="label">Booking ID:</span>
                    <span class="value">#{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio:</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date:</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
            </div>

            <div class="reason-box">
                <strong>Cancellation Reason:</strong><br>
                {{.Reason}}
            </div>

            <center>
                <a href="{{.AppURL}}/studios" class="button">Browse Studios</a>
            </center>

            <p>If you have any questions, please contact our support team.</p>
        </div>
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

        "payment_verified": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #10B981; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .success-badge { background: #D1FAE5; color: #065F46; padding: 10px 20px; border-radius: 20px; display: inline-block; margin: 20px 0; font-size: 18px; }
        .payment-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #e5e7eb; }
        .label { font-weight: bold; color: #6B7280; }
        .value { color: #111827; }
        .amount { font-size: 24px; color: #10B981; font-weight: bold; }
        .button { display: inline-block; background: #10B981; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; color: #6B7280; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>💸 Payment Verified!</h1>
        </div>
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            
            <center>
                <div class="success-badge">
                    ✅ Your payment has been verified!
                </div>
            </center>

            <p>Your {{.PaymentType}} payment has been successfully verified by our admin.</p>
            
            <div class="payment-details">
                <h3>Payment Details</h3>
                <div class="detail-row">
                    <span class="label">Payment Type:</span>
                    <span class="value">{{.PaymentType}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Amount:</span>
                    <span class="amount">{{.Amount}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Booking ID:</span>
                    <span class="value">#{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio:</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Date:</span>
                    <span class="value">{{.BookingDate}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Booking Status:</span>
                    <span class="value"><strong>{{.BookingStatus}}</strong></span>
                </div>
            </div>

            <center>
                <a href="{{.AppURL}}/bookings/{{.BookingID}}" class="button">View Booking</a>
            </center>

            <p>Thank you for your payment! We look forward to seeing you at the studio. 🎵</p>
        </div>
        <div class="footer">
            <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`,

        "payment_rejected": `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #EF4444; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .payment-details { background: white; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .detail-row { display: flex; justify-content: space-between; padding: 10px 0; border-bottom: 1px solid #e5e7eb; }
        .label { font-weight: bold; color: #6B7280; }
        .value { color: #111827; }
        .alert { background: #FEE2E2; border-left: 4px solid #EF4444; padding: 15px; margin: 20px 0; border-radius: 4px; }
        .button { display: inline-block; background: #EF4444; color: white; padding: 12px 30px; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; color: #6B7280; font-size: 12px; margin-top: 30px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>❌ Payment Rejected</h1>
        </div>
        <div class="content">
            <p>Hi <strong>{{.CustomerName}}</strong>,</p>
            <p>Unfortunately, your {{.PaymentType}} payment has been rejected.</p>
            
            <div class="payment-details">
                <h3>Payment Details</h3>
                <div class="detail-row">
                    <span class="label">Payment Type:</span>
                    <span class="value">{{.PaymentType}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Amount:</span>
                    <span class="value">{{.Amount}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Booking ID:</span>
                    <span class="value">#{{.BookingID}}</span>
                </div>
                <div class="detail-row">
                    <span class="label">Studio:</span>
                    <span class="value">{{.StudioName}}</span>
                </div>
            </div>

            <div class="alert">
                <strong>Rejection Reason:</strong><br>
                {{.Reason}}
            </div>

            <p><strong>What to do next:</strong></p>
            <ul>
                <li>Check the rejection reason above</li>
                <li>Upload a new payment proof with correct information</li>
                <li>Make sure the payment amount matches the required amount</li>
                <li>Use clear and readable payment proof image</li>
            </ul>

            <center>
                <a href="{{.AppURL}}/bookings/{{.BookingID}}" class="button">Upload New Payment</a>
            </center>

            <p>If you have any questions, please contact our support team.</p>
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