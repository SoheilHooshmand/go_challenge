package main

import (
    "math"
    "os"
    "testing"
    "time"
)

const tolerance = 0.01 // Set a small tolerance for floating-point comparisons

// Unit test for haversine formula
func TestHaversine(t *testing.T) {
    lat1, lon1 := 51.5007, 0.1246
    lat2, lon2 := 40.6892, 74.0445
    expectedDistance := 5574.840456848555 // Expected distance in kilometers

    result := haversine(lat1, lon1, lat2, lon2)
    if math.Abs(result-expectedDistance) > tolerance {
        t.Errorf("Haversine calculation incorrect, got: %f, expected: %f", result, expectedDistance)
    }
}

// Unit test for fare calculation based on time of day
func TestTimeOfDayFare(t *testing.T) {
    // Test during daytime
    daytimeTimestamp := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC).Unix()
    if rate := timeOfDayFare(daytimeTimestamp); rate != movingRateDay {
        t.Errorf("Daytime rate incorrect, got: %f, expected: %f", rate, movingRateDay)
    }

    // Test during nighttime
    nighttimeTimestamp := time.Date(2023, 10, 1, 2, 0, 0, 0, time.UTC).Unix()
    if rate := timeOfDayFare(nighttimeTimestamp); rate != movingRateNight {
        t.Errorf("Nighttime rate incorrect, got: %f, expected: %f", rate, movingRateNight)
    }
}

// Unit test for calculateFare function
func TestCalculateFare(t *testing.T) {
    points := []DeliveryPoint{
        {ID: "1", Lat: 51.5007, Lng: 0.1246, Timestamp: 1696111200}, // Day
        {ID: "1", Lat: 40.6892, Lng: 74.0445, Timestamp: 1696114800}, // 1-hour later, 1318.39 km apart
    }

    expectedFare := flagFare + (5574.840456848555 * movingRateDay)
    fare := calculateFare(points)
    if math.Abs(fare-expectedFare) > tolerance {
        t.Errorf("Fare calculation incorrect, got: %f, expected: %f", fare, expectedFare)
    }
}

// End-to-end test for fare calculation and CSV writing
func TestFareCalculationAndCSVOutput(t *testing.T) {
    // Mock the input data
    deliveries := map[string][]DeliveryPoint{
        "1": {
            {ID: "1", Lat: 52.2296756, Lng: 21.0122287, Timestamp: 1696111200},
            {ID: "1", Lat: 41.8919300, Lng: 12.5113300, Timestamp: 1696114800},
        },
    }

    // Perform fare calculation and output
    fareEstimates := make(map[string]float64)
    for id, points := range deliveries {
        fare := calculateFare(points)
        fareEstimates[id] = fare
    }

    // Write to CSV (mock file)
    err := writeToCSV("test_fare_estimates.csv", fareEstimates)
    if err != nil {
        t.Errorf("Error writing CSV: %v", err)
    }

    // Check if the file was created
    if _, err := os.Stat("test_fare_estimates.csv"); os.IsNotExist(err) {
        t.Errorf("CSV file was not created")
    } else {
        // Clean up the test file
        os.Remove("test_fare_estimates.csv")
    }
}
