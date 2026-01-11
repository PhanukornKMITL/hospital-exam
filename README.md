# Hospital Exam API

## 🧭 Overview
- Purpose: Middleware API for hospital staff to search and display patient information while enforcing strict same-hospital isolation.
- Architecture: Go (Gin) API behind Nginx, PostgreSQL for persistence, Docker Compose for local orchestration.
- Features: Staff account creation and JWT login; patient search by ID (national_id or passport_id) and filter-based search with optional fields; simple hospital management.
- Security: All patient operations require `Authorization: Bearer <token>`; JWT-based auth; results restricted to the authenticated staff’s hospital.
- Documentation: Interactive Swagger available at `/swagger/index.html`.
- Testing: Unit tests cover positive and negative scenarios for Hospital, Staff, and Patient services.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.25+ (for development)

### Setup & Run

1. **Clone the repository**
```bash
git clone https://github.com/PhanukornKMITL/hospital-exam.git
cd hospital-exam
git checkout dev
```

2. **Start all services**
```bash
docker-compose up -d
```

3. **Access the API**
- HTTP: `http://localhost` or `http://hospital-a.api.co.th` (Required add domain in next step)
- HTTPS: `https://localhost` or `https://hospital-a.api.co.th` (Required add domain in next step)
- Swagger UI: `http://localhost/swagger/index.html` or `http://hospital-a.api.co.th/swagger/index.html` (Required add domain in next step)

> 📝 **Note:** 
> - SSL certificates are automatically generated on first run for both domains
> - To use `hospital-a.api.co.th`, add it to your hosts file (see Optional section below)
> - No manual SSL setup required!

---

### 🌐 Optional: Use Custom Domain

If you want to access via custom domain (`hospital-a.api.co.th`) instead of localhost:

**For Linux/macOS:**
```bash
echo "127.0.0.1 hospital-a.api.co.th" | sudo tee -a /etc/hosts
```

**For Windows:**
1. Open PowerShell **as Administrator** (Right-click → Run as Administrator)
2. Run this command:
```powershell
Add-Content -Path C:\Windows\System32\drivers\etc\hosts -Value "127.0.0.1 hospital-a.api.co.th"
```

> **Important:** You MUST run PowerShell as Administrator, otherwise you'll get "Access Denied" error.

**Then access via:**
- HTTPS: `https://hospital-a.api.co.th`
- Swagger: `https://hospital-a.api.co.th/swagger/index.html`

> 📝 **Note:** The auto-generated SSL certificate already supports both localhost and hospital-a.api.co.th.

---

## 📋 API Endpoints

### Health Check
```bash
curl http://localhost/health
```

### Hospitals
- `GET /hospital` - Get all hospitals
- `POST /hospital` - Create hospital
- `PUT /hospital/:id` - Update hospital
- `DELETE /hospital/:id` - Delete hospital

### Staff
- `GET /staff` - Get all staff
- `POST /staff/create` - Create staff
- `POST /staff/login` - Staff login
- `PUT /staff/:id` - Update staff
- `DELETE /staff/:id` - Delete staff

### Patients (All require auth)
- `GET /patient` - Get all patients
- `POST /patient/create` - Create patient
- `PUT /patient/:id` - Update patient
- `DELETE /patient/:id` - Delete patient
- `GET /patient/search/:id` - Search patient by ID (national_id or passport_id)
- `POST /patient/search` - Search patients with filters

**Auth header example (required for all patient endpoints):**
```bash
curl -H "Authorization: Bearer <JWT_TOKEN>" http://localhost/patient
```

### 🧪 Testing
Run unit tests (services):
```bash
go test ./tests/unit/service/... -v
```

### Stop services
```bash
docker-compose down
```

### View logs
```bash
docker-compose logs -f
docker logs hospital-api
docker logs hospital-nginx
docker logs hospital-postgres
```

### Rebuild
```bash
docker-compose up -d --build
```

### 🔄 Clean Reset (Fresh Install)
To completely reset the project as if it's a fresh installation:

**For Linux/macOS:**
```bash
# Stop and remove all containers, volumes, and generated files
docker-compose down -v
rm -rf ssl/*
docker-compose up -d
```

**For Windows (PowerShell):**
```powershell
# Stop and remove all containers, volumes, and generated files
docker-compose down -v
Remove-Item -Path "ssl\*" -Force -ErrorAction SilentlyContinue
docker-compose up -d
```

> **What this does:**
> - Stops and removes all Docker containers
> - Deletes all database data (volumes)
> - Removes any generated SSL certificates
> - Restarts with fresh state (auto-generates dummy SSL)

## 📦 Services
- **API**: Go (Gin framework) - Internal port 8080
- **Nginx**: Reverse proxy - Ports 80 (HTTP) & 443 (HTTPS)
- **PostgreSQL**: Database - Port 5432
