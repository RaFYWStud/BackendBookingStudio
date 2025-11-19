package database

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// User model - sudah ada, tambahkan UpdatedAt
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

// BookingStatus enum
type BookingStatus string

const (
    BookingStatusPending   BookingStatus = "pending"
    BookingStatusConfirmed BookingStatus = "confirmed"
    BookingStatusPaid      BookingStatus = "paid"
    BookingStatusCompleted BookingStatus = "completed"
    BookingStatusCancelled BookingStatus = "cancelled"
    BookingStatusExpired   BookingStatus = "expired"
)

// Booking model
type Booking struct {
    ID              int           `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    UserID          int           `gorm:"column:user_id;not null;index"`
    User            *User         `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    StudioID        int           `gorm:"column:studio_id;not null;index"`
    Studio          *Studio       `gorm:"foreignKey:StudioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    BookingDate     time.Time     `gorm:"column:booking_date;type:date;not null;index"`
    StartTime       time.Time     `gorm:"column:start_time;type:time;not null"`
    EndTime         time.Time     `gorm:"column:end_time;type:time;not null"`
    DurationHours   int           `gorm:"column:duration_hours;not null"`
    TotalPrice      int           `gorm:"column:total_price;not null"`
    DPAmount        int           `gorm:"column:dp_amount;not null"` // 30% of total
    RemainingAmount int           `gorm:"column:remaining_amount;not null"`
    DPDeadline      time.Time     `gorm:"column:dp_deadline;index"` // 24h after created
    Status          BookingStatus `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
    CreatedAt       time.Time     `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt       time.Time     `gorm:"column:updated_at;autoUpdateTime"`
}

// Add composite index for availability check
func (Booking) TableName() string {
    return "bookings"
}

// PaymentMethod model
type PaymentMethod struct {
    ID            int       `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    Name          string    `gorm:"column:name;type:varchar(100);not null"`
    BankName      string    `gorm:"column:bank_name;type:varchar(100);not null"`
    AccountNumber string    `gorm:"column:account_number;type:varchar(50);not null"`
    AccountName   string    `gorm:"column:account_name;type:varchar(255);not null"`
    IsActive      bool      `gorm:"column:is_active;default:true"`
    CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// PaymentType enum
type PaymentType string

const (
    PaymentTypeDP   PaymentType = "dp"
    PaymentTypeFull PaymentType = "full_payment"
)

// PaymentStatus enum
type PaymentStatus string

const (
    PaymentStatusPending  PaymentStatus = "pending"
    PaymentStatusVerified PaymentStatus = "verified"
    PaymentStatusRejected PaymentStatus = "rejected"
)

// Payment model
type Payment struct {
    ID               int            `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    BookingID        int            `gorm:"column:booking_id;not null;index"`
    Booking          *Booking       `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    PaymentType      PaymentType    `gorm:"column:payment_type;type:varchar(20);not null;index"`
    Amount           int            `gorm:"column:amount;not null"`
    ProofURL         string         `gorm:"column:proof_url;type:text"`
    Status           PaymentStatus  `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
    VerifiedBy       *int           `gorm:"column:verified_by"`
    VerifiedByUser   *User          `gorm:"foreignKey:VerifiedBy;references:ID"`
    VerifiedAt       *time.Time     `gorm:"column:verified_at"`
    RejectionReason  string         `gorm:"column:rejection_reason;type:text"`
    PaidAt           *time.Time     `gorm:"column:paid_at"`
    CreatedAt        time.Time      `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt        time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

// Cancellation model
type Cancellation struct {
    ID           int       `gorm:"column:id;primaryKey;autoIncrement;not null;<-:create"`
    BookingID    int       `gorm:"column:booking_id;not null;uniqueIndex"`
    Booking      *Booking  `gorm:"foreignKey:BookingID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Reason       string    `gorm:"column:reason;type:text"`
    RefundAmount int       `gorm:"column:refund_amount;not null;default:0"`
    RefundMethod string    `gorm:"column:refund_method;type:varchar(100)"`
    RefundStatus string    `gorm:"column:refund_status;type:varchar(50);default:'pending'"` // pending, processed, completed
    CancelledAt  time.Time `gorm:"column:cancelled_at;not null"`
    CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
}