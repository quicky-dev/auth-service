package auth

import (
	"fmt"
	"os"
)

func formatVerificationURL(verificationCode string, userID string) string {
	if os.Getenv("PRODUCTION") == "" {
		return fmt.Sprintf("http://localhost:3000/verify/email?v=%s&u=%s", verificationCode, userID)
	} else {
		return fmt.Sprintf("https://auth.quicky.dev/verify/email?v=%s&u=%s", verificationCode, userID)
	}
}
