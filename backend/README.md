# University of Kigali Invoices & Payments API

A comprehensive backend API for managing student invoices, payments, and financial records for the University of Kigali, based on the UNSTOPPABLE API specification.

## Project Structure

The project follows a **service-oriented architecture** where each service has its own controllers, models, and tests.

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── auth/                    # JWT authentication
│   ├── config/                  # Configuration management
│   ├── jobs/                    # Background job scheduler
│   ├── middleware/              # HTTP middleware (CORS, Auth)
│   └── routes/                  # Route definitions
├── services/                    # Service-oriented architecture
│   ├── auth/                    # Authentication service
│   │   ├── controllers/
│   │   └── models/
│   ├── student/                 # Student service
│   │   ├── controllers/
│   │   └── models/
│   ├── invoice/                 # Invoice service
│   │   ├── controllers/
│   │   └── models/
│   ├── payment/                 # Payment service
│   │   ├── controllers/
│   │   └── models/
│   ├── receipt/                 # Receipt service
│   │   ├── controllers/
│   │   └── models/
│   ├── fee/                     # Fee type service
│   │   ├── controllers/
│   │   └── models/
│   ├── finance/                 # Finance service
│   │   ├── controllers/
│   │   └── models/
│   ├── admin/                   # Admin service
│   │   ├── controllers/
│   │   └── models/
│   └── webhook/                 # Webhook service
│       └── controllers/
├── shared/                      # Shared components
│   ├── repositories/            # Database operations
│   └── utils/                   # Utility functions
├── migrations/                  # Database migrations
└── .env.example                 # Environment variables template
```

## Architecture Overview

### Service-Oriented Design

Each service is self-contained with:
- **Controllers**: Handle HTTP requests and responses
- **Models**: Define data structures specific to the service

### Shared Components

- **Repositories**: Database operations shared across services
- **Middleware**: Common HTTP middleware (CORS, Auth)
- **Routes**: Centralized route registration
- **Config**: Configuration management
- **Jobs**: Background job scheduler for payment status checks and overdue invoices

## Features

- **User Management**: Role-based authentication (student, finance, admin)
- **Student Management**: Student records and profile management
- **Fee Structure**: Configurable fee types and amounts
- **Invoice Management**: Create, update, and track invoices
- **Payment Processing**: MoPay integration with automatic status updates
- **Receipt Generation**: Automatic receipt generation for paid payments
- **Audit Trail**: Complete audit log for all operations
- **Background Jobs**: Payment status monitoring and overdue invoice tracking
- **API Documentation**: RESTful API following UNSTOPPABLE specification

## Prerequisites

- Go 1.27 or higher
- MySQL 8.0 or higher (XAMPP MySQL is supported)
- Git

## Quick Start (Single-Command Setup)

The easiest way to get started is using the provided setup script:

```bash
cd backend
./setup.sh
```

This script will:
- Check for Go installation
- Install Go dependencies
- Detect and configure XAMPP MySQL (if available)
- Create the database
- Run migrations
- Seed test data (users, programmes, semesters, students)
- Build the application

After setup completes, start the server:
```bash
make run
# or
./bin/server
```

The server will start on http://localhost:9107

### Seed Data

To populate the database with test data:
```bash
cd migrations
go run seed.go
```

This will create:
- 3 test users (student, finance, admin)
- 4 programmes
- 2 semesters
- 1 student record

### Test Credentials

**Student Account:**
- Email: `student@uk.ac.rw`
- Password: `password123`
- Role: student
- Student Number: UOK/2024/0001

**Finance Account:**
- Email: `finance@uk.ac.rw`
- Password: `password123`
- Role: finance

**Admin Account:**
- Email: `admin@uk.ac.rw`
- Password: `password123`
- Role: admin

**Login Example:**
The login endpoint accepts either email or student number as the identifier:
```bash
# Login with email
curl -X POST http://localhost:9107/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"student@uk.ac.rw","password":"password123"}'

# Login with student number
curl -X POST http://localhost:9107/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"UOK/2024/0001","password":"password123"}'
```

**Finance Account:**
- Email: `finance@uk.ac.rw`
- Password: `password123`
- Role: finance

**Admin Account:**
- Email: `admin@uk.ac.rw`
- Password: `password123`
- Role: admin

## Manual Installation

If you prefer manual setup or need to customize the process:

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd TEAM3/backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

   Required environment variables:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=university_invoices
   SERVER_PORT=8080
   JWT_SECRET=your-secret-key-change-in-production
   MOPAY_API_URL=https://api.mopay.ishema.rw
   MOPAY_API_TOKEN=your-mopay-token
   MOPAY_WALLET_PHONE=250790989830
   PAYMENT_CALLBACK_URL=
   PAYMENT_TIMEOUT_MS=25000
   ```

4. **Create database**
   ```bash
   # For standard MySQL:
   mysql -u root -p
   CREATE DATABASE university_invoices;
   EXIT;

   # For XAMPP MySQL:
   /opt/lampp/bin/mysql -u root -e "CREATE DATABASE university_invoices;"
   ```

5. **Run migrations**
   ```bash
   cd migrations
   go run migrate.go
   ```

## Makefile Commands

```bash
make setup      # Run complete setup (recommended)
make build      # Build the application
make run        # Run the application
make test       # Run tests
make clean      # Clean build artifacts
make migrate    # Run database migrations
make deps       # Download dependencies
make help       # Show all available commands
```

## Running the Application

### Development Mode
```bash
go run ./cmd/server/main.go
```

### Production Build
```bash
go build -o bin/server ./cmd/server/main.go
./bin/server
```

## API Endpoints

### Base Path
All endpoints are prefixed with `/api/v1`

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `POST /api/v1/auth/forgot-password` - Request password reset
- `GET /api/v1/auth/me` - Current user info

### Student
- `GET /api/v1/student/dashboard` - Student dashboard

### Invoices
- `GET /api/v1/invoices` - List invoices
- `POST /api/v1/invoices` - Create invoice
- `GET /api/v1/invoices/{invoiceId}` - Get invoice details
- `PUT /api/v1/invoices/{invoiceId}` - Update draft invoice
- `POST /api/v1/invoices/{invoiceId}/issue` - Issue invoice
- `POST /api/v1/invoices/{invoiceId}/payments` - Start MoPay payment

### Payments
- `GET /api/v1/payments` - List payments
- `GET /api/v1/payments/{paymentId}` - Get payment details
- `POST /api/v1/payments/{paymentId}/cancel` - Cancel pending payment
- `GET /api/v1/payments/{paymentId}/notes` - List payment notes (finance only)
- `POST /api/v1/payments/{paymentId}/notes` - Add payment note (finance only)

### Receipts
- `GET /api/v1/receipts/{receiptId}` - Get receipt
- `GET /api/v1/receipts/{receiptId}/pdf` - Download receipt PDF

### Finance
- `GET /api/v1/finance/dashboard` - Finance overview
- `GET /api/v1/finance/reports/summary` - Financial summary
- `GET /api/v1/finance/reports/accounts` - Account summary per student
- `GET /api/v1/finance/students` - Search students
- `GET /api/v1/finance/programmes` - List programmes
- `GET /api/v1/finance/semesters` - List semesters

### Fee Types
- `GET /api/v1/fee-types` - List fee types
- `POST /api/v1/fee-types` - Create fee type (finance only)
- `GET /api/v1/fee-types/{feeTypeId}` - Get fee type
- `PUT /api/v1/fee-types/{feeTypeId}` - Update fee type (finance only)

### Admin
- `GET /api/v1/admin/users` - List users (admin only)
- `POST /api/v1/admin/users` - Invite user (admin only)
- `PATCH /api/v1/admin/users/{userId}/access` - Change role or access (admin only)
- `GET /api/v1/admin/roles` - Role permissions (admin only)
- `GET /api/v1/admin/audit-logs` - Audit log (admin only)

### Webhooks
- `POST /api/v1/webhooks/mopay` - MoPay status notification

### Health Check
- `GET /health` - Health check endpoint

## Database Schema

### Tables
- `users` - User accounts and authentication
- `programmes` - Academic programmes
- `semesters` - Academic semesters
- `students` - Student information
- `fee_types` - Fee structure definitions
- `invoices` - Invoice records
- `invoice_items` - Invoice line items
- `payments` - Payment records
- `receipts` - Receipt records
- `payment_notes` - Payment follow-up notes
- `audit_logs` - Audit trail
- `schema_migrations` - Migration tracking

## User Roles

- **student** - Can view their own invoices and payments, make payments
- **finance** - Can manage invoices, payments, fee types, and reports
- **admin** - Full system access including user management and audit logs

## Payment Flow

1. Student initiates payment via `POST /api/v1/invoices/{invoiceId}/payments`
2. Payment is created with status "Pending" and MoPay request is initiated
3. Student approves payment on their phone via MoPay
4. MoPay webhook updates payment status
5. If payment is "Paid", invoice balance is updated and receipt is generated
6. Background jobs monitor pending payments and mark expired ones

## Security Notes

- Passwords are hashed using bcrypt
- JWT authentication with 24-hour expiry
- Role-based access control on all endpoints
- SQL injection prevention through parameterized queries
- CORS enabled for cross-origin requests
- Phone numbers are masked in API responses
- Audit logging for all critical operations

## Background Jobs

The application includes background jobs that run periodically:

1. **Payment Status Check** (every minute)
   - Checks for pending payments that have expired
   - Marks expired payments as "Expired"

2. **Overdue Invoice Check** (daily)
   - Checks for invoices past their due date
   - Marks unpaid invoices as "Overdue"

## Testing

Run all tests:
```bash
go test ./...
```

The project currently has no test files. Tests can be added for individual services as needed.

## Development

### Adding a New Service

1. Create service directory structure:
   ```bash
   mkdir -p services/newservice/controllers services/newservice/models
   ```

2. Add models in `services/newservice/models/`
3. Add controller in `services/newservice/controllers/`
4. Add repository in `shared/repositories/` (if needed)
5. Register routes in `internal/routes/routes.go`

### Code Organization

- **Services**: Feature-specific implementations
- **Shared**: Common utilities and database access
- **Internal**: Internal application components (config, middleware, routes, jobs)

## UNSTOPPABLE API Specification

This implementation follows the UNSTOPPABLE API specification for the University of Kigali Team 3 project. The specification includes:

- MySQL database schema with exact table structures
- RESTful API endpoints with proper role-based access
- MoPay payment gateway integration
- Automatic payment status handling via webhooks
- Receipt generation for paid payments
- Complete audit trail
- Background job scheduling

## License

This project is for University of Kigali internal use.

## Support

For support and questions, please contact the development team.
