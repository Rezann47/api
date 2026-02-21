package testutil

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewTestDB test için temiz bir PostgreSQL bağlantısı açar.
// DB_TEST_DSN env var'ı set edilmelidir.
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("DB_TEST_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=postgres password=test dbname=yks_test sslmode=disable"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("test db connect: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})

	return db
}

// TruncateTables test sonrası tabloları temizler
func TruncateTables(t *testing.T, db *gorm.DB, tables ...string) {
	t.Helper()
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Logf("truncate %s: %v", table, err)
		}
	}
}
