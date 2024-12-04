package handlers

import "log"

func sendEmailToUser(email string, content string) string {
	return "success"
}

func sendWarningEmail(email string) {
	content := "Someone just used your session with a new IP address. If this wasn't you, change your password"
	if sendEmailToUser(email, content) == "success" {
		log.Printf("sent a warning email to %s", email)
	}
}
