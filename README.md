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
- Go 1.21+ (for development)

### Setup & Run

1. **Clone the repository**
```bash
git clone <repository-url>
cd hospital-exam
```

2. **Start all services**
```bash
docker-compose up -d
```

3. **Access the API**
- HTTP (localhost): `http://localhost`
- Swagger UI: `http://localhost/swagger/index.html`

**That's it!** ✅ No additional setup required for localhost testing.

---

### 🔒 Optional: HTTPS with Custom Domain

If you want to use HTTPS with a custom domain (`https://hospital-a.api.co.th`):

1. **Generate SSL certificate**
```bash
./generate-ssl.sh
```

2. **Add domain to hosts file**
```bash
echo "127.0.0.1 hospital-a.api.co.th" | sudo tee -a /etc/hosts
```

3. **Restart services**
```bash
docker-compose restart
```

4. **Access via HTTPS**
- HTTPS: `https://hospital-a.api.co.th`
- Swagger: `https://hospital-a.api.co.th/swagger/index.html`

> **Note:** Browser will show SSL warning for self-signed certificate. Click "Proceed" to continue.

---

## 📋 API Endpoints

### Health Check
```bash
curl http://localhost/health
```

### Hospitals
- `GET /hospital` - Get all hospitals
- `POST /hospital` - Create hospital

### Staff
- `GET /staff` - Get all staff
- `POST /staff/create` - Create staff
- `POST /staff/login` - Staff login
- `DELETE /staff/:id` - Delete staff

### Patients
- `GET /patient` - Get all patients (requires auth)
- `POST /patient/create` - Create patient (requires auth)
- `GET /patient/search/:id` - Search patient by ID (requires auth)
- `POST /patient/search` - Search patients (requires auth)

Auth header example (required for patient endpoints):
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

## 📦 Services
- **API**: Go (Gin framework) - Internal port 8080
- **Nginx**: Reverse proxy - Ports 80 (HTTP) & 443 (HTTPS)
- **PostgreSQL**: Database - Port 5432
