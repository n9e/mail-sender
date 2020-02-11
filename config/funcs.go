package config

import (
	"fmt"
	"os"
	"time"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/mail"
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
	c := Get()

	mailer := mail.NewSMTP(
		c.Smtp.FromMail,
		c.Smtp.FromName,
		c.Smtp.Username,
		c.Smtp.Password,
		c.Smtp.ServerHost,
		c.Smtp.ServerPort,
		c.Smtp.UseSSL,
		c.Smtp.StartTLS,
	)

	if len(args) == 0 {
		fmt.Println("mail address not given")
		os.Exit(1)
	}

	err := mailer.Send(mail.Mail{
		Tos:     args,
		Subject: "mail from mail-sender",
		Content: fmt.Sprintf("%v", time.Now()),
	})

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
