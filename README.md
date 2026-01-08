# Hospital Exam API

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

2. **Generate SSL certificate for local development**
```bash
./generate-ssl.sh
```

3. **Start all services**
```bash
docker-compose up -d
```

4. **Access the API**
- HTTPS: `https://localhost`
- HTTP: `http://localhost` (redirects to HTTPS)

> **Note:** Browser will show SSL warning because we use self-signed certificate. Click "Proceed" or "Accept Risk" to continue.

### Optional: Use Custom Domain

If you want to use `https://hospital-a.api.co.th` instead of `localhost`:

1. Edit `/etc/hosts` file:
```bash
sudo nano /etc/hosts
```

2. Add this line:
```
127.0.0.1 hospital-a.api.co.th
```

3. Access via: `https://hospital-a.api.co.th`

## 📋 API Endpoints

### Health Check
```bash
curl -k https://localhost/health
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

## 🛠 Development

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
