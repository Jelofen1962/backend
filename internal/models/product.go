// backend/internal/models/product.go
package models

import (
	"database/sql"
	"time"
)

// Category corresponds to the "categories" table.
type Category struct {
	ID            string         `json:"id"`
	Name          string         `json:"name"`
	Description   sql.NullString `json:"description,omitempty"`
	NameEN        sql.NullString `json:"name_en,omitempty"`
	NameFI        sql.NullString `json:"name_fi,omitempty"`
	DescriptionEN sql.NullString `json:"description_en,omitempty"`
	DescriptionFI sql.NullString `json:"description_fi,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

// Product corresponds to the "products" table.
type Product struct {
	ID             string         `json:"id"`
	CategoryID     sql.NullString `json:"category_id,omitempty"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description,omitempty"`
	Price          float64        `json:"price"`
	InventoryCount int            `json:"inventory_count"`
	NameEN         sql.NullString `json:"name_en,omitempty"`
	NameFI         sql.NullString `json:"name_fi,omitempty"`
	DescriptionEN  sql.NullString `json:"description_en,omitempty"`
	DescriptionFI  sql.NullString `json:"description_fi,omitempty"`
	OriginEN       sql.NullString `json:"origin_en,omitempty"`
	OriginFI       sql.NullString `json:"origin_fi,omitempty"`
	UnitEN         sql.NullString `json:"unit_en,omitempty"`
	UnitFI         sql.NullString `json:"unit_fi,omitempty"`
	BadgeEN        sql.NullString `json:"badge_en,omitempty"`
	BadgeFI        sql.NullString `json:"badge_fi,omitempty"`
	FeaturesEN     sql.NullString `json:"features_en,omitempty"` // Assuming JSONB is read as a string
	FeaturesFI     sql.NullString `json:"features_fi,omitempty"` // We can unmarshal this later if needed
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// ProductImage corresponds to the "product_images" table.
type ProductImage struct {
	ID        string         `json:"id"`
	ProductID string         `json:"product_id"`
	URL       string         `json:"url"`
	AltText   sql.NullString `json:"alt_text,omitempty"`
	AltEN     sql.NullString `json:"alt_en,omitempty"`
	AltFI     sql.NullString `json:"alt_fi,omitempty"`
	IsPrimary bool           `json:"is_primary"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
