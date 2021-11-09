package mail

import (
	"errors"
	"fmt"
	"strings"

	//"fmt"
	"net/smtp"
)
//
type GmailLoginAuth struct {
	username, password string
}

func newLoginAuth(username, password string) smtp.Auth {
	return &GmailLoginAuth{username, password}
}

func (a *GmailLoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *GmailLoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

type MailSender interface {
	SendMail(subject string, msg string, to ...string) error
}

type GmailSender struct {
	userMail     string
	userPassword string
	auth         smtp.Auth
}

func (s *GmailSender) SendMail(subject string, msg string, to ...string) error {
	if to ==nil || len(to) == 0{
		return fmt.Errorf("the to is empty")
	}
		m := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n" , strings.Join(to , ",") , subject,msg)
	err := smtp.SendMail("smtp.gmail.com:587", s.auth, s.userMail, to, []byte(m))
	if err != nil {
		return err
	}
	return nil
}

// NewGmailSender create a GmailSender
func NewMailSender(userMail, username, password string) MailSender{
	return &GmailSender{
		userMail: userMail,
		auth:     newLoginAuth(username, password),
	}
}
