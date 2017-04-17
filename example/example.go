package main

import (
	"log"

	"github.com/bykovme/tlsmail"
)

const subject = "utf-8 mail: mail subject, тема сообщения, メール件名, 邮件主题"

const body = `
hi, привет, こんにちは, 嗨，

utf-8 body: message, сообщение, メッセージ, 信息

Message sent using GO package github.com/bykovme/tlsmail

Enjoy!
`

func main() {

	mail := tlsmail.TLSMail{
		Host:     "mail.your_favorite_hosting_provider.here",           // smtp & auth host
		Port:     "465",                                                // Default port
		Sender:   "noreply@mail_from.here",                             // sender mail id
		Password: "123456",                                             // sender mail password
		TO:       []string{"mail1@mail_to.here", "mail2@mail_to.here"}, // recipients in TO, can be a list
		CC:       []string{"mail3@mail_cc.here", "mail4@mail_cc.here"}, // recipients in CC, can be a list
		Subject:  subject,                                              // Subject in UTF-8
		Body:     body,                                                 // Mail in UTF-8
	}

	err := mail.Send()

	if err != nil {
		log.Println("Mail send failure: " + err.Error())
	} else {
		log.Println("Mail sent successfully")
	}
}
