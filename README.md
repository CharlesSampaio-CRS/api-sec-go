# api-sec-go

**api-sec-go** is a JWT authentication API built with Go, providing secure user authentication and authorization backed by MongoDB.

## Technologies Used

- **Go (Golang) 1.23+** — The programming language used to build the API.
- **MongoDB** — NoSQL database for storing user data and tokens.
- **JWT (github.com/golang-jwt/jwt/v5)** — JSON Web Tokens for stateless authentication.
- **Gin (github.com/gin-gonic/gin)** — HTTP web framework for routing and middleware support.
- **Swaggo (github.com/swaggo/gin-swagger)** — Auto-generated API documentation (Swagger UI).
- **dotenv (github.com/joho/godotenv)** — For environment variable management.
- **MongoDB official driver (go.mongodb.org/mongo-driver)** — MongoDB Go driver for database interactions.
- Additional libraries for validation, JSON processing, security, and compression.

## Features

- User registration and login with JWT token issuance.
- Token verification middleware for protected routes.
- Swagger documentation available at `/docs`.
- Environment configuration via `.env` file.

## .ENV 
### 
MONGODB_URI=mongodb://localhost:27017/db-exemplo

### 
JWT_SECRET=your_super_secret_jwt_key

### 
PORT=8080

## Prerequisites

- Go 1.23 or higher installed.
- MongoDB instance running and accessible.
- `git` installed to clone the repository.
## Get Started 
- Clone project https://github.com/CharlesSampaio-CRS/api-sec-go 
- cd api-sec-go 
- go run main.go
