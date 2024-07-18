package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func createMail() Mail {
	domain := os.Getenv("SMTP_DOMAIN")
	if domain == "" {
		log.Fatalf("SMTP_DOMAIN is not set")
	}

	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Fatalf("SMTP_HOST is not set")
	}
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatalf("Error converting SMTP_PORT to int: %v", err)
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	encryption := os.Getenv("SMTP_ENCRYPTION")
	fromAddress := os.Getenv("SMTP_FROM_ADDRESS")
	if fromAddress == "" {
		log.Fatalf("SMTP_FROM_ADDRESS is not set")
	}
	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		log.Fatalf("SMTP_FROM_NAME is not set")
	}

	return Mail{
		Domain: domain,
		Host: host,
		Port: port,
		Username: username,
		Password: password,
		Encryption: encryption,
		FromAddress: fromAddress,
		FromName: fromName,
	}
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)

	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()

	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextHTML, formattedMessage)
	email.AddAlternative(mail.TextPlain, plainMessage)

	if len(msg.Attachments) > 0 {
		for _, v := range msg.Attachments {
			email.AddAttachment(v)
		}
	}

	err = email.Send(smtpClient)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) getEncryption(encryption string) mail.Encryption {
	switch encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.html.gohtml"

	tmpl, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	
	if err = tmpl.ExecuteTemplate(&doc, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := doc.String()

	return m.inlineCSS(formattedMessage)
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := "./templates/mail.plain.gohtml"

	tmpl, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var doc bytes.Buffer
	
	if err = tmpl.ExecuteTemplate(&doc, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := doc.String()

	return plainMessage, nil
}

func (m *Mail) inlineCSS(html string) (string, error) {
	options := premailer.Options{
		RemoveClasses:true,
		CssToAttributes:true,
		KeepBangImportant:true,
	}

	prem, err := premailer.NewPremailerFromString(html, &options)
	if err != nil {
		return "", err
	}

	return prem.Transform()
}
