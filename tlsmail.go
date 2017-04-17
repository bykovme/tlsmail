package tlsmail

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
)

// TLSMail - mail structure
type TLSMail struct {
	Host     string
	Port     string
	Sender   string
	Password string
	TO       []string
	CC       []string
	Subject  string
	Body     string
}

// ServerName - build server name
func (tlsmail *TLSMail) ServerName() string {
	return tlsmail.Host + ":" + tlsmail.Port
}

// EncodedSubject - encode subject according to RFC2047, so UTF can be used
func (tlsmail *TLSMail) EncodedSubject() string {
	addr := mail.Address{
		Address: "",
		Name:    tlsmail.Subject,
	}
	return strings.Trim(addr.String(), " <>@")
}

// BuildMessage - build mail message
func (tlsmail *TLSMail) BuildMessage() string {

	header := make(map[string]string)
	header["From"] = tlsmail.Sender
	header["To"] = strings.Join(tlsmail.TO, ";")
	if len(tlsmail.CC) > 0 {
		header["CC"] = strings.Join(tlsmail.CC, ";")
	}
	header["Subject"] = tlsmail.EncodedSubject()
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(tlsmail.Body))

	return message
}

// CheckMandatoryFields - check if all mandatory fields are set
func (tlsmail *TLSMail) CheckMandatoryFields() error {
	if tlsmail.Sender == "" {
		return errors.New("Check error: Sender cannot be empty")
	}
	if tlsmail.Password == "" {
		return errors.New("Check error: Password cannot be empty")
	}
	if tlsmail.Host == "" {
		return errors.New("Check error: Host cannot be empty")
	}
	if tlsmail.Port == "" {
		return errors.New("Check error: Port cannot be empty")
	}
	if len(tlsmail.TO) == 0 {
		return errors.New("Check error: there should be at least one recipient in TO")
	}
	for _, mail := range tlsmail.TO {
		if mail == "" {
			return errors.New("Check error: inappropriate mail in TO")
		}
	}
	if tlsmail.Subject == "" {
		return errors.New("Check error: Subject cannot be empty")
	}
	if tlsmail.Body == "" {
		return errors.New("Check error: Body cannot be empty")
	}
	return nil
}

// Send - send mail
func (tlsmail *TLSMail) Send() error {
	err := tlsmail.CheckMandatoryFields()
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", tlsmail.Sender, tlsmail.Password, tlsmail.Host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         tlsmail.Host,
	}

	conn, err := tls.Dial("tcp", tlsmail.ServerName(), tlsconfig)
	if err != nil {
		return err
	}

	smtpClient, err := smtp.NewClient(conn, tlsmail.Host)
	if err != nil {
		return err
	}
	defer smtpClient.Quit()

	if err = smtpClient.Auth(auth); err != nil {
		return err
	}

	if err = smtpClient.Mail(tlsmail.Sender); err != nil {
		return err
	}

	receivers := append(tlsmail.TO, tlsmail.CC...)
	for _, receiverMail := range receivers {
		if err = smtpClient.Rcpt(receiverMail); err != nil {
			return err
		}
	}

	dataWriter, err := smtpClient.Data()
	if err != nil {
		return err
	}

	_, err = dataWriter.Write([]byte(tlsmail.BuildMessage()))
	if err != nil {
		return err
	}

	err = dataWriter.Close()
	if err != nil {
		return err
	}

	return nil
}
