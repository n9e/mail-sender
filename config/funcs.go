package config

import (
	"crypto/tls"
	"fmt"
	"os"

	"github.com/toolkits/pkg/logger"
	"gopkg.in/gomail.v2"
)

// InitLogger init logger toolkits
func InitLogger() {
	c := Get().Logger

	lb, err := logger.NewFileBackend(c.Dir)
	if err != nil {
		fmt.Println("cannot init logger:", err)
		os.Exit(1)
	}

	lb.SetRotateByHour(true)
	lb.SetKeepHours(c.KeepHours)

	logger.SetLogging(c.Level, lb)
}

func TestSMTP(args []string) {
	if len(args) == 0 {
		fmt.Println("mail address not given")
		os.Exit(1)
	}

	c := Get()

	d := gomail.NewDialer(c.Smtp.Host, c.Smtp.Port, c.Smtp.User, c.Smtp.Pass)
	if c.Smtp.InsecureSkipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", c.Smtp.User)
	m.SetHeader("To", args...)
	m.SetHeader("Subject", "Hello! 中文标题 N9E test")
	m.SetBody("text/html", "Hello <b>N9E User</b> and <i>中文内容</i>! ")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
