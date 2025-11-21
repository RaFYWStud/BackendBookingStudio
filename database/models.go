package database

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// User model
type User struct {
    ID        int       `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    Name      string    `gorm:"column:name;not null"`
    Email     string    `gorm:"column:email;uniqueIndex;not null"`
    Password  string    `gorm:"column:password;not null"`
    Role      string    `gorm:"column:role;type:varchar(50);not null;default:'customer'"` // customer, admin
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// StringArray type for JSONB arrays
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
    return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
    if value == nil {
        *a = []string{}
        return nil
    }
    bytes, ok := value.([]byte)
    if !ok {
        return nil
    }
    return json.Unmarshal(bytes, a)
}

// Studio model
type Studio struct {
    ID             int         `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    Name           string      `gorm:"column:name;not null"`
    Description    string      `gorm:"column:description;type:text"`
    Location       string      `gorm:"column:location;not null;index"`
    PricePerHour   int         `gorm:"column:price_per_hour;not null;index"`
    ImageURL       string      `gorm:"column:image_url;type:text"`
    Facilities     StringArray `gorm:"column:facilities;type:jsonb"`
    OperatingHours string      `gorm:"column:operating_hours;type:varchar(100)"` // '09:00-22:00'
    IsActive       bool        `gorm:"column:is_active;default:true;index"`
    CreatedAt      time.Time   `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt      time.Time   `gorm:"column:updated_at;autoUpdateTime"`
}

// BookingStatus enum - SIMPLIFIED
type BookingStatus string

const (
    BookingStatusPending   BookingStatus = "pending"   // Menunggu pembayaran
    BookingStatusConfirmed BookingStatus = "confirmed" // Sudah bayar (dikonfirmasi admin)
    BookingStatusCompleted BookingStatus = "completed" // Selesai digunakan
    BookingStatusCancelled BookingStatus = "cancelled" // Dibatalkan
)

// Booking model - SIMPLIFIED
type Booking struct {
    ID            int           `gorm:"primaryKey;autoIncrement" json:"id"`
    UserID        int           `gorm:"not null;index" json:"user_id"`
    StudioID      int           `gorm:"not null;index" json:"studio_id"`
    BookingDate   time.Time     `gorm:"type:date;not null;index" json:"booking_date"`
    StartTime     time.Time     `gorm:"type:time;not null" json:"start_time"`
    EndTime       time.Time     `gorm:"type:time;not null" json:"end_time"`
    DurationHours int           `gorm:"not null" json:"duration_hours"`
    TotalPrice    int           `gorm:"not null" json:"total_price"`
    Status        BookingStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
    AdminNotes    string        `gorm:"type:text" json:"admin_notes"` // Catatan pembayaran dari admin
    CreatedAt     time.Time     `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt     time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

    // Relations
    User   *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
    Studio *Studio `gorm:"foreignKey:StudioID;constraint:OnDelete:CASCADE" json:"studio,omitempty"`
}

func (Booking) TableName() string {
    return "bookings"
}