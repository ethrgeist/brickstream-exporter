package database

import (
	"github.com/ethrgeist/brickstream-exporter/models"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"reflect"
)

var (
	DB *gorm.DB
)

func registerIDGenerator(db *gorm.DB) {
	db.Callback().Create().Before("gorm:before_create").Register("custom:before_create", func(tx *gorm.DB) {
		if tx.Statement.Schema == nil {
			return
		}

		pkField := tx.Statement.Schema.PrioritizedPrimaryField
		if pkField == nil || pkField.FieldType.Kind() != reflect.String {
			return
		}

		ctx := tx.Statement.Context
		rv := tx.Statement.ReflectValue

		switch rv.Kind() {
		case reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				if elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if pkValue, _ := pkField.ValueOf(ctx, elem); pkValue == "" {
					id, _ := gonanoid.New(14)
					pkField.Set(ctx, elem, id)
				}
			}
		case reflect.Struct:
			if pkValue, _ := pkField.ValueOf(ctx, rv); pkValue == "" {
				id, _ := gonanoid.New(14)
				pkField.Set(ctx, rv, id)
			}
		default:
			return
		}
	})
}

func DbConn() error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		models.Site{},
		models.Device{},
		models.Counter{},
	)
	if err != nil {
		return err
	}

	registerIDGenerator(db)

	DB = db
	return nil
}
