package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"langtools/services"

)

type EmailRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("Recibiendo petición para enviar correo")

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al leer los datos", http.StatusBadRequest)
		return
	}

	emailService := services.NewEmailService()
	if err := emailService.SendEmail(req.Name, req.Email, req.Phone); err != nil {
		http.Error(w, "Error enviando el correo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Correo enviado con éxito"))
}
