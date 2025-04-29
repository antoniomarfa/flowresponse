package handles

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"flowresponse/api"
	"flowresponse/models"
)

func HandleToken(w http.ResponseWriter, r *http.Request) {
	// Asegurarse de que el método sea POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Leer el cuerpo de la solicitud
	var tokenRequest models.TokenRequest
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body)) // Reusar el body para r.FormValue

	if err := json.Unmarshal(body, &tokenRequest); err != nil || tokenRequest.Token == "" {
		// Si falla, probamos leerlo desde el formulario (x-www-form-urlencoded)
		r.ParseForm()
		token := r.FormValue("token")
		if token == "" {
			http.Error(w, "Token no encontrado en JSON ni en formulario", http.StatusBadRequest)
			return
		}
		tokenRequest.Token = token
	}

	// Aquí puedes llamar a tu función principal con el token recibido
	// Simulamos que el token es válido y devolvemos una respuesta JSON
	response := api.Flow(tokenRequest.Token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
