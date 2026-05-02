package database

import (
    "fmt"
    "time"

    "velora/config"
    "velora/internal/models"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// Connect opens a PostgreSQL database connection and configures GORM.
func Connect(cfg *config.Config) (*gorm.DB, error) {
    dialector := postgres.Open(cfg.DatabaseURL)
    db, err := gorm.Open(dialector, &gorm.Config{
        Logger: logger.Default.LogMode(logger.Warn),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetConnMaxLifetime(30 * time.Minute)

    if err := AutoMigrate(db); err != nil {
        return nil, err
    }

    return db, nil
}

// Migrate applies schema changes for all tenant-aware models.
func Migrate(db *gorm.DB) error {
    if err := db.AutoMigrate(
        &models.Tenant{},
        &models.User{},
        &models.APIKey{},
        &models.Agent{},
        &models.Workflow{},
        &models.WorkflowRun{},
        &models.Lead{},
        &models.Contact{},
        &models.Note{},
        &models.Integration{},
        &models.Subscription{},
        &models.Invoice{},
        &models.UsageRecord{},
        &models.Partner{},
        &models.Referral{},
        &models.AuditLog{},
        &models.Notification{},
    ); err != nil {
        return fmt.Errorf("database migration failed: %w", err)
    }
    return nil
}
