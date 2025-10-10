package auth

import (
    "fmt"
    "os"
    "strconv"
    "gopkg.in/mail.v2"
)

// SendVerificationEmail sends a verification email via Mailtrap SMTP with a token link
func SendVerificationEmail(to string, link string) error {
    host := getenvDefault("MAILTRAP_HOST", "live.smtp.mailtrap.io")
    portStr := getenvDefault("MAILTRAP_PORT", "587")
    username := os.Getenv("MAILTRAP_USERNAME")
    password := os.Getenv("MAILTRAP_PASSWORD")
    from := getenvDefault("MAIL_FROM", "no-reply@example.com")

    port, err := strconv.Atoi(portStr)
    if err != nil {
        return fmt.Errorf("invalid MAILTRAP_PORT: %w", err)
    }

    m := mail.NewMessage()
    m.SetHeader("From", from)
    m.SetHeader("To", to)
    m.SetHeader("Subject", "Verify your email")
    m.SetBody("text/plain", fmt.Sprintf("Welcome! Please verify your email by visiting: %s", link))

    d := mail.NewDialer(host, port, username, password)
    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }
    return nil
}

func getenvDefault(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}
