# CollP Backend
```bash
# รัน กatabase ก่อน ด้วย docker
docker run -d --name collp-postgres \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=1234 \
  -e POSTGRES_DB=collp_backend \
  -p 5433:5432 postgres:15-alpine
# จากนั้น run project ขึ้นมาได้เลย database จะ migrate อัตโนมัติ
# ตาม folder models  database จะสร้างตารางตามนั้น
```
## 🚀 วิธี Run Project (Quick Start)
## รัน database ก่อน ด้วย docker
### ⚡ แบบง่ายที่สุด
```bash
# ติดตั้ง Air สำหรับ hot reload
go install github.com/cosmtrek/air@latest

# รัน project
air
```

### 🔧 วิธีอื่นๆ ในการรัน
```bash
# 1. รันแบบธรรมดา
go run cmd/server/main.go

# 2. Build แล้วรัน
go build -o app cmd/server/main.go
./app

# 3. ใช้ CompileDaemon (auto-rebuild)
go install github.com/githubnemo/CompileDaemon@latest
CompileDaemon -build="go build -o app cmd/server/main.go" -command="./app"

# 4. ใช้ Docker
docker-compose up --build
```

### 📋 เช็คก่อนรัน
```bash
# ตรวจสอบ Go version
go version

# ติดตั้ง dependencies
go mod tidy

# เช็คว่ามีไฟล์ .env และ rsa.pem แล้วหรือยัง
ls -la .env rsa.pem
```

---

A Go backend application for CollP platform built with Gin framework, GORM, and PostgreSQL.

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── config/
│   └── config.go            # Database and configuration setup
├── controllers/
│   ├── auth_controller.go   # Authentication controllers
│   ├── main_controller.go   # Main menu controllers
│   └── user_controller.go   # User management controllers
├── middleware/
│   └── logger.go            # JWT authentication middleware
├── models/
│   └── user.go              # Database models
├── repositories/
│   ├── main_repo.go         # Data access layer
│   └── user_repo.go         # User repository
├── routes/
│   └── router.go            # Route definitions
├── services/
│   ├── main_usecase.go      # Business logic layer
│   └── user_usecase.go      # User service logic
├── utils/
│   ├── functions.go         # Utility functions
│   ├── hash.go              # Password hashing utilities
│   └── jwt.go               # JWT token utilities
├── validators/
│   └── user_validator.go    # Input validation
├── .env.example             # Environment variables template
├── air.toml                 # Hot reload configuration
├── docker-compose.yml       # Docker compose setup
├── Dockerfile               # Docker configuration
└── rsa.pem                  # RSA private key for JWT
```

## Features

- **RESTful API** with Gin framework
- **PostgreSQL** database with GORM ORM
- **JWT Authentication** with RSA key signing
- **Google OAuth 2.0** integration
- **Password hashing** with bcrypt
- **Input validation** and sanitization
- **CORS support** and security middleware
- **Rate limiting** middleware
- **Hot reload** with Air for development
- **Docker support** with docker-compose

## Setup

### 💡 Setup แบบง่าย (สำหรับคนขี้เกียจ)
```bash
# 1. Clone project
git clone <your-repo-url>
cd CollP-Backend

# 2. ติดตั้ง dependencies
go mod tidy

# 3. Copy .env
cp .env.example .env
# แล้วแก้ไข .env ตามต้องการ

# 4. สร้าง database (ถ้ายังไม่มี)
docker run -d --name collp-postgres \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=1234 \
  -e POSTGRES_DB=collp_backend \
  -p 5433:5432 postgres:15-alpine

# 5. รัน project
air
```

### Prerequisites

- Go 1.24.5 or later
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Update the following variables:
- Database credentials (DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
- Google OAuth credentials (GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET)
- JWT configuration
- Server port and frontend URL

### Development Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up PostgreSQL database:**
   ```bash
   # Using Docker (ใช้ port 5433 ตาม .env ของคุณ)
   docker run -d \
     --name collp-postgres \
     -e POSTGRES_USER=admin \
     -e POSTGRES_PASSWORD=1234 \
     -e POSTGRES_DB=collp_backend \
     -p 5433:5432 \
     postgres:15-alpine
   ```

3. **Generate RSA key pair for JWT:**
   ```bash
   # Generate private key (ไฟล์นี้มีอยู่แล้วในโปรเจคของคุณ)
   openssl genrsa -out rsa.pem 2048
   
   # Generate public key (optional)
   openssl rsa -in rsa.pem -pubout -out rsa_public.pem
   ```

4. **Run the application:**
   ```bash
   # Development mode with hot reload (แนะนำ)
   air
   
   # หรือ build และ run แบบธรรมดา
   go build -o tmp/main cmd/server/main.go
   ./tmp/main
   
   # หรือ run ตรงๆ
   go run cmd/server/main.go
   ```

📌 **Server จะรันที่:** `http://localhost:8080`

### Docker Setup

1. **Using Docker Compose (แบบง่าย):**
   ```bash
   # รันทุกอย่างพร้อมกัน (database + backend)
   docker-compose up --build
   
   # รันแบบ background
   docker-compose up -d
   
   # หยุด services
   docker-compose down
   ```

🐳 **Docker จะรัน:** Backend + PostgreSQL พร้อมกัน

## API Endpoints

### Public Endpoints
- `GET /api/auth/google/login` - Initiate Google OAuth login
- `GET /api/auth/google/callback` - Google OAuth callback
- `POST /api/collp/login` - User login
- `POST /api/collp/register` - User registration

### Protected Endpoints (Requires JWT)
- `GET /api/collp/main-menu` - Get main menu items

## Database Models

### User Model
```go
type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null"`
    Name      string         `json:"name" gorm:"not null"`
    GoogleID  string         `json:"google_id" gorm:"uniqueIndex"`
    Avatar    string         `json:"avatar"`
    IsActive  bool           `json:"is_active" gorm:"default:true"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

## Security Features

- **JWT tokens** with RSA256 signing
- **Password hashing** with bcrypt
- **CORS protection** with configurable origins
- **Rate limiting** (100 requests per minute)
- **Security headers** (XSS protection, content type nosniff, etc.)
- **Input validation** and sanitization

## Development Tools

- **Air** for hot reloading during development
- **Docker** for containerized development and deployment
- **GORM** for database operations with auto-migration
- **Gin** framework for high-performance HTTP routing

## 📝 คำสั่งที่ใช้บ่อย

### การรัน Project
```bash
air                           # รันแบบ hot reload
go run cmd/server/main.go     # รันแบบธรรมดา
docker-compose up --build    # รันด้วย Docker
```

### การจัดการ Database  
```bash
# เริ่ม PostgreSQL
docker run -d --name collp-postgres \
  -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=1234 \
  -e POSTGRES_DB=collp_backend -p 5433:5432 postgres:15-alpine

# หยุด database
docker stop collp-postgres

# ลบ database
docker rm collp-postgres
```

### อื่นๆ
```bash
go mod tidy                   # ติดตั้ง dependencies
go test ./...                 # รัน tests
go build -o app cmd/server/main.go  # build executable
```

---

## License

MIT License for keep secret key.
    GOOGLE_CLIENT_ID=your client id
    GOOGLE_CLIENT_SECRET=your secret key
    GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
    GOOGLE_USERINFO=https://www.googleapis.com/oauth2/v2/userinfo
    GOOGLE_USERINFO_EMAIL=https://www.googleapis.com/auth/userinfo.email
    GOOGLE_USERINFO_PROFILE=https://www.googleapis.com/auth/userinfo.profile
3. Create file rsa.pem and in file and private key for system middleware.