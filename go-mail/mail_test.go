package mail

import "testing"
import "os"
import "io/ioutil"
import "encoding/json"

type testMailJSON struct {
	SenderAccount  string   `json:"senderAccount"`
	SenderIdentity string   `json:"senderIdentity"`
	SenderPassword string   `json:"senderPassword"`
	Host           string   `json:"host"`
	ServerAddr     string   `json:"serverAddr"`
	Msg            string   `json:"msg"`
	Subject        string   `json:"subject"`
	To             []string `json:"to"`
}

var config testMailJSON

func readConfig(t *testing.T) {
	f, err := os.Open("./mail_test.json")
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

}
func TestMail(t *testing.T) {
	readConfig(t)
	groupSender := NewSender()

	err := groupSender.
		Login(config.SenderAccount, config.SenderPassword, config.Host, config.ServerAddr).
		SetNickname(config.SenderIdentity).
		Send(config.Subject, config.Msg, config.To...).
		Done()

	if err != nil {
		t.Errorf("group send failed. Errors: %v", err)
	}

	// singleSender := NewSender()
	// err = singleSender.
	// 	Login(config.SenderAccount, config.SenderPassword, config.Host, config.ServerAddr).
	// 	Send(config.Subject, config.Msg, config.To[0]).
	// 	Send(config.Subject, config.Msg, config.To[1]).
	// 	Done()
	// if err != nil {
	// 	t.Errorf("group send failed. Errors: %v", err)
	// }

}
