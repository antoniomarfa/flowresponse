package handles

import (
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Decodificar el cuerpo JSON
	var tokenRequest models.TokenRequest
	if err := json.Unmarshal(body, &tokenRequest); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Validar que el token no esté vacío
	if tokenRequest.Token == "" {
		http.Error(w, "El token es requerido", http.StatusBadRequest)
		return
	}

	// Aquí puedes llamar a tu función principal con el token recibido
	// Simulamos que el token es válido y devolvemos una respuesta JSON
	response := api.Flow(tokenRequest.Token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
