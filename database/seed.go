package database

import (
	"fmt"
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

    // Seed Payment Methods
    if err := seedPaymentMethods(db); err != nil {
        return fmt.Errorf("failed to seed payment methods: %w", err)
    }

    // Seed Default Admin
    if err := seedDefaultAdmin(db); err != nil {
        return fmt.Errorf("failed to seed admin user: %w", err)
    }

    // Seed Sample Studios (optional, untuk testing)
    if err := seedSampleStudios(db); err != nil {
        return fmt.Errorf("failed to seed sample studios: %w", err)
    }

    log.Println("✅ Database seeding completed!")
    return nil
}

// seedPaymentMethods inserts default payment methods
func seedPaymentMethods(db *gorm.DB) error {
    var count int64
    db.Model(&PaymentMethod{}).Count(&count)

    // Only seed if table is empty
    if count > 0 {
        log.Println("⏭️  Payment methods already exist, skipping...")
        return nil
    }

    log.Println("📦 Seeding payment methods...")

    paymentMethods := []PaymentMethod{
        {
            Name:          "BCA",
            BankName:      "Bank Central Asia",
            AccountNumber: "1234567890",
            AccountName:   "Studio Booking System",
            IsActive:      true,
        },
        {
            Name:          "Mandiri",
            BankName:      "Bank Mandiri",
            AccountNumber: "0987654321",
            AccountName:   "Studio Booking System",
            IsActive:      true,
        },
        {
            Name:          "BNI",
            BankName:      "Bank Negara Indonesia",
            AccountNumber: "5555666677",
            AccountName:   "Studio Booking System",
            IsActive:      true,
        },
        {
            Name:          "BRI",
            BankName:      "Bank Rakyat Indonesia",
            AccountNumber: "8888999900",
            AccountName:   "Studio Booking System",
            IsActive:      true,
        },
    }

    if err := db.Create(&paymentMethods).Error; err != nil {
        return err
    }

    log.Printf("✅ Created %d payment methods\n", len(paymentMethods))
    return nil
}

// seedDefaultAdmin creates default admin user
func seedDefaultAdmin(db *gorm.DB) error {
    var count int64
    db.Model(&User{}).Where("role = ?", "admin").Count(&count)

    // Only create if no admin exists
    if count > 0 {
        log.Println("⏭️  Admin user already exists, skipping...")
        return nil
    }

    log.Println("👤 Creating default admin user...")

    admin := User{
        Name:     "Admin",
        Email:    "admin@studiobooking.com",
        Password: hashPassword("admin123"), // Change this in production!
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

    // Only seed if table is empty
    if count > 0 {
        log.Println("⏭️  Studios already exist, skipping...")
        return nil
    }

    log.Println("🎸 Seeding sample studios...")

    studios := []Studio{
        {
            Name:        "Studio Rock A",
            Description: "Studio musik profesional dengan peralatan lengkap untuk band rock. Dilengkapi dengan drum set, amplifier, dan sistem akustik berkualitas tinggi.",
            Location:    "Jakarta Selatan",
            PricePerHour: 150000,
            ImageURL:    "https://images.unsplash.com/photo-1598488035139-bdbb2231ce04",
            Facilities: StringArray{
                "AC",
                "Drum Set",
                "Guitar Amplifier",
                "Bass Amplifier",
                "Microphone",
                "Soundproof",
                "Mixing Console",
            },
            OperatingHours: "09:00-22:00",
            IsActive:       true,
        },
        {
            Name:        "Studio Acoustic B",
            Description: "Studio akustik nyaman untuk recording vokal dan instrumen akustik. Ideal untuk singer-songwriter dan podcast.",
            Location:    "Jakarta Pusat",
            PricePerHour: 100000,
            ImageURL:    "https://images.unsplash.com/photo-1519508234439-4f23643125c1",
            Facilities: StringArray{
                "AC",
                "Acoustic Guitar",
                "Condenser Mic",
                "Audio Interface",
                "Soundproof",
                "Comfortable Seating",
            },
            OperatingHours: "10:00-21:00",
            IsActive:       true,
        },
        {
            Name:        "Studio Premium C",
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
            Location:    "Jakarta Timur",
            PricePerHour: 75000,
            ImageURL:    "https://images.unsplash.com/photo-1563330232-57114bb0823c",
            Facilities: StringArray{
                "Fan",
                "Basic Drum Set",
                "Guitar Amplifier",
                "Bass Amplifier",
                "Microphone",
            },
            OperatingHours: "12:00-20:00",
            IsActive:       true,
        },
        {
            Name:        "Studio Recording E",
            Description: "Studio khusus recording dengan engineer berpengalaman. Hasil rekaman berkualitas profesional.",
            Location:    "Tangerang",
            PricePerHour: 200000,
            ImageURL:    "https://images.unsplash.com/photo-1598653222000-6b7b7a552625",
            Facilities: StringArray{
                "AC",
                "Recording Engineer",
                "Professional Microphones",
                "Audio Interface",
                "Mixing & Mastering",
                "Soundproof Recording Booth",
                "Monitoring Speakers",
            },
            OperatingHours: "09:00-22:00",
            IsActive:       true,
        },
    }

    if err := db.Create(&studios).Error; err != nil {
        return err
    }

    log.Printf("✅ Created %d sample studios\n", len(studios))
    return nil
}

// SeedTestBookings creates sample booking data for development/testing
// Call this manually if needed: go run main.go seed-test
func SeedTestBookings(db *gorm.DB) error {
    log.Println("📅 Seeding test bookings...")

    // Get first user and studio
    var user User
    var studio Studio

    if err := db.First(&user).Error; err != nil {
        return fmt.Errorf("no user found: %w", err)
    }

    if err := db.First(&studio).Error; err != nil {
        return fmt.Errorf("no studio found: %w", err)
    }

    // Sample bookings (you can customize these)
    // Note: Implement this based on your needs

    log.Println("✅ Test bookings created!")
    return nil
}