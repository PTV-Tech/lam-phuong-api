package email

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// Service handles email sending via SMTP relay
type Service struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
	fromName     string
	useTLS       bool // Use TLS for SMTP connection
}

// NewService creates a new email service with TLS enabled by default
func NewService(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, fromName string) *Service {
	return NewServiceWithTLS(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, fromName, true)
}

// NewServiceWithTLS creates a new email service with configurable TLS
func NewServiceWithTLS(smtpHost, smtpPort, smtpUsername, smtpPassword, fromEmail, fromName string, useTLS bool) *Service {
	return &Service{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
		fromName:     fromName,
		useTLS:       useTLS,
	}
}

// SendVerificationEmail sends an email verification email to the user
func (s *Service) SendVerificationEmail(toEmail, verificationToken, baseURL string) error {
	verificationURL := fmt.Sprintf("%s/api/auth/verify-email?token=%s", baseURL, verificationToken)

	subject := "Verify Your Email Address"
	body := fmt.Sprintf(`Hello,

Thank you for registering! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you did not create an account, please ignore this email.

Best regards,
%s`, verificationURL, s.fromName)

	return s.sendEmail(toEmail, subject, body)
}

// sendEmail sends an email using SMTP relay
// Authentication is optional - works with open relays or authenticated SMTP servers
func (s *Service) sendEmail(toEmail, subject, body string) error {
	// If SMTP is not configured, log and skip sending (for development)
	if s.smtpHost == "" || s.smtpPort == "" {
		fmt.Printf("[EMAIL] Would send email to %s\n", toEmail)
		fmt.Printf("[EMAIL] Subject: %s\n", subject)
		fmt.Printf("[EMAIL] Body: %s\n", body)
		return nil
	}

	// Validate email addresses
	if !isValidEmail(toEmail) {
		return fmt.Errorf("invalid recipient email address: %s", toEmail)
	}
	if !isValidEmail(s.fromEmail) {
		return fmt.Errorf("invalid sender email address: %s", s.fromEmail)
	}

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	// Create email message with proper headers
	from := s.fromEmail
	if s.fromName != "" {
		from = fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail)
	}

	// Build email message with proper headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "8bit"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	msg := []byte(message)

	// Connect to SMTP server
	var client *smtp.Client
	var err error

	// Determine connection type based on port and TLS setting
	// Port 465: Direct TLS (implicit TLS)
	// Port 587/25: STARTTLS (plain connection then upgrade)
	// Other ports: Use TLS setting to determine

	port := s.smtpPort
	if port == "" {
		port = "25"
	}

	// Port 465 uses direct TLS connection
	if port == "465" && s.useTLS {
		// Use direct TLS connection for port 465
		tlsConfig := &tls.Config{
			ServerName: s.smtpHost,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, s.smtpHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
	} else {
		// Use plain connection first (for ports 25, 587, or when TLS is disabled)
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer conn.Close()

		client, err = smtp.NewClient(conn, s.smtpHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}

		// Upgrade to STARTTLS if TLS is enabled and server supports it
		if s.useTLS {
			if ok, _ := client.Extension("STARTTLS"); ok {
				tlsConfig := &tls.Config{
					ServerName: s.smtpHost,
				}
				if err = client.StartTLS(tlsConfig); err != nil {
					return fmt.Errorf("failed to start TLS: %w", err)
				}
			} else {
				// Server doesn't support STARTTLS, but TLS was requested
				// This is OK - continue without TLS upgrade
			}
		}
	}

	defer client.Close()

	// Authenticate if credentials are provided
	// Only authenticate if server supports AUTH extension
	if s.smtpUsername != "" && s.smtpPassword != "" {
		// Check if server supports authentication
		if ok, authMethods := client.Extension("AUTH"); ok {
			// Check if PLAIN auth is supported
			supportsPlain := false
			if authMethods != "" {
				// authMethods might be something like "PLAIN LOGIN" or "PLAIN"
				if strings.Contains(strings.ToUpper(authMethods), "PLAIN") {
					supportsPlain = true
				}
			} else {
				// If no methods listed, assume PLAIN is supported
				supportsPlain = true
			}

			if supportsPlain {
				auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)
				if err = client.Auth(auth); err != nil {
					// Authentication failed
					// EOF error usually means server closed connection - might not support auth on plain connection
					// or credentials are wrong
					return fmt.Errorf("SMTP authentication failed: %w", err)
				}
			}
		}
		// If server doesn't support AUTH extension, continue without authentication (open relay)
	}

	// Set sender
	if err = client.Mail(s.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email data
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		writer.Close()
		return fmt.Errorf("failed to write email data: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	// Quit
	err = client.Quit()
	if err != nil {
		return fmt.Errorf("failed to quit SMTP session: %w", err)
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if parts[0] == "" || parts[1] == "" {
		return false
	}
	return true
}

// GenerateVerificationToken generates a secure random token for email verification
func GenerateVerificationToken() (string, error) {
	b := make([]byte, 32) // 64 character hex string
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
