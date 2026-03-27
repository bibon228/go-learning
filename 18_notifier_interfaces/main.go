package main

import "fmt"

type Notifier interface {
	Send(message string) string
}

type EmailNotifier struct {
	Email string
}

func (e *EmailNotifier) Send(message string) string {
	return fmt.Sprintf("Отправка письма на %s: %s", e.Email, message)
}

type SMSNotifier struct {
	Phone string
}

func (s *SMSNotifier) Send(message string) string {
	return fmt.Sprintf("Отправка СМС на %s: %s", s.Phone, message)
}

type TelegramNotifier struct {
	UserName string
}

func (t *TelegramNotifier) Send(message string) string {
	return fmt.Sprintf("Отправка СМС на %s: %s", t.UserName, message)
}
func Alert(n Notifier, f Formatter, message string) {
	fancyText := f.text(message)
	result := n.Send(fancyText)
	fmt.Println(result)
}

type Formatter interface {
	text(message string) string
}
type HTMLFormatter struct{}

func (h *HTMLFormatter) text(message string) string {
	return fmt.Sprintf("<b>%s</b>", message)
}

type JSONFormatter struct{}

func (j *JSONFormatter) text(message string) string {
	return fmt.Sprintf("{\"message\": \"%s\"}", message)
}

func main() {
	// Твой код Notification System
	email := &EmailNotifier{
		Email: "[EMAIL_ADDRESS]",
	}
	sms := &SMSNotifier{
		Phone: "1234567890",
	}
	telegram := &TelegramNotifier{
		UserName: "@username",
	}
	Alert(email, &HTMLFormatter{}, "Hello from Go!")
	Alert(sms, &JSONFormatter{}, "Hello from Go!")
	Alert(telegram, &HTMLFormatter{}, "Hello from Go!")
}
