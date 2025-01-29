# Receipt Processor Challenge

A REST API to process receipts and calculate reward points based on predefined rules.  
**Built with**: Go, Gin, Docker

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
  - [Run with Go](#1-run-with-go)
  - [Run with Docker](#2-run-with-docker)
- [API Documentation](#api-documentation)
- [Testing the Service](#testing-the-service)
- [Rules Implementation](#rules-implementation)
- [Important Notes](#important-notes)

---

## Prerequisites

- **Go 1.19+** ([Installation Guide](https://go.dev/doc/install))
- **Docker** (optional, [Installation Guide](https://docs.docker.com/get-docker/))

---

## Quick Start

### 1. Run with Go

```bash
# Clone the repository
git clone https://github.com/<your-username>/reciept-processor.git
cd Receipt-processor

# Install dependencies
go mod tidy

# Start the server (port 8080)
go run main.go

2. Run with Docker
bash
Copy
# Build the image
docker build -t Receipt-processor .

# Run the container (maps port 8080)
docker run -p 8080:8080 Receipt-processor
API Documentation
Endpoint 1: Process Receipt
http
Copy
POST /receipts/process
Request:
Submit a receipt JSON (see /examples directory).

Response:
200 OK with generated receipt ID:

json
Copy
{ "id": "adb6b560-0eef-42bc-9d16-df48f30e89b2" }
Endpoint 2: Get Points
http
Copy
GET /receipts/{id}/points
Response:
200 OK with points:

json
Copy
{ "points": 109 }
Testing the Service
Using CURL
Process a sample receipt:

bash
Copy
curl -X POST -H "Content-Type: application/json" \
-d @examples/simple-receipt.json \
http://localhost:8080/receipts/process
Get points (replace <id> with returned ID):

bash
Copy
curl http://localhost:8080/receipts/<id>/points
Example Outputs
Example File	Expected Points
simple-receipt.json	28
morning-receipt.json	15
Rules Implementation
Points are calculated using these rules:

1 point per alphanumeric character in retailer name

50 points if total is round dollar amount

25 points if total is multiple of 0.25

5 points per 2 items

Bonus for item descriptions divisible by 3

Special LLM Rule: +5 points if total > $10.00 (applies to LLM-generated code)

6 points for odd purchase day

10 points for purchases between 2-4 PM

Important Notes
🚨 Data Storage: Uses in-memory storage (data lost on restart)

🔌 Port: Service runs on http://localhost:8080

⚠️ Validation: Strict input validation for dates/times/patterns

🤖 LLM Rule: Explicitly implemented as per requirements

Copy

---

### Key Features:
- Clear **step-by-step instructions** for both Go and Docker
- Ready-to-use **CURL examples** with sample files
- **Rule implementation** transparency
- **Visual hierarchy** with emoji markers for important notes
- Table format for example outputs

This format ensures engineers can quickly:
1. Run the service locally or via Docker
2. Test endpoints with pre-configured examples
3. Verify rule implementation details
4. Understand architectural constraints
New chat
```
