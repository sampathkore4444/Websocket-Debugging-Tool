package core

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	db *gorm.DB
}

func NewDatabase(config *Config) (*Database, error) {
	// Check if running in Docker
	if os.Getenv("DATABASE_URL") != "" {
		// Use DATABASE_URL for Docker/Production
		return NewDatabaseFromURL(os.Getenv("DATABASE_URL"))
	}

	// Use individual config values for local development
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	user := config.GetString("database.user")
	password := config.GetString("database.password")
	dbname := config.GetString("database.name")
	sslmode := config.GetString("database.sslmode")

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	return connectToPostgres(dsn, config)
}

func NewDatabaseFromURL(databaseURL string) (*Database, error) {
	return connectToPostgres(databaseURL, nil)
}

func connectToPostgres(dsn string, config *Config) (*Database, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB for connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Production-grade connection pool settings
	sqlDB.SetMaxOpenConns(25)                  // Maximum open connections
	sqlDB.SetMaxIdleConns(10)                  // Maximum idle connections
	sqlDB.SetConnMaxLifetime(5 * 60 * 1000)    // Connection lifetime: 5 minutes
	sqlDB.SetConnMaxIdleTime(1 * 60 * 1000)   // Idle connection timeout: 1 minute

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to PostgreSQL database")

	return &Database{db: db}, nil
}

func (d *Database) GetDB() *gorm.DB {
	return d.db
}

func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	err := db.AutoMigrate(
		&Session{},
		&Message{},
		&Connection{},
		&FuzzTest{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
