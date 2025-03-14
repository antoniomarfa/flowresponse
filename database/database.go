package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(dsn string) (*gorm.DB, error) {

	// Abrir la conexión con GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Obtener la instancia de *sql.DB subyacente para verificar la conexión
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Verificar la conexión con la base de datos
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
