# Generate Gmail API Refresh Token

This guide explains how to generate a refresh token for the Gmail API using a simple Go script.

## Prerequisites

- Go installed on your system
- Gmail OAuth credentials (`client_secret.json`) downloaded from Google Cloud Console
- Port 8080 available (or modify the port in the script)

## Steps

### 1. Download OAuth Credentials

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate to **APIs & Services** > **Credentials**
3. Click on your OAuth 2.0 Client ID
4. Download the JSON file and save it as `client_secret.json` in your project root

### 2. Configure Redirect URI

Make sure your OAuth 2.0 Client ID has `http://localhost:8080/oauth2callback` added as an authorized redirect URI:

1. In Google Cloud Console, go to **APIs & Services** > **Credentials**
2. Click on your OAuth 2.0 Client ID
3. Under **Authorized redirect URIs**, add: `http://localhost:8080/oauth2callback`
4. Click **Save**

### 3. Create the Script

Create a file named `generate_token.go` in your project root with the following code:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func main() {
	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatal(err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatal(err)
	}

	// Tạo URL OAuth
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Visit the URL below to authorize:")
	fmt.Println(authURL)

	// Webserver lắng callback
	http.HandleFunc("/oauth2callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		tok, err := config.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal("Token exchange error:", err)
		}

		fmt.Println("Your REFRESH TOKEN:")
		fmt.Println(tok.RefreshToken)

		w.Write([]byte("Done! Copy the refresh token from your console."))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### 4. Install Dependencies

```bash
go get golang.org/x/oauth2/google
go get google.golang.org/api/gmail/v1
```

### 5. Run the Script

```bash
go run generate_token.go
```

### 6. Get Your Refresh Token

1. The script will print an authorization URL in the console
2. Open the URL in your browser
3. Sign in with the Google account you want to use for sending emails
4. Grant the necessary permissions
5. You will be redirected to `http://localhost:8080/oauth2callback`
6. The refresh token will be printed in the console
7. Copy the refresh token and save it securely

### 7. Use the Refresh Token

Add the refresh token to your environment variables:

```env
EMAIL_CLIENT_ID=your-client-id.apps.googleusercontent.com
EMAIL_CLIENT_SECRET=your-client-secret
EMAIL_REFRESH_TOKEN=your-refresh-token-here
EMAIL_FROM_EMAIL=your-email@gmail.com
EMAIL_FROM_NAME=Lam Phuong
```

## Important Notes

- **Security**: Never commit `client_secret.json` or the refresh token to version control
- **Port Conflict**: If port 8080 is already in use, modify the script to use a different port (e.g., `:8081`) and update the redirect URI accordingly
- **Token Validity**: Refresh tokens don't expire unless revoked, so you only need to generate it once
- **Account**: The refresh token is tied to the Google account you authorize with

## Troubleshooting

### "Port already in use"
- Change the port in the script: `http.ListenAndServe(":8081", nil)`
- Update the redirect URI in Google Cloud Console to match the new port

### "Redirect URI mismatch"
- Make sure the redirect URI in Google Cloud Console exactly matches `http://localhost:PORT/oauth2callback`
- The port must match what's in the script

### "Token exchange error"
- Make sure you copied the entire authorization code from the URL
- The code expires quickly, so complete the flow immediately after authorization

## Alternative: Using OAuth 2.0 Playground

You can also use the [OAuth 2.0 Playground](https://developers.google.com/oauthplayground/) to generate a refresh token:

1. Go to https://developers.google.com/oauthplayground/
2. Click the gear icon (⚙️) in the top right
3. Check "Use your own OAuth credentials"
4. Enter your Client ID and Client Secret
5. In the left panel, find "Gmail API v1" and select `https://www.googleapis.com/auth/gmail.send`
6. Click "Authorize APIs"
7. Sign in and grant permissions
8. Click "Exchange authorization code for tokens"
9. Copy the "Refresh token" value

