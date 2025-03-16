package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"flowresponse/database"
	"flowresponse/handles"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar el archivo .env

	dsn := os.Getenv("DSN")
	ApiUrl := os.Getenv("APIURL")

	fmt.Println("dsn entrada ", dsn)
	if ApiUrl == "" {
		err := godotenv.Load(".env")
		fmt.Println("cargo el archivo env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		dsn = os.Getenv("DSN")
	}

	// Acceder a las variables de entorno
	fmt.Println("nuevo dsn des env", dsn)
	db, err := database.GetDB(dsn)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Obtener la instancia de *sql.DB para cerrar la conexión
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener la instancia de *sql.DB: %v", err)
	}
	defer sqlDB.Close() // Cerrar la conexión al finalizar

	// Iniciar el servidor HTTP en una goroutine
	http.HandleFunc("/token", handles.HandleToken)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
