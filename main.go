package main

import (
	"fmt"
	"os"

	"github.com/RaFYWStud/BackendBookingStudio/config"
	dbConfig "github.com/RaFYWStud/BackendBookingStudio/config/database"
	"github.com/RaFYWStud/BackendBookingStudio/config/pkg/token"
	"github.com/RaFYWStud/BackendBookingStudio/config/server"
	dbMigration "github.com/RaFYWStud/BackendBookingStudio/database"
)

func main() {
    config.Load()
    token.Load()

    // Handle CLI commands
    if len(os.Args) > 1 {
        switch os.Args[1] {
        case "migrate":
            runMigrations()
            return
        case "reset":
            runReset()
            return
        case "seed":
            runSeedOnly()
            return
        default:
            fmt.Println("Unknown command. Use: migrate | reset | seed")
            return
        }
    }

    // 🔥 AUTO-MIGRATE saat server start
    db, _, err := dbConfig.ConnectDB()
    if err != nil {
        panic(fmt.Errorf("failed to connect database: %w", err))
    }

    // Check if migration needed
    if !db.Migrator().HasTable(&dbMigration.User{}) {
        fmt.Println("🔄 First run detected, running auto-migration...")
        if err := dbMigration.RunMigration(db); err != nil {
            panic(fmt.Errorf("auto-migration failed: %w", err))
        }
    } else {
        fmt.Println("✅ Database already migrated")
    }

    server.Run()
}

func runMigrations() {
    db, _, err := dbConfig.ConnectDB()
    if err != nil {
        panic(err)
    }

    if err := dbMigration.RunMigration(db); err != nil {
        panic(err)
    }
}

func runReset() {
    db, _, err := dbConfig.ConnectDB()
    if err != nil {
        panic(err)
    }

    fmt.Println("🗑️  Dropping all tables...")
    err = db.Migrator().DropTable(
        &dbMigration.Booking{},
        &dbMigration.Studio{},
        &dbMigration.User{},
    )
    if err != nil {
        panic(err)
    }

    fmt.Println("🔄 Recreating tables with AutoMigrate...")
    if err := dbMigration.RunMigration(db); err != nil {
        panic(err)
    }

    fmt.Println("✅ Database reset completed")
}

func runSeedOnly() {
    db, _, err := dbConfig.ConnectDB()
    if err != nil {
        panic(err)
    }

    fmt.Println("🌱 Running seed only...")
    if err := dbMigration.Seed(db); err != nil {
        panic(err)
    }
    fmt.Println("✅ Seeding completed")
}