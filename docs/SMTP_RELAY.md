# SMTP Relay Email Service

This document describes the SMTP relay email service implementation for sending emails.

## Overview

The email service uses SMTP relay to send emails with support for:
- TLS/SSL encryption
- STARTTLS upgrade
- SMTP authentication (PLAIN auth)
- Proper email headers
- Email address validation

## Configuration

### Environment Variables

```bash
# SMTP Server Configuration
EMAIL_SMTP_HOST=smtp.gmail.com          # SMTP server hostname
EMAIL_SMTP_PORT=587                     # SMTP server port (587 for STARTTLS, 465 for TLS)
EMAIL_SMTP_USERNAME=your-email@gmail.com  # Optional - leave empty for open relays
EMAIL_SMTP_PASSWORD=your-app-password   # Optional - leave empty for open relays

# Email Settings
EMAIL_FROM_EMAIL=noreply@lamphuong.com  # From email address
EMAIL_FROM_NAME=Lam Phuong             # From name
EMAIL_BASE_URL=http://localhost:8080   # Base URL for verification links
EMAIL_USE_TLS=true                      # Use TLS (default: true)
```

### Authentication

**SMTP authentication is optional**. The service will:
- **With credentials**: Authenticate using PLAIN auth if username and password are provided
- **Without credentials**: Send emails without authentication (for open relays or internal SMTP servers)

This is useful for:
- Internal SMTP relays that don't require authentication
- Open SMTP relays (use with caution)
- Testing environments

### Common SMTP Providers

#### Gmail
```bash
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-app-password  # Generate app password in Google Account settings
EMAIL_USE_TLS=true
```

#### Outlook/Office 365
```bash
EMAIL_SMTP_HOST=smtp.office365.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-email@outlook.com
EMAIL_SMTP_PASSWORD=your-password
EMAIL_USE_TLS=true
```

#### SendGrid
```bash
EMAIL_SMTP_HOST=smtp.sendgrid.net
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=apikey
EMAIL_SMTP_PASSWORD=your-sendgrid-api-key
EMAIL_USE_TLS=true
```

#### AWS SES
```bash
EMAIL_SMTP_HOST=email-smtp.us-east-1.amazonaws.com  # Use your region
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-ses-smtp-username
EMAIL_SMTP_PASSWORD=your-ses-smtp-password
EMAIL_USE_TLS=true
```

#### Custom SMTP Relay (with authentication)
```bash
EMAIL_SMTP_HOST=your-smtp-server.com
EMAIL_SMTP_PORT=587  # or 465 for TLS, 25 for plain (not recommended)
EMAIL_SMTP_USERNAME=your-username
EMAIL_SMTP_PASSWORD=your-password
EMAIL_USE_TLS=true   # Set to false for plain connection (will try STARTTLS)
```

#### Open SMTP Relay (no authentication)
```bash
EMAIL_SMTP_HOST=your-smtp-server.com
EMAIL_SMTP_PORT=25   # or 587
EMAIL_SMTP_USERNAME=  # Leave empty
EMAIL_SMTP_PASSWORD=  # Leave empty
EMAIL_USE_TLS=false   # Usually false for open relays
```

#### Internal SMTP Relay
```bash
EMAIL_SMTP_HOST=mail.internal.company.com
EMAIL_SMTP_PORT=25
EMAIL_SMTP_USERNAME=  # Leave empty - internal relay doesn't require auth
EMAIL_SMTP_PASSWORD=  # Leave empty
EMAIL_USE_TLS=false
```

## Features

### 1. TLS Support

The service supports two TLS modes:

- **Direct TLS** (port 465): Establishes TLS connection immediately
- **STARTTLS** (port 587): Starts with plain connection, then upgrades to TLS

The service automatically detects and uses STARTTLS if available when `use_tls=false`.

### 2. SMTP Authentication

Uses PLAIN authentication method:
- Authenticates only if username and password are provided
- Works with most SMTP relays

### 3. Email Headers

Properly formatted email headers:
- `From`: Sender name and email
- `To`: Recipient email
- `Subject`: Email subject
- `MIME-Version`: 1.0
- `Content-Type`: text/plain; charset=UTF-8
- `Content-Transfer-Encoding`: 8bit

### 4. Email Validation

Basic email address validation:
- Checks for @ symbol
- Validates local and domain parts are not empty

### 5. Error Handling

Comprehensive error handling:
- Connection errors
- Authentication failures
- SMTP command errors
- Detailed error messages for debugging

## Usage

### Sending Verification Email

The service automatically sends verification emails when:
- User registers via `POST /api/auth/register`
- Admin creates user via `POST /api/users`

### Development Mode

If SMTP is not configured (empty host/port):
- Email content is logged to console
- No actual email is sent
- Useful for local development

## SMTP Relay Flow

1. **Connect to SMTP Server**
   - Direct TLS: `tls.Dial()` for port 465
   - STARTTLS: `net.Dial()` then upgrade for port 587

2. **Authenticate**
   - Uses `smtp.PlainAuth()` if credentials provided
   - Calls `client.Auth()`

3. **Send Email**
   - `client.Mail()` - Set sender
   - `client.Rcpt()` - Set recipient
   - `client.Data()` - Write email content
   - `client.Quit()` - Close connection

## Troubleshooting

### Connection Issues

**Error**: `failed to connect to SMTP server`
- Check `EMAIL_SMTP_HOST` and `EMAIL_SMTP_PORT`
- Verify firewall allows outbound SMTP connections
- Check if server requires VPN

### Authentication Failures

**Error**: `SMTP authentication failed`
- Verify `EMAIL_SMTP_USERNAME` and `EMAIL_SMTP_PASSWORD`
- For Gmail: Use app-specific password, not regular password
- Check if account has 2FA enabled (requires app password)

### TLS Issues

**Error**: `failed to start TLS` or certificate errors
- Try setting `EMAIL_USE_TLS=false` to use STARTTLS
- Verify port matches TLS mode (465 for TLS, 587 for STARTTLS)
- Check if server supports TLS

### Email Not Received

- Check spam/junk folder
- Verify recipient email address is correct
- Check SMTP server logs
- Verify `EMAIL_FROM_EMAIL` is valid and not blocked

## Security Considerations

1. **Credentials**: Store SMTP credentials in environment variables, never in code
2. **TLS**: Always use TLS in production (`EMAIL_USE_TLS=true`)
3. **App Passwords**: Use app-specific passwords for Gmail/Google accounts
4. **Rate Limiting**: Consider implementing rate limiting for email sending
5. **Email Validation**: Validate email addresses before sending

## Testing

### Test Email Configuration

1. Set environment variables:
   ```bash
   export EMAIL_SMTP_HOST=smtp.gmail.com
   export EMAIL_SMTP_PORT=587
   export EMAIL_SMTP_USERNAME=test@gmail.com
   export EMAIL_SMTP_PASSWORD=your-app-password
   export EMAIL_FROM_EMAIL=test@gmail.com
   export EMAIL_FROM_NAME=Test
   export EMAIL_BASE_URL=http://localhost:8080
   ```

2. Register a new user - verification email should be sent

3. Check email inbox (and spam folder)

### Development Testing

Without SMTP configuration, emails are logged to console:
```
[EMAIL] Would send email to user@example.com
[EMAIL] Subject: Verify Your Email Address
[EMAIL] Body: ...
```

## Implementation Details

### Connection Types

**TLS Connection** (`useTLS=true`):
```go
conn, err := tls.Dial("tcp", addr, tlsConfig)
client, err := smtp.NewClient(conn, smtpHost)
```

**STARTTLS Connection** (`useTLS=false`):
```go
conn, err := net.Dial("tcp", addr)
client, err := smtp.NewClient(conn, smtpHost)
if ok, _ := client.Extension("STARTTLS"); ok {
    client.StartTLS(tlsConfig)
}
```

### Email Format

```
From: Lam Phuong <noreply@lamphuong.com>
To: user@example.com
Subject: Verify Your Email Address
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit

Hello,
...
```

## Future Enhancements

1. **HTML Email Support**: Add HTML email templates
2. **Email Queue**: Implement email queue for better reliability
3. **Retry Logic**: Add retry mechanism for failed sends
4. **Email Templates**: Support for multiple email templates
5. **Delivery Tracking**: Track email delivery status
6. **Bounce Handling**: Handle bounced emails

