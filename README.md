

# Go + Gin + MongoDB + New Relic APM Example

This project is a REST API example written in Go, using the Gin framework and MongoDB, and is observable with New Relic APM.

## Features
- **/api/data**: Lists items from MongoDB (GET)
- **/api/add**: Adds a new item to MongoDB (POST, JSON: `{ "name": "value" }`)
- HTTP and DB operations are tracked with New Relic
- Automatically inserts a sample item if the collection is empty at startup

## Setup
1. **Go and MongoDB must be installed**
2. Get a license key from your New Relic account
3. Enter your license key in the `ConfigLicense` section in `main.go`
4. Install dependencies:
   ```sh
   go mod tidy
   ```
5. Start MongoDB:
   ```sh
   mongod
   ```
6. Start the application:
   ```sh
   go run main.go
   ```

## API Usage
- **GET /api/data**: Lists all items
- **POST /api/add**: Adds a new item
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{"name":"example"}' http://localhost:8080/api/add
  ```

## New Relic
- When the application starts, it automatically sends data to New Relic
- DB operations (find, insert) appear as Datastore in the New Relic APM UI

## Notes
- Do not share your license key in `main.go`!
- This is for demo purposes; additional security and configuration are required for production.
