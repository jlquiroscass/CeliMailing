package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gorilla/mux"
)

const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

// Mail The mail Type (more like an object)
type Mail struct {
	Subject 	string `json:"subject,omitempty"`
	To      	string `json:"to,omitempty"`
	Body    	string `json:"body,omitempty"`
	Sender  	string `json:"sender,omitempty"`
	TextBody  	string `json:"txtBody,omitempty"`
}

// send mail
func sendMail(w http.ResponseWriter, r *http.Request) {
	var mail Mail
	_ = json.NewDecoder(r.Body).Decode(&mail)
	log.Print("Sending mail to " + mail.To)
	var Recipient = mail.To
	var Subject = mail.Subject
	var HtmlBody = mail.Body
	var Sender = mail.Sender
	var TextBody = mail.TextBody

	// Create a new session in the eu-west-1 region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Print(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Print(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Print(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Print(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Print(err.Error())
		}

		return
	}

	log.Print("Email Sent to address: " + Recipient)
	log.Print(result)
	json.NewEncoder(w).Encode(mail)
}

// main function to boot up everything
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mail", sendMail).Methods("POST")
	log.Fatal(http.ListenAndServe(":6565", router))
}
