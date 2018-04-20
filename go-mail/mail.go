package mail

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"strings"
	"sync"
)

type Sender interface {
	Login(user, pw, host, serverAddr string) Sender
	Send(subject, message string, to ...string) Sender
	SetNickname(string) Sender
	Done() []error
}

func NewSender() Sender {
	return new(emailAccount)
}

type receiver struct {
	subject  string
	msg      string
	accounts []string
}

type emailAccount struct {
	smtp.Auth
	user       string
	nickname   string
	serverAddr string
	receivers  []*receiver
}

func newReceiver(subject, msg string, accounts ...string) *receiver {
	r := new(receiver)
	r.subject = subject
	r.msg = msg
	r.accounts = accounts
	return r
}

func (e *emailAccount) Login(user, pw, host, serverAddr string) Sender {
	e.Auth = smtp.PlainAuth("", user, pw, host)
	e.user = user
	e.serverAddr = serverAddr
	return e
}

func (e *emailAccount) SetNickname(nickname string) Sender {
	e.nickname = nickname
	return e
}

// Send
// 每一次Send代表发送一次邮件
func (e *emailAccount) Send(subject, message string, accounts ...string) Sender {
	e.receivers = append(e.receivers, newReceiver(subject, message, accounts...))
	return e
}

func send(auth smtp.Auth, serverAddr, from, subject, msg, nickname string, to ...string) error {
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s<%s>", nickname, from)
	header["To"] = strings.Join(to, ";")
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(subject)))
	header["Content-Type"] = "text/plain; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"
	var body string
	for k, v := range header {
		body += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	body += "\r\n" + base64.StdEncoding.EncodeToString([]byte(msg))
	return smtp.SendMail(serverAddr, auth, from, to, []byte(body))
}

func (e *emailAccount) Done() (errs []error) {
	var group sync.WaitGroup
	if e.nickname == "" {
		e.nickname = e.user
	}

	for _, r := range e.receivers {
		group.Add(1)
		go func(r *receiver) {
			if err := send(e.Auth, e.serverAddr, e.user, r.subject, r.msg, e.nickname, r.accounts...); err != nil {
				errs = append(errs, err)
			}
			group.Done()
		}(r)
	}

	group.Wait()
	return
}
