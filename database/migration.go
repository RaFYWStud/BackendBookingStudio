package database

import (
	"fmt"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
    fmt.Println("🚀 Running migrations...")
    
    if err := db.AutoMigrate(
        &User{},
        &Studio{},
        &Booking{},
    ); err != nil {
        return fmt.Errorf("gagal migrasi: %w", err)
    }
    
    fmt.Println("✅ Migrations completed")

    // Add composite index for bookings
    if err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_bookings_studio_date_time 
        ON bookings(studio_id, booking_date, start_time, end_time)
    `).Error; err != nil {
        fmt.Printf("⚠️  Warning: Failed to create composite index: %v\n", err)
    }

    fmt.Println("🌱 Seeding database...")
    if err := Seed(db); err != nil {
        return fmt.Errorf("gagal seeding: %w", err)
    }

    return nil
}