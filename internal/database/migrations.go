package database

import (
    "velora/internal/models"

    "gorm.io/gorm"
)

// AutoMigrate registers all models and ensures schema consistency.
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
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
    )
}
