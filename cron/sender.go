package cron

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"path"
	"strings"
	"time"

	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"
	"gopkg.in/gomail.v2"

	"github.com/n9e/mail-sender/config"
	"github.com/n9e/mail-sender/dataobj"
	"github.com/n9e/mail-sender/redisc"
)

var semaphore chan int
var mailer *gomail.Dialer

func SendMails() {
	c := config.Get()

	// 如果发送SMTP的并发太大，怕SMTP服务器受不了
	semaphore = make(chan int, c.Consumer.Worker)

	mailer = gomail.NewDialer(c.Smtp.Host, c.Smtp.Port, c.Smtp.User, c.Smtp.Pass)

	if c.Smtp.InsecureSkipVerify {
		mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	for {
		messages := redisc.Pop(1, c.Consumer.Queue)
		if len(messages) == 0 {
			time.Sleep(time.Duration(300) * time.Millisecond)
			continue
		}

		sendMails(messages)
	}
}

func sendMails(messages []*dataobj.Message) {
	for _, message := range messages {
		semaphore <- 1
		go sendMail(message)
	}
}

func sendMail(message *dataobj.Message) {
	defer func() {
		<-semaphore
	}()

	cnt := len(message.Tos)
	tos := make([]string, 0, cnt)
	for i := 0; i < cnt; i++ {
		item := strings.TrimSpace(message.Tos[i])
		if item == "" {
			continue
		}
		tos = append(tos, item)
	}

	subject := genSubject(message)
	content := genContent(message)

	if len(tos) == 0 {
		logger.Warningf("hashid: %d: subject: %s, tos is empty", message.Event.HashId, subject)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Get().Smtp.User)
	m.SetHeader("To", message.Tos...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	err := mailer.DialAndSend(m)

	logger.Infof("hashid: %d: subject: %s, tos: %v, error: %v", message.Event.HashId, subject, message.Tos, err)
	logger.Infof("hashid: %d: endpoint: %s, metric: %s, tags: %s", message.Event.HashId, message.ReadableEndpoint, strings.Join(message.Metrics, ","), message.ReadableTags)
}

var ET = map[string]string{
	"alert":    "告警",
	"recovery": "恢复",
}

func genSubject(message *dataobj.Message) string {
	subject := ""
	if message.IsUpgrade {
		subject = "[报警已升级]" + subject
	}

	return fmt.Sprintf("[P%d %s]%s - %s", message.Event.Priority, ET[message.Event.EventType], message.Event.Sname, message.ReadableEndpoint)
}

func parseEtime(etime int64) string {
	t := time.Unix(etime, 0)
	return t.Format("2006-01-02 15:04:05")
}

func genContent(message *dataobj.Message) string {
	fp := path.Join(runner.Cwd, "etc", "mail.html")
	t, err := template.ParseFiles(fp)
	if err != nil {
		payload := fmt.Sprintf("InternalServerError: cannot parse %s %v", fp, err)
		logger.Errorf(payload)
		return fmt.Sprintf(payload)
	}

	var body bytes.Buffer
	err = t.Execute(&body, map[string]interface{}{
		"IsAlert":   message.Event.EventType == "alert",
		"Status":    ET[message.Event.EventType],
		"Sname":     message.Event.Sname,
		"Endpoint":  message.ReadableEndpoint,
		"Metric":    strings.Join(message.Metrics, ","),
		"Tags":      message.ReadableTags,
		"Value":     message.Event.Value,
		"Info":      message.Event.Info,
		"Etime":     parseEtime(message.Event.Etime),
		"Elink":     message.EventLink,
		"Slink":     message.StraLink,
		"Clink":     message.ClaimLink,
		"IsUpgrade": message.IsUpgrade,
		"Bindings":  message.Bindings,
	})

	if err != nil {
		logger.Errorf("InternalServerError: %v", err)
		return fmt.Sprintf("InternalServerError: %v", err)
	}

	return body.String()
}
