package services

import (
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	Sender   string
	Password string
	Host     string
	Port     int
	Receiver string
}

func NewEmailService() *EmailService {
	return &EmailService{
		Sender:   os.Getenv("EMAIL_SENDER"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Host:     "smtp.office365.com",
		Port:     587,
		Receiver: os.Getenv("EMAIL_RECEIVER"),
	}
}

func (s *EmailService) SendEmail(name, email, phone string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Sender)
	m.SetHeader("To", s.Receiver)
	m.SetHeader("Subject", "Nuevo Registro en el Bot")
	m.SetBody("text/html",
		"<h3>Nuevo usuario registrado</h3>"+
			"<p><strong>Nombre:</strong> "+name+"</p>"+
			"<p><strong>Correo:</strong> "+email+"</p>"+
			"<p><strong>Teléfono:</strong> "+phone+"</p>",
	)

	d := gomail.NewDialer(s.Host, s.Port, s.Sender, s.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Println("Error enviando correo:", err)
		return err
	}

	log.Println("Correo enviado correctamente")
	return nil
}
