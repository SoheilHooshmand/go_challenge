# Fare Calculation System

## 1. Objective
The goal of the system is to estimate the fare for deliveries based on geographical points (latitude, longitude, timestamp) using the Haversine formula to compute distances between points. The fare is calculated by considering both the distance traveled and the time spent idle (at speeds below 10 km/h). The system reads delivery data from a CSV file, processes the data concurrently for multiple deliveries, and writes the calculated fares back to a CSV file.

## 2. Key Components

### DeliveryPoint Structure
This struct represents a point in the delivery journey, containing the delivery ID, latitude, longitude, and a timestamp.

### Fare Rates
Constants like `movingRateDay`, `movingRateNight`, `idleRate`, `flagFare`, and `minFare` represent rates used for fare calculation. These vary based on the time of day and whether the vehicle is moving or idling.

### Haversine Formula
The `haversine()` function calculates the great-circle distance between two geographical points. This is essential for determining the distance traveled between consecutive points in the delivery route.

### Reading CSV Files
The `readfile()` function reads the delivery points from a CSV file and organizes them into a map where the delivery ID is the key, and the value is a slice of `DeliveryPoint` structs.

### Fare Calculation
The `calculateFare()` function computes the fare by iterating over each delivery's points, calculating the distance traveled, and adjusting the fare based on the speed of travel (idle or moving).

### Point Filtering
The `filterPoints()` function ensures that only valid points (with reasonable speeds, i.e., ≤ 100 km/h) are considered for fare calculation.

### Concurrency
Using Go's `sync.WaitGroup` and goroutines, the fare calculation for each delivery is performed concurrently to improve performance, especially with large datasets.

### Writing to CSV
The `writeToCSV()` function writes the calculated fare estimates for each delivery back into a new CSV file.

## 3. Design Decisions

### Concurrency
By utilizing goroutines for each delivery, the system processes multiple deliveries in parallel, which improves efficiency, especially when dealing with large datasets.

### Filtering
A filtering mechanism ensures that data points with unrealistic speeds are excluded, improving the accuracy of the fare calculation.

### Modularity
Each function handles a distinct task, promoting clean code and separation of concerns.

