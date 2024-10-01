package main

import (
    "encoding/csv"
    "fmt"
    "math"
    "os"
    "strconv"
    "sync"
    "time"
)

type DeliveryPoint struct {
    ID        string
    Lat       float64
    Lng       float64
    Timestamp int64
}

const (
    movingRateDay   = 0.74
    movingRateNight = 1.30 
    idleRate        = 11.90 
    flagFare        = 1.30
    minFare         = 3.47
)


func haversine(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 
    dLat := (lat2 - lat1) * (math.Pi / 180)
    dLon := (lon2 - lon1) * (math.Pi / 180)
    lat1 = lat1 * (math.Pi / 180)
    lat2 = lat2 * (math.Pi / 180)

    a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}


func readfile(filePath string) map[string][]DeliveryPoint {
    file, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return nil
    }
    defer file.Close()

    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1

    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        return nil
    }

    deliveries := make(map[string][]DeliveryPoint)

    for i, record := range records {
        if i == 0 {
            continue
        }

        id := record[0]
        lat, _ := strconv.ParseFloat(record[1], 64)
        lng, _ := strconv.ParseFloat(record[2], 64)
        timestamp, _ := strconv.ParseInt(record[3], 10, 64)

        point := DeliveryPoint{
            ID:        id,
            Lat:       lat,
            Lng:       lng,
            Timestamp: timestamp,
        }

        deliveries[id] = append(deliveries[id], point)
    }

    return deliveries
}


func timeOfDayFare(timestamp int64) float64 {
    t := time.Unix(timestamp, 0).UTC()
    hour := t.Hour()
    if hour >= 5 && hour < 24 {
        return movingRateDay
    }
    return movingRateNight
}

func calculateFare(points []DeliveryPoint) float64 {
    totalFare := flagFare 
    totalIdleTime := 0.0

    for i := 1; i < len(points); i++ {
        p1 := points[i-1]
        p2 := points[i]

        distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
        timeDiff := float64(p2.Timestamp - p1.Timestamp) / 3600.0

        if timeDiff > 0 {
            speed := distance / timeDiff

            if speed > 10 { 
                fareRate := timeOfDayFare(p1.Timestamp)
                totalFare += distance * fareRate
            } else { 
                totalIdleTime += timeDiff
            }
        }
    }

   
    totalFare += totalIdleTime * idleRate

    
    if totalFare < minFare {
        totalFare = minFare
    }

    return totalFare
}


func filterPoints(pointsID map[string][]DeliveryPoint) map[string][]DeliveryPoint {
    validatemap := make(map[string][]DeliveryPoint)

    for _, v := range pointsID {
        if len(v) > 0 {
            validatemap[v[0].ID] = append(validatemap[v[0].ID], v[0])
            for i := 1; i < len(v); i++ {
                p1 := v[i-1]
                p2 := v[i]

                distance := haversine(p1.Lat, p1.Lng, p2.Lat, p2.Lng)
                timeDiff := float64(p2.Timestamp - p1.Timestamp) / 3600.0
                if timeDiff > 0 {
                    speed := distance / timeDiff
                    if speed <= 100 {
                        validatemap[v[i].ID] = append(validatemap[v[i].ID], p2)
                    }
                }
            }
        }
    }
    return validatemap
}


func writeToCSV(filename string, results map[string]float64) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    
    writer.Write([]string{"id_delivery", "fare_estimate"})

    
    for id, fare := range results {
        writer.Write([]string{id, fmt.Sprintf("%.2f", fare)})
    }

    return nil
}

func main() {
    filePath := "./sample_data.csv"
    deliveries := readfile(filePath)

    
    results := make(chan struct {
        ID   string
        Fare float64
    }, len(deliveries))

   
    var wg sync.WaitGroup

  
    fareEstimates := make(map[string]float64)

    
    for id, points := range deliveries {
        wg.Add(1)
        go func(id string, points []DeliveryPoint) {
            defer wg.Done()

            filterpoint := filterPoints(map[string][]DeliveryPoint{id: points})
            fare := calculateFare(filterpoint[id])

            results <- struct {
                ID   string
                Fare float64
            }{ID: id, Fare: fare}
        }(id, points)
    }

   
    go func() {
        wg.Wait()
        close(results)
    }()

   
    for result := range results {
        fareEstimates[result.ID] = result.Fare
    }

    
    err := writeToCSV("fare_estimates.csv", fareEstimates)
    if err != nil {
        fmt.Println("Error writing to CSV:", err)
    } else {
        fmt.Println("Fare estimates written to fare_estimates.csv")
    }
}
