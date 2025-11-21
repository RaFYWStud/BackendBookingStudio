package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hashPassword(password string) string {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic("failed to hash password: " + err.Error())
    }
    return string(hash)
}

func Seed(db *gorm.DB) error {
    log.Println("🌱 Starting database seeding...")

    // Seed Default Admin
    if err := seedDefaultAdmin(db); err != nil {
        return err
    }

    // Seed Sample Studios (optional, untuk testing)
    if err := seedSampleStudios(db); err != nil {
        return err
    }

    log.Println("✅ Database seeding completed!")
    return nil
}

// seedDefaultAdmin creates default admin user
func seedDefaultAdmin(db *gorm.DB) error {
    var count int64
    db.Model(&User{}).Where("role = ?", "admin").Count(&count)

    if count > 0 {
        log.Println("⏭️  Admin user already exists, skipping...")
        return nil
    }

    log.Println("👤 Creating default admin user...")

    admin := User{
        Name:     "Admin",
        Email:    "admin@studiobooking.com",
        Password: hashPassword("admin123"),
        Role:     "admin",
    }

    if err := db.Create(&admin).Error; err != nil {
        return err
    }

    log.Println("✅ Admin user created:")
    log.Println("   Email: admin@studiobooking.com")
    log.Println("   Password: admin123")
    log.Println("   ⚠️  CHANGE THIS PASSWORD IN PRODUCTION!")

    return nil
}

// seedSampleStudios creates sample studio data for testing
func seedSampleStudios(db *gorm.DB) error {
    var count int64
    db.Model(&Studio{}).Count(&count)

    if count > 0 {
        log.Println("⏭️  Studios already exist, skipping...")
        return nil
    }

    log.Println("🎸 Creating sample studios...")

    studios := []Studio{
        {
            Name:        "Studio Premium A",
            Description: "Studio premium dengan peralatan kelas dunia. Cocok untuk produksi musik profesional dan rekaman album.",
            Location:    "Jakarta Barat",
            PricePerHour: 250000,
            ImageURL:    "https://images.unsplash.com/photo-1598653222000-6b7b7a552625",
            Facilities: StringArray{
                "AC",
                "Professional Drum Set",
                "Multiple Amplifiers",
                "Grand Piano",
                "Pro Microphones",
                "Soundproof",
                "Professional Mixing Console",
                "Recording Booth",
                "Lounge Area",
            },
            OperatingHours: "08:00-23:00",
            IsActive:       true,
        },
        {
            Name:        "Studio Budget D",
            Description: "Studio terjangkau untuk latihan band dan jamming session. Peralatan standar dengan harga bersahabat.",
            Location:    "Jakarta Utara",
            PricePerHour: 100000,
            ImageURL:    "https://images.unsplash.com/photo-1519892300165-cb5542fb47c7",
            Facilities: StringArray{
                "AC",
                "Drum Set",
                "Guitar Amplifier",
                "Bass Amplifier",
                "Microphones",
                "Mixing Console",
            },
            OperatingHours: "09:00-21:00",
            IsActive:       true,
        },
    }

    if err := db.Create(&studios).Error; err != nil {
        return err
    }

    log.Printf("✅ Created %d sample studios\n", len(studios))
    return nil
}