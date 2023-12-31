# JSON Server in Golang

This simple Golang server accepts HTTP POST requests with JSON data on port 8080. It validates the presence of a predefined "message" field in the JSON data. If the "message" field is present, it outputs the message to the server console and responds with a JSON success message. If the "message" field is missing or empty, it returns an HTTP error code 400 with a "Invalid JSON message" response.

## Usage

1. Install Golang: [https://golang.org/dl/](https://golang.org/dl/)
2. Run the server: `go run main.go`
3. Use Postman or any HTTP client to send POST requests with JSON data to `http://localhost:8080/`.

Example JSON data:
```json
{
  "message": "Hello, server! This is JSON data from Postman."
}
```

The server will respond with:
```json
{
  "status": "success",
  "message": "Data successfully received"
}
```

If the "message" field is missing, empty or written incorrectly, the server will return:
```json
{
  "status": "400",
  "message": "Invalid JSON message"
}
```
