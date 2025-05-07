package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type TokenRequest struct {
	Token string `json:"token"`
}

type TokenResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Create---Req  request struct
type CreatePagosReq struct {
	ID            string    `gorm:"primaryKey;autoIncrement"`
	Tipocom       string    `json:"tipocom"`
	IngresoId     int64     `json:"ingreso_id"`
	Identificador string    `json:"identificador"`
	Fecha         time.Time `json:"fecha"`
	SaleId        int64     `json:"sale_id"`
	Rutalumn      string    `json:"rutalumn"`
	Transaccion   string    `json:"transaccion"`
	Tipo          string    `json:"tipo"`
	Monto         float64   `json:"monto"`
	Nrotarjeta    string    `json:"nrotarjeta"`
	Codigoauto    string    `json:"codigoauto"`
	Fechaauto     time.Time `json:"fechaauto"`
	Tipopago      string    `json:"tipopago"`
	Nrocuota      int       `json:"nrocuota"`
	Fechatransac  time.Time `json:"fechatransac"`
	Author        string    `json:"author"`
	Activo        int       `json:"activo"`
	CompanyId     int64     `json:"company_id"`
	Cuotapagada   int       `json:"cuotapagada"`
	Cuotafecha    string    `json:"cuotafecha"`
	CreatedDate   time.Time `gorm:"autoCreateTime"`
	UpdatedDate   time.Time `gorm:"autoUpdateTime"`
}

func (CreatePagosReq) TableName() string {
	return "pagos" // Nombre de la tabla en la base de datos
}

type IngresoReq struct {
	ID            int64     `json:"-"`
	Tipocomp      string    `json:"tipocomp"`
	Fecha         time.Time `json:"fecha"`
	Identificador string    `json:"identificador"`
	SaleId        int64     `json:"sale_id"`
	CursoId       int64     `json:"curso_id"`
	Rutapo        string    `json:"rutapo"`
	Rutalum       string    `json:"rutalum"`
	Fpago         string    `json:"fpago"`
	Monto         float32   `json:"monto"`
	Activo        int       `json:"activo"`
	StatusPago    string    `json:"status_pago"`
	Author        string    `json:"author"`
	CompanyId     int64     `json:"company_id"`
	TokenFlow     string    `json:"token_flow"`
	Nrocuotas     int       `json:"nrocuotas"`
	Valorcuota    float32   `json:"valorcuota"`
	Fechainicial  time.Time `json:"fechainicial"`
	CreatedDate   time.Time `gorm:"autoCreateTime"`
	UpdatedDate   time.Time `gorm:"autoUpdateTime"`
}

func (IngresoReq) TableName() string {
	return "ingresos" // Nombre de la tabla en la base de datos
}

type UpdateIngresoReq struct {
	ID            string     `json:"-"`
	Tipocomp      *string    `json:"tipocomp"`
	Fecha         *time.Time `json:"fecha"`
	Identificador *string    `json:"identificador"`
	SaleId        *int64     `json:"sale_id"`
	CursoId       *int64     `json:"curso_id"`
	Rutapo        *string    `json:"rutapo"`
	Rutalum       *string    `json:"rutalum"`
	Fpago         *string    `json:"fpago"`
	Monto         *float32   `json:"monto"`
	Activo        *int       `json:"activo"`
	StatusPago    *string    `json:"status_pago"`
	Author        *string    `json:"author"`
	CompanyId     *int64     `json:"company_id"`
	TokenFlow     *string    `json:"token_flow"`
	Nrocuotas     *int       `json:"nrocuotas"`
	Valorcuota    *float32   `json:"valorcuota"`
	Fechainicial  *time.Time `json:"fechainicial"`
	CreatedDate   *time.Time `gorm:"autoCreateTime"`
	UpdatedDate   *time.Time `gorm:"autoUpdateTime"`
}

func (UpdateIngresoReq) TableName() string {
	return "ingresos" // Nombre de la tabla en la base de datos
}

type PaymentResponse struct {
	Amount        string `json:"amount"`
	CommerceOrder string `json:"commerceOrder"`
	Currency      string `json:"currency"`
	FlowOrder     string `json:"flowOrder"`
	Merchantid    string `json:"merchantId"`
	Optional      struct {
		Venta  string `json:"venta"`
		Alumno string `json:"alumno"`
	} `json:"optional"`
	Payer       string `json:"payer"`
	PaymentData struct {
		Amount         string `json:"amount"`
		Balance        string `json:"balance"`
		Conversiondate string `json:"conversionDate"`
		Conversionrate string `json:"conversionRate"`
		Currency       string `json:"currency"`
		Date           string `json:"date"`
		Fee            string `json:"fee"`
		Media          string `json:"media"`
		Transferdate   string `json:"transferDate"`
	} `json:"paymentData"`
	PendingInfo struct {
		Date  string `json:"date"`
		Media string `json:"media"`
	} `json:"pending_info"`
	RequestDate string `json:"requestDate"`
	Status      string `json:"status"`
	Subject     string `json:"subject"`
}

type CompanyconResp struct {
	ID            string `json:"id"`
	CompanyId     int64  `json:"company_id"`
	ActiveFlow    int    `json:"active_flow"`
	ApikeyFlow    string `json:"apikey_flow"`
	SecretkeyFlow string `json:"secretkey_flow"`
}

func (CompanyconResp) TableName() string {
	return "company_connection" // Nombre de la tabla en la base de datos
}

type GatewaysResp struct {
	ID               string                  `json:"id"`
	CompanyId        int64                   `json:"company_id"`
	GatewayId        int64                   `json:"gateway_id"`
	AdditionalConfig GatewayAdditionalConfig `gorm:"type:jsonb" json:"additional_config"` // GORM serializa automáticamente
	Active           int                     `json:"active"`
	CreatedDate      time.Time               `gorm:"autoCreateTime"`
	UpdatedDate      time.Time               `gorm:"autoUpdateTime"`
}

func (GatewaysResp) TableName() string {
	return "gateways" // Nombre de la tabla en la base de datos
}

type GatewayAdditionalConfig struct {
	FlowAPIKey         string `json:"flow_apikey"`
	FlowSecretKey      string `json:"flow_secretkey"`
	TrbkCommercialCode string `json:"trbk_commercialcode"`
	TrbkKeySecret      string `json:"trbk_keysecret"`
}

func (c GatewayAdditionalConfig) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *GatewayAdditionalConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
