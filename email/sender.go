package email

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendWelcomeEmail(to string, userID string) error {
	// Se não tiver configuração SMTP, apenas simula o envio
	from := os.Getenv("SMTP_USER")
	if from == "" {
		fmt.Printf("📧 [email] Email de boas-vindas simulado para %s (ID: %s)\n", to, userID)
		return nil
	}

	password := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

	if password == "" || host == "" || port == "" {
		fmt.Printf("📧 [email] Email de boas-vindas simulado para %s (ID: %s) - Configuração SMTP incompleta\n", to, userID)
		return nil
	}

	addr := host + ":" + port
	auth := smtp.PlainAuth("", from, password, host)

	subject := "Subject: Bem-vindo!\n"
	body := fmt.Sprintf("Olá,\n\nSeu cadastro foi realizado com sucesso!\nID do usuário: %s\n\nAbraços,\nEquipe\n", userID)
	msg := []byte(subject + "\n" + body)

	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return err
	}

	fmt.Printf("📧 [email] Email enviado com sucesso para %s\n", to)
	return nil
}
