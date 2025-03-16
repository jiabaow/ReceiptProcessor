# ReceiptProcessor

## Getting Started
Follow these steps to run the application using Docker.

### Clone the Repo
After clone the repo go to /ReceiptProcessor.

### Build the Docker Image
```docker build -t receipt-processor .```

### Run the Docker Contaner
```docker run -p 80880:8080 receipt-processor```
The server will start on port 8080.

## Testing the Application
You can use the provided argument in test.txt to test the API endpoints using curl.

## API Endpoints
- POST /receipts/process. Process a new receipt.
- GET /receipts/{id}/points. Retrieve points for a specific receipt.

## Notes
- Receipts are stored in memory. Restarting the server will clear all receipts.
- Ensure the purchase date is in YYYY-MM-DD format.
- Purchase time should be in HH:MM (24-hour) format.