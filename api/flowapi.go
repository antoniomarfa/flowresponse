package api

import (
	"flowresponse/database"
	"flowresponse/models"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/antoniomarfa/traveltools/utils"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Flow(token string) models.TokenResponse {

	//variables de respuesta
	status_servicio := "Ok"
	error_servicio := "terminado correctamente"

	dsn := os.Getenv("DSN")
	ApiUrl := os.Getenv("APIURL")

	fmt.Println("dsn entrada api ", dsn)
	if ApiUrl == "" {
		err := godotenv.Load(".env")
		fmt.Println("cargo el archivo env api")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		dsn = os.Getenv("DSN")
		ApiUrl = os.Getenv("APIURL")
	}

	// Acceder a las variables de entorno

	db, err := database.GetDB(dsn)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
		status_servicio = "No"
		error_servicio = "Error al conectar a la base de datos: "
	}

	// Obtener la instancia de *sql.DB para cerrar la conexión
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error al obtener la instancia de *sql.DB: %v", err)
		status_servicio = "No"
		error_servicio = "Error al obtener la instancia de *sql.DB "
	}
	defer sqlDB.Close() // Cerrar la conexión al finalizar
	//----------------------
	//Buscar el ingreso por token flow
	var ingreso models.IngresoReq
	maxRetries := 7
	found := false

	for i := 0; i < maxRetries; i++ {
		err := db.Where("token_flow = ?", token).First(&ingreso).Error
		if err == nil {
			found = true
			break
		}
		if err != gorm.ErrRecordNotFound {
			fmt.Println("Error inesperado al buscar el ingreso:", err)
			status_servicio = "No"
			error_servicio = "Error inesperado al buscar el ingreso"
			return models.TokenResponse{
				Status:  status_servicio,
				Message: error_servicio,
			}
		}
		fmt.Printf("Intento %d/%d: ingreso no encontrado, esperando...\n", i+1, maxRetries)
		time.Sleep(500 * time.Millisecond)
	}

	if !found {
		fmt.Println("No se encontró ingreso con el token tras varios intentos")
		status_servicio = "No"
		error_servicio = "Ingreso no encontrado tras varios intentos"
		return models.TokenResponse{
			Status:  status_servicio,
			Message: error_servicio,
		}
	}

	//---------------------
	//con el company id del ingeso buscar las key de flow
	var flowcon models.CompanyconResp
	if err := db.Where("company_id = ?", ingreso.CompanyId).First(&flowcon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("No se encontró ningún registro con el company_id proporcionado")
			status_servicio = "No"
			error_servicio = "No se encontró ningún registro con el company_id proporcionado "
		} else {
			fmt.Println("Error al buscar el ingreso:", err)
			status_servicio = "No"
			error_servicio = "Error al buscar el ingreso "
		}
	}

	// Crear una instancia de FlowApi
	api := &utils.FlowApi{} // Usamos & para obtener un puntero a FlowApi

	//Usar la instancia para pasar las keys
	// Configurar FlowApi
	api.SetApiKey(flowcon.ApikeyFlow)
	api.SetSecretKey(flowcon.SecretkeyFlow)
	api.SetApiURL(ApiUrl)

	// Parámetros para el método Send
	// Llamar a la API de Flow
	service := "payment/getStatus"
	method := "GET"
	params := map[string]string{"token": token}

	// Llamar al método Send
	response, err := api.Send(service, params, method)
	if err != nil {
		fmt.Println("Error:", err)

	}
	fmt.Println("respuesta ", response)
	// Extraer los valores del response y asignarlos a la estructura PaymentResponse
	continuaOperacion := true
	paymentResponse, err := parseResponse(response)
	if err != nil {
		fmt.Println("Error parsing response", err)
		status_servicio = "No"
		error_servicio = "Error al convertir el body a JSON "
		continuaOperacion = false
	}

	// Procesar la respuesta
	if continuaOperacion {
		status := paymentResponse.Status
		switch status {
		case "1":
			fmt.Println("Pendiente de Pago")
		case "2":
			fmt.Println("Pagado")
			// Convertir la cadena a float64
			montoIngreso, err := strconv.ParseFloat(paymentResponse.Amount, 32)
			if err != nil {
				fmt.Println("Error al convertir:", err)

			}
			transaccion := paymentResponse.FlowOrder
			nroIngreso := paymentResponse.CommerceOrder
			fechaIngreso := paymentResponse.RequestDate[:10] // Formato YYYY-MM-DD
			fechaTrans := paymentResponse.RequestDate[:10]
			fechaAuto := paymentResponse.RequestDate
			tipoPago := "FW"
			media := paymentResponse.PaymentData.Media
			rutAl := paymentResponse.Optional.Alumno

			fecha_Ingreso, err := time.Parse("2006-01-02", fechaIngreso)
			if err != nil {
				fmt.Println("Error al parsear la fecha:", err)
				status_servicio = "No"
				error_servicio = "Error al parsear la fecha fechaingres"

			}
			fecha_Trans, err := time.Parse("2006-01-02", fechaTrans)
			if err != nil {
				fmt.Println("Error al parsear la fecha:", err)
				status_servicio = "No"
				error_servicio = "Error al parsear la fecha fechatrans"

			}
			fecha_Auto, err := time.Parse("2006-01-02 15:04:05", fechaAuto)
			if err != nil {
				fmt.Println("Error al parsear la fecha:", err)
				status_servicio = "No"
				error_servicio = "Error al parsear la fecha fechaauto"

			}

			if ingreso.Nrocuotas > 0 {

				fecha_inicial := ingreso.Fechainicial.Format("2006-01-02")
				valorCuota := float64(ingreso.Valorcuota)
				// Procesar cuotas
				//	dia := fecha_inicial[8:10]
				mes := fecha_inicial[5:7]
				agno := fecha_inicial[:4]

				for i := 0; i < ingreso.Nrocuotas; i++ {
					cuota := i + 1
					pago := models.CreatePagosReq{
						Tipocom:       "COW",
						IngresoId:     ingreso.ID,
						Identificador: nroIngreso,
						Fecha:         fecha_Ingreso,
						SaleId:        ingreso.SaleId,
						Rutalumn:      rutAl,
						Transaccion:   transaccion,
						Tipo:          tipoPago,
						Monto:         valorCuota,
						Nrotarjeta:    "",
						Codigoauto:    "",
						Fechaauto:     fecha_Auto,
						Tipopago:      media,
						Nrocuota:      0,
						Fechatransac:  fecha_Trans,
						Activo:        1,
						Author:        "",
						Cuotapagada:   cuota,
						Cuotafecha:    agno + mes,
					}

					if err := db.Create(&pago).Error; err != nil {
						fmt.Println("Error al crear el pago: ", err)
						status_servicio = "No"
						error_servicio = "Error al crear el pago"

					}

					// Incrementar el mes
					mesInt, _ := strconv.Atoi(mes)
					mesInt++
					if mesInt > 12 {
						mesInt = 1
						agnoInt, _ := strconv.Atoi(agno)
						agnoInt++
						agno = strconv.Itoa(agnoInt)
					}
					mes = fmt.Sprintf("%02d", mesInt)
				}
			} else {
				// Insertar un solo pago
				pago := models.CreatePagosReq{
					Tipocom:       "COW",
					IngresoId:     ingreso.ID,
					Identificador: nroIngreso,
					Fecha:         fecha_Ingreso,
					SaleId:        ingreso.SaleId,
					Rutalumn:      rutAl,
					Transaccion:   transaccion,
					Tipo:          tipoPago,
					Monto:         montoIngreso,
					Nrotarjeta:    "",
					Codigoauto:    "",
					Fechaauto:     fecha_Auto,
					Tipopago:      media,
					Nrocuota:      0,
					Fechatransac:  fecha_Trans,
					Activo:        1,
					Author:        "",
					Cuotapagada:   0,
					Cuotafecha:    "",
				}

				if err := db.Create(&pago).Error; err != nil {
					fmt.Println("Error al crear el pago: ", err)
					status_servicio = "No"
					error_servicio = "Error al crear el pago"
				}
			}
		case "3":
			fmt.Println("Transacción Rechazada")
		case "4":
			fmt.Println("Transacción Anulada")
		default:
			fmt.Println("Estado desconocido")
		}

		// Actualizar el estado del ingreso
		statusText := ""
		switch status {
		case "1":
			statusText = "Pendiente de Pago"
		case "2":
			statusText = "Pagado"
		case "3":
			statusText = "Transacción Rechazada"
		case "4":
			statusText = "Transacción Anulada"
		default:
			statusText = "Estado desconocido"
		}

		if err := db.Model(&ingreso).Update("status_pago", statusText).Error; err != nil {
			fmt.Println("Error al actualizar el estado del pago:", err)
			status_servicio = "No"
			error_servicio = "Error al actualizar el estado del pago en ingreso"
		}
	}
	// Si todo está bien, devuelve un TokenResponse y nil como error
	return models.TokenResponse{
		Status:  status_servicio,
		Message: error_servicio,
	}
}

// Función para parsear el response a la estructura PaymentResponse
func parseResponse(response map[string]interface{}) (models.PaymentResponse, error) {
	var paymentResponse models.PaymentResponse

	// Asignar valores directamente desde el map
	paymentResponse.Amount = response["amount"].(string)
	paymentResponse.CommerceOrder = response["commerceOrder"].(string)
	paymentResponse.Currency = response["currency"].(string)
	paymentResponse.FlowOrder = fmt.Sprint(response["flowOrder"])   // Convertir a string
	paymentResponse.Merchantid = fmt.Sprint(response["merchantId"]) // Si es nil, se asigna un valor vacío

	// Optional
	optional := response["optional"].(map[string]interface{})
	paymentResponse.Optional.Venta = fmt.Sprint(optional["venta"])
	paymentResponse.Optional.Alumno = optional["alumno"].(string)

	// Payer
	paymentResponse.Payer = response["payer"].(string)

	// PaymentData
	paymentData := response["paymentData"].(map[string]interface{})
	paymentResponse.PaymentData.Amount = fmt.Sprint(paymentData["amount"])
	paymentResponse.PaymentData.Balance = fmt.Sprint(paymentData["balance"])
	paymentResponse.PaymentData.Conversiondate = fmt.Sprint(paymentData["conversionDate"])
	paymentResponse.PaymentData.Conversionrate = fmt.Sprint(paymentData["conversionRate"])
	paymentResponse.PaymentData.Currency = fmt.Sprint(paymentData["currency"])
	paymentResponse.PaymentData.Date = fmt.Sprint(paymentData["date"])
	paymentResponse.PaymentData.Fee = fmt.Sprint(paymentData["fee"])
	paymentResponse.PaymentData.Media = fmt.Sprint(paymentData["media"])
	paymentResponse.PaymentData.Transferdate = fmt.Sprint(paymentData["transferDate"])

	// PendingInfo
	pendingInfo := response["pending_info"].(map[string]interface{})
	paymentResponse.PendingInfo.Date = fmt.Sprint(pendingInfo["date"])
	paymentResponse.PendingInfo.Media = fmt.Sprint(pendingInfo["media"])

	// RequestDate
	paymentResponse.RequestDate = response["requestDate"].(string)

	// Status
	paymentResponse.Status = fmt.Sprint(response["status"])

	// Subject
	paymentResponse.Subject = response["subject"].(string)

	return paymentResponse, nil
}
