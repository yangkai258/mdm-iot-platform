// Package migrations contains data migration utilities for multi-tenancy.
//
// Run this script after applying the SQL migrations:
//
//	go run ./migrations/data_migration.go
//
// This script backfills tenant_id for all existing business tables,
// assigning them to the default tenant (00000000-0000-0000-0000-000000000001).
//
// It is idempotent: already-migrated rows (tenant_id IS NOT NULL) are skipped.
package main

import (
	"fmt"
	"log"
	"os"

	"mdm-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultTenantID = "00000000-0000-0000-0000-000000000001"

// backfillQuery executes an UPDATE ... WHERE tenant_id IS NULL and logs the affected rows.
func backfillQuery(db *gorm.DB, table, tenantID string) error {
	result := db.Exec(fmt.Sprintf(`UPDATE %s SET tenant_id = $1 WHERE tenant_id IS NULL`, table), tenantID)
	if result.Error != nil {
		return fmt.Errorf("backfill %s: %w", table, result.Error)
	}
	if result.RowsAffected > 0 {
		log.Printf("[data_migration] %s: %d rows updated", table, result.RowsAffected)
	} else {
		log.Printf("[data_migration] %s: no rows to update", table)
	}
	return nil
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback to individual env vars
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "mdm_user")
		pass := getEnv("DB_PASSWORD", "mdm_password")
		dbname := getEnv("DB_NAME", "mdm_db")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, pass, dbname)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("[data_migration] Starting tenant_id backfill...")
	log.Printf("[data_migration] Default tenant ID: %s", defaultTenantID)

	// Verify default tenant exists
	var count int64
	if err := db.Model(&models.Tenant{}).Where("id = ?", defaultTenantID).Count(&count).Error; err != nil {
		log.Fatalf("Failed to verify default tenant: %v", err)
	}
	if count == 0 {
		log.Fatalf("Default tenant %s does not exist. Please run 001_multi_tenant_fixed.sql first.", defaultTenantID)
	}
	log.Println("[data_migration] Default tenant verified.")

	// All tables that need tenant_id backfill (matching 001_multi_tenant_fixed.sql)
	tables := []string{
		// Core
		"sys_users",
		"devices",
		"members",
		// OTA
		"ota_packages",
		"ota_deployments",
		"ota_progress",
		"command_histories",
		// Alerts
		"device_alert_rules",
		"device_alerts",
		"geofence_rules",
		"geofence_alerts",
		"alert_notifications",
		// Notifications
		"notifications",
		"notification_templates",
		"announcements",
		// Policy & Compliance
		"policies",
		"policy_configs",
		"policy_bindings",
		"compliance_policies",
		"compliance_violations",
		// Apps
		"apps",
		"app_versions",
		"app_distributions",
		"app_install_records",
		"app_licenses",
		// Member commerce
		"coupons",
		"coupon_grants",
		"promotions",
		"stores",
		// Logs & Knowledge
		"sys_operation_logs",
		"sys_login_logs",
		"knowledge",
		// Pet
		"pet_profiles",
	}

	for _, table := range tables {
		if err := backfillQuery(db, table, defaultTenantID); err != nil {
			log.Fatalf("[data_migration] FATAL: %v", err)
		}
	}

	log.Println("[data_migration] Backfill complete.")

	// Initialize package_quotas for the default tenant if not already present
	var pqCount int64
	db.Model(&models.PackageQuota{}).Where("tenant_id = ?", defaultTenantID).Count(&pqCount)
	if pqCount == 0 {
		log.Println("[data_migration] Initializing package_quotas for default tenant...")
		// Get free package ID
		var freePkg models.Package
		if err := db.Where("package_code = ?", "free").First(&freePkg).Error; err != nil {
			log.Printf("[data_migration] WARNING: Could not find 'free' package: %v (skipping quota init)", err)
		} else {
			quotas := []struct {
				quotaType string
				limit     int
			}{
				{"user", 2},
				{"device", 5},
				{"store", 1},
				{"dept", 1},
				{"ota_deployment", 1},
				{"app", 1},
				{"notification", 100},
				{"alert", 50},
			}
			for _, q := range quotas {
				pq := models.PackageQuota{
					TenantID:   defaultTenantID,
					PackageID: freePkg.ID,
					QuotaType: q.quotaType,
					QuotaUsed: 0,
				}
				if err := db.Exec(`INSERT INTO package_quotas (tenant_id, package_id, quota_type, quota_limit, quota_used, updated_at)
					VALUES (?, ?, ?, ?, 0, NOW())
					ON CONFLICT (tenant_id, quota_type) DO NOTHING`,
					pq.TenantID, pq.PackageID, pq.QuotaType, q.limit).Error; err != nil {
					log.Printf("[data_migration] WARNING: Failed to insert quota %s: %v", q.quotaType, err)
				} else {
					log.Printf("[data_migration] package_quota[%s] = %d initialized", q.quotaType, q.limit)
				}
			}
		}
	} else {
		log.Printf("[data_migration] package_quotas already initialized (%d records)", pqCount)
	}

	log.Println("[data_migration] Done. All data migrated to default tenant.")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
