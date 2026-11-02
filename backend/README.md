# University of Kigali Invoices & Payments API

A comprehensive backend API for managing student invoices, payments, and financial records for the University of Kigali.

## Project Structure

The project follows a **service-oriented architecture** where each service has its own controllers, models, and tests.

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   │   └── config.go
│   ├── middleware/              # HTTP middleware
│   │   ├── cors.go
│   │   └── auth.go
│   └── routes/                  # Route definitions
│       └── routes.go
├── services/                    # Service-oriented architecture
│   ├── user/                    # User service
│   │   ├── controllers/
│   │   │   ├── user_controller.go
│   │   │   └── user_controller_test.go
│   │   └── models/
│   │       ├── user.go
│   │       └── audit.go
│   ├── student/                 # Student service
│   │   ├── controllers/
│   │   │   ├── student_controller.go
│   │   │   └── student_controller_test.go
│   │   └── models/
│   │       └── student.go
│   ├── invoice/                 # Invoice service
│   │   ├── controllers/
│   │   │   ├── invoice_controller.go
│   │   │   └── invoice_controller_test.go
│   │   └── models/
│   │       └── invoice.go
│   ├── payment/                 # Payment service
│   │   ├── controllers/
│   │   │   ├── payment_controller.go
│   │   │   └── payment_controller_test.go
│   │   └── models/
│   │       ├── payment.go
│   │       └── receipt.go
│   └── fee/                     # Fee service
│       ├── controllers/
│       │   ├── fee_controller.go
│       │   └── fee_controller_test.go
│       └── models/
│           └── fee.go
├── shared/                      # Shared components
│   ├── repositories/            # Database operations (shared across services)
│   │   ├── user_repository.go
│   │   ├── student_repository.go
│   │   ├── invoice_repository.go
│   │   ├── invoice_item_repository.go
│   │   ├── payment_repository.go
│   │   ├── fee_repository.go
│   │   ├── receipt_repository.go
│   │   └── audit_repository.go
├── migrations/                  # Database migrations
│   ├── 001_create_users_table.sql
│   ├── 002_create_students_table.sql
│   ├── 003_create_fee_types_table.sql
│   ├── 004_create_invoices_table.sql
│   ├── 005_create_invoice_items_table.sql
│   ├── 006_create_payments_table.sql
│   ├── 007_create_receipts_table.sql
│   ├── 008_create_audit_logs_table.sql
│   └── migrate.go
├── .env.example                 # Environment variables template
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Architecture Overview

### Service-Oriented Design

Each service is self-contained with:
- **Controllers**: Handle HTTP requests and responses
- **Models**: Define data structures specific to the service
- **Tests**: Unit tests for each controller endpoint

### Shared Components

- **Repositories**: Database operations shared across services
- **Middleware**: Common HTTP middleware (CORS, Auth)
- **Routes**: Centralized route registration
- **Config**: Configuration management

## Features

- **User Management**: Role-based authentication (student, staff, admin, finance, registrar)
- **Student Management**: Student records and profile management
- **Fee Structure**: Configurable fee types and amounts
- **Invoice Management**: Create, update, and track invoices
- **Payment Processing**: Record and track payments with multiple payment methods
- **Receipt Generation**: Automatic receipt generation for payments
- **Audit Trail**: Complete audit log for all operations
- **API Documentation**: RESTful API design

## Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Git

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd TEAM3/backend
   ```

2. **Install dependencies**
   ```bash
   make install-deps
   # or manually:
   go get github.com/go-sql-driver/mysql
   go get golang.org/x/crypto/bcrypt
   go mod tidy
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
   SERVER_PORT=9000
   ```

4. **Create database**
   ```bash
   mysql -u root -p
   CREATE DATABASE university_invoices;
   EXIT;
   ```

5. **Run migrations**
   ```bash
   make migrate
   # or manually:
   cd migrations
   go run migrate.go
   ```

## Running the Application

### Development Mode
```bash
make run
# or
go run ./cmd/server/main.go
```

### Production Build
```bash
make build
./bin/server
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

### Users
- `GET /api/users` - List all users
- `GET /api/users?id={id}` - Get user by ID
- `POST /api/users` - Create new user
- `PUT /api/users?id={id}` - Update user
- `DELETE /api/users?id={id}` - Delete user

### Students
- `GET /api/students` - List all students
- `GET /api/students?id={id}` - Get student by ID
- `GET /api/students/by-id?student_id={student_id}` - Get student by student ID
- `POST /api/students` - Create new student
- `PUT /api/students?id={id}` - Update student
- `DELETE /api/students?id={id}` - Delete student

### Fee Types
- `GET /api/fees` - List all fee types
- `GET /api/fees?id={id}` - Get fee type by ID
- `POST /api/fees` - Create new fee type
- `PUT /api/fees?id={id}` - Update fee type
- `DELETE /api/fees?id={id}` - Delete fee type

### Invoices
- `GET /api/invoices` - List all invoices
- `GET /api/invoices?id={id}` - Get invoice by ID
- `GET /api/invoices/student?student_id={id}` - Get invoices by student
- `POST /api/invoices` - Create new invoice
- `PUT /api/invoices?id={id}` - Update invoice
- `DELETE /api/invoices?id={id}` - Delete invoice

### Payments
- `GET /api/payments` - List all payments
- `GET /api/payments?id={id}` - Get payment by ID
- `GET /api/payments/invoice?invoice_id={id}` - Get payments by invoice
- `GET /api/payments/student?student_id={id}` - Get payments by student
- `POST /api/payments` - Create new payment
- `PUT /api/payments?id={id}` - Update payment
- `DELETE /api/payments?id={id}` - Delete payment

## Database Schema

### Tables
- `users` - User accounts and authentication
- `students` - Student information
- `fee_types` - Fee structure definitions
- `invoices` - Invoice records
- `invoice_items` - Invoice line items
- `payments` - Payment records
- `receipts` - Receipt records
- `audit_logs` - Audit trail
- `schema_migrations` - Migration tracking

## User Roles

- **student** - Can view their own invoices and payments
- **staff** - Can view and manage invoices and payments
- **admin** - Full system access
- **finance** - Financial operations and reporting
- **registrar** - Academic records access

## Testing

Run tests for a specific service:
```bash
# Test user service
go test ./services/user/controllers/...

# Test student service
go test ./services/student/controllers/...

# Test all services
go test ./services/.../...
```

## Development

### Adding a New Service

1. Create service directory structure:
   ```bash
   mkdir -p services/newservice/controllers services/newservice/models
   ```

2. Add models in `services/newservice/models/`
3. Add controller in `services/newservice/controllers/`
4. Add test file in `services/newservice/controllers/`
5. Add repository in `shared/repositories/` (if needed)
6. Register routes in `internal/routes/routes.go`

### Code Organization

- **Services**: Feature-specific implementations
- **Shared**: Common utilities and database access
- **Internal**: Internal application components (config, middleware, routes)

## Security Notes

- Passwords are hashed using bcrypt
- JWT authentication (to be implemented)
- Role-based access control (to be implemented)
- SQL injection prevention through parameterized queries
- CORS enabled for cross-origin requests

## TODO

- [ ] Implement JWT authentication
- [ ] Implement role-based access control middleware
- [ ] Add input validation
- [ ] Add API rate limiting
- [ ] Add comprehensive error handling
- [ ] Implement mock repositories for tests
- [ ] Add integration tests
- [ ] Add API documentation (Swagger/OpenAPI)
- [ ] Add financial reports endpoint
- [ ] Add email notifications
- [ ] Add payment gateway integration

## License

This project is for University of Kigali internal use.

## Support

For support and questions, please contact the development team.
