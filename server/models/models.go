package models

import "time"

type User struct {
    UserID          string      `json:"user_id"`
    FirstName       string      `json:"first_name"`
    LastName        string      `json:"last_name"`
    Email           string      `json:"email"`
    PhoneNumber     string      `json:"phone_number"`
    Password        string      `json:"password"`
    CreatedAt       time.Time   `json:"created_at"`
}

type Organization struct {
    OrganizationID string    `json:"organization_id"`
    Name           string    `json:"name"`
    CreatedAt      time.Time `json:"created_at"`
}

type UserOrganization struct {
    ID             string    `json:"id"`
    UserID         string    `json:"user_id"`
    OrganizationID string    `json:"organization_id"`
    Role           string    `json:"role"`
    CreatedAt      time.Time `json:"created_at"`
}

type Folder struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    OrganizationID  string    `json:"organization_id"`
    ParentFolderID  *string    `json:"parent_folder_id,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    Deleted         bool      `json:"deleted"`
    DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type File struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    FolderID        *string    `json:"folder_id"`
    FilePath        string    `json:"file_path"`
    FileSize        int64     `json:"file_size"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
    OrganizationID  string    `json:"organization_id"`
    Deleted         bool      `json:"deleted"`
    DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

type Fleet struct {
    ID             string    `json:"id"`
    Name           string    `json:"name"`
    OrganizationID string    `json:"organization_id"`
    CreatedAt      time.Time `json:"created_at"`
}

type Device struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    SerialNumber  string    `json:"serial_number"`
    FleetID       *string    `json:"fleet_id"`
    IPAddress     string    `json:"ip_address"`
    OrganizationID string   `json:"organization_id"`
    CreatedAt     time.Time `json:"created_at"`
}

type UserDevice struct {
    ID                string    `json:"id"`
    UserOrganizationID string    `json:"user_organization_id"`
    DeviceID          string    `json:"device_id"`
    CreatedAt         time.Time `json:"created_at"`
}

type Plan struct {
    ID            string  `json:"id"`
    Name          string  `json:"name"`
    Description   string  `json:"description"`
    Price         float64 `json:"price"`
    Interval      string  `json:"interval"`
    StorageLimit  int64   `json:"storage_limit"`
}

type Subscription struct {
    ID                   string    `json:"id"`
    OrganizationID       string    `json:"organization_id"`
    PlanID               string    `json:"plan_id"`
    Status               string    `json:"status"`
    StripeSubscriptionID string    `json:"stripe_subscription_id"`
    StartDate            time.Time `json:"start_date"`
    EndDate              time.Time `json:"end_date"`
    PaypalSubscriptionID string    `json:"paypal_subscription_id"`
    CreatedAt            time.Time `json:"created_at"`
    UpdatedAt            time.Time `json:"updated_at"`
}
