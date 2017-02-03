package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"time"
)

type Notifier struct {
	smtpServer      string
	smtpPort        string
	username        string
	password        string
	from            string
	to              string
	retries         int
	lastAlertFile   string
	minuteThreshold int
	timeFormat      string
}

// cameraCheck returns an error if there is an issue with the cameras,
// nil otherwise
func cameraCheck() error {
	return fmt.Errorf("could not find camera files")
	// log.Println("no camera errors found")
	// return nil
}

// email sends an email notification
// it does a max of n.retries sending the alert
func (n Notifier) email(message string) {
	n._email(message, 0)
}

func (n Notifier) _email(message string, attempt int) {
	if attempt > n.retries {
		log.Println("max attempts reached sending message: %s", message)
	}

	auth := smtp.PlainAuth("",
		n.username,
		n.password,
		n.smtpServer,
	)

	subject := "Camera Alert!"
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s",
		n.from,
		n.to,
		subject,
		message)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", n.smtpServer, n.smtpPort),
		auth,
		n.username,
		[]string{n.to},
		[]byte(msg))

	if err != nil {
		log.Println("Error attempting to send a mail, will retry:", err)
		time.Sleep(time.Second * time.Duration(5^attempt))
		n._email(message, attempt+1)
		return
	}

	err = ioutil.WriteFile(
		n.lastAlertFile,
		[]byte(time.Now().UTC().Format(n.timeFormat)),
		0644)
	if err != nil {
		log.Println("error writing to last-alert-file: %s", err)
	}
	log.Print("Message send successful")
}

// lastAlert gets the time we last alerted
// if this is the first run returns 0
func (n Notifier) lastAlertTime() (time.Time, error) {
	content, err := ioutil.ReadFile(n.lastAlertFile)
	if os.IsNotExist(err) {
		fmt.Println("here 1")
		return time.Time{}, nil
	}
	if err != nil {
		fmt.Println("here 2")
		log.Println("Error getting time of last alert", err)
	}
	fmt.Println("here 3")
	return time.Parse(n.timeFormat, string(content))
}

func (n *Notifier) parseConfig() {
	var smtpServer = flag.String("smtp-server", "localhost", "the email server to send alert")
	var smtpPort = flag.String("smtp-port", "587", "the email server port to use")
	var username = flag.String("username", "", "user to send with")
	var password = flag.String("password", "", "password for user")
	var from = flag.String("from", "", "who the notification is from")
	var to = flag.String("to", "", "who to send the notification to")
	var retries = flag.Int("retries", 6, "max number of times to retry sending notification")
	var lastAlertFile = flag.String("last-alert-file", "last_alert", "file holds timestamp of when the last alert was sent. This is used to prevent an email flood.")
	var minuteThreshold = flag.Int("minute-threshold", 10, "Minutes to wait before sending another alert. This is used to prevent an email flood.")
	flag.Parse()

	n.smtpServer = *smtpServer
	n.smtpPort = *smtpPort
	n.password = *password
	n.username = *username
	n.from = *from
	n.to = *to
	n.retries = *retries
	n.lastAlertFile = *lastAlertFile
	n.minuteThreshold = *minuteThreshold
	n.timeFormat = "Jan 2, 2006 15:04:05"
}

func main() {
	notifier := Notifier{}
	notifier.parseConfig()

	lastAlert, err := notifier.lastAlertTime()
	minutesSinceLastAlert := time.Since(lastAlert).Minutes()
	if err != nil {
		log.Println("error getting time of last alert")
	} else if minutesSinceLastAlert < 10 {
		log.Printf("Already sent an alert within last 10m (%.2fm ago), not sending another", minutesSinceLastAlert)
		return
	}

	err = cameraCheck()
	if err != nil {
		notifier.email(err.Error())
	}
}
