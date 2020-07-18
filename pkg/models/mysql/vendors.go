package mysql

import (
	"database/sql"

	"github.com/deldrone/server/pkg/models"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// VendorModel wraps a connection pool
type VendorModel struct {
	DB *sql.DB
}

// Insert creates a new vendor by inserting values into the database. Returns an error
func (m *VendorModel) Insert(name, address, email, password, phone string, pincode int, gpsLat, gpsLong float64) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO vendors (Vendor_name, vendor_address, vendor_pincode, vendor_email, vendor_hash_pwd, vendor_phone, vendor_gps_lat, vendor_gps_long)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = m.DB.Exec(stmt, name, address, pincode, email, hashedPassword, phone, gpsLat, gpsLong)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return models.ErrDuplicateEmail
			}
		}
	}

	return err
}

// Authenticate verifies the credentials and returns userid if valid details are provided.
func (m *VendorModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get fetches the details of the customer using its id
func (m *VendorModel) Get(id int) *models.Customer {
	return nil
}