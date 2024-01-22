package mail

import (
	"log"
	"sync"
	"github.com/go-gomail/gomail"
)

func SendEmail(sender, subject string, emailList []string, AttachFile map[string][]string) error {
	// Buffer the result channel
	resultChannel := make(chan EmailResult, len(emailList))

	// Use sync.WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Iterate over each recipient email
	for _, recipientEmail := range emailList {
		// Increment the WaitGroup counter
		wg.Add(1)

		// Launch a goroutine to send email
		go sendEmailAsync(sender, subject, recipientEmail, AttachFile[recipientEmail], &wg, resultChannel)
	}

	// Close the result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect the results from the channel
	var sendError error
	for result := range resultChannel {
		if result.Error != nil {
			log.Printf("Error sending email to %s: %v", result.RecipientEmail, result.Error)
			sendError = result.Error
		}
	}

	return sendError
}

func sendEmailAsync(sender, subject, email string, files []string, wg *sync.WaitGroup, resultChannel chan<- EmailResult) {
	defer wg.Done()

	// Create a new message for each goroutine
	message := gomail.NewMessage()
	message.SetAddressHeader("From", smtpConfig.User, sender)
	message.SetHeader("To", email)
	message.SetHeader("Reply-To", sender)
	message.SetHeader("Subject", subject)

	body := `
Dear Partner,
The attached file is your secured information file.
If you have any questions or concerns regarding the files attached or any other matter, please feel free to reach out to me.
Best regards,
DEV-TRACTHAI
	`
	message.SetBody("text/plain", body)

	// Attach files for the current recipient
	for _, filePath := range files {
		message.Attach(filePath)
	}

	// Send the email for the current recipient
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
	err := d.DialAndSend(message)

	// Send the result through the channel
	resultChannel <- EmailResult{
		RecipientEmail: email,
		Error:          err,
	}
}

type EmailResult struct {
	RecipientEmail string
	Error          error
}

// SMTPConfig contains SMTP configuration values
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

var smtpConfig = SMTPConfig{
	Host:     "smtp.gmail.com",
	Port:     587,
	User:     "report.trac@gmail.com",
	Password: "mcoqvwpabjtdoxvw",
}







// package mail

// import (
// 	"log"
// 	"github.com/go-gomail/gomail"
// )

// func SendEmail(sender, subject string, emailList []string, AttachFile map[string][]string) error {
// 	for _, recipientEmail := range emailList {
// 		err := sendEmail(sender, subject, recipientEmail, AttachFile[recipientEmail])
// 		if err != nil {
// 			log.Printf("Error sending email to %s: %v", recipientEmail, err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// func sendEmail(sender, subject, email string, files []string) error {
// 	message := gomail.NewMessage()
// 	message.SetAddressHeader("From", smtpConfig.User, sender)
// 	message.SetHeader("To", email)
// 	message.SetHeader("Reply-To", sender)
// 	message.SetHeader("Subject", subject)

// 	body := `
// Dear Partner,
// The attached file is your secured information file.
// If you have any questions or concerns regarding the files attached or any other matter, please feel free to reach out to me.
// Best regards,
// DEV-TRACTHAI
// 	`
// 	message.SetBody("text/plain", body)

// 	// Attach files for the current recipient
// 	for _, filePath := range files {
// 		message.Attach(filePath)
// 	}

// 	// Send the email for the current recipient
// 	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Password)
// 	return d.DialAndSend(message)
// }


// type EmailResult struct {
// 	RecipientEmail string
// 	Error          error
// }

// // SMTPConfig contains SMTP configuration values
// type SMTPConfig struct {
// 	Host     string
// 	Port     int
// 	User     string
// 	Password string
// }

// var smtpConfig = SMTPConfig{
// 	Host:     "smtp.gmail.com",
// 	Port:     587,
// 	User:     "report.trac@gmail.com",
// 	Password: "mcoqvwpabjtdoxvw",
// }