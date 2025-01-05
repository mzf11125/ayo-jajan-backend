package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	_ "github.com/lib/pq"
)

// BLE setup
func setupBLE() (ble.Client, error) {
	device, err := ble.NewDevice()
	if err != nil {
		return nil, fmt.Errorf("could not initialize BLE device: %v", err)
	}

	// Scan for nearby BLE devices (Tile, AirTag, etc.)
	devices, err := device.Scan(10 * time.Second)
	if err != nil {
		return nil, fmt.Errorf("error scanning BLE devices: %v", err)
	}

	// Assuming the device's MAC address is known
	var targetDevice ble.Device
	for _, dev := range devices {
		if dev.Address == "MAC_ADDRESS" {
			targetDevice = dev
			break
		}
	}

	if targetDevice == nil {
		return nil, fmt.Errorf("device not found")
	}

	client, err := device.Connect(targetDevice)
	if err != nil {
		return nil, fmt.Errorf("could not connect to device: %v", err)
	}

	return client, nil
}

// Database setup
func setupDB() (*sql.DB, error) {
	connStr := "user=youruser dbname=streetfood password=yourpassword host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("could not connect to the database: %v", err)
	}

	return db, nil
}

// Store Vendor Location
func storeLocation(db *sql.DB, name string, lat, lon float64) error {
	query := `
		INSERT INTO street_food_vendors (name, location)
		VALUES ($1, ST_SetSRID(ST_Point($2, $3), 4326))`
	_, err := db.Exec(query, name, lat, lon)
	if err != nil {
		return fmt.Errorf("could not insert vendor data: %v", err)
	}
	return nil
}

// Tracking function
func trackVendor() {
	// Setup BLE and database
	client, err := setupBLE()
	if err != nil {
		log.Fatalf("BLE Setup Error: %v", err)
	}
	defer client.Close()

	db, err := setupDB()
	if err != nil {
		log.Fatalf("Database Setup Error: %v", err)
	}
	defer db.Close()

	// Example: Track vendor for 5 minutes
	for {
		// Fetch Vendor Location (simulating for now)
		// Normally, you'd get these coordinates from the BLE device's data.
		lat, lon := 40.7128, -74.0060 // Sample coordinates (New York)

		// Store vendor info in the DB
		err = storeLocation(db, "Vendor Name", lat, lon)
		if err != nil {
			log.Printf("Error storing vendor location: %v", err)
		}

		// Wait for the next cycle
		time.Sleep(10 * time.Second) // Adjust the time as needed
	}
}

func main() {
	// Start tracking vendors
	trackVendor()
}
