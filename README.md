# University of Kigali Invoices & Payments System

A comprehensive system for managing student invoices, payments, and financial records for the University of Kigali.

## Project Overview

This system allows:
- Students to view invoices, track payments, and check outstanding balances
- Authorized university staff to manage invoices and record payments
- Finance office to generate financial reports
- Role-based access control for different user types

## Project Structure

```
TEAM3/
├── frontend/           # Frontend application (to be implemented)
└── backend/            # Go backend API with service-oriented architecture
    ├── cmd/           # Application entry point
    ├── internal/      # Internal packages (config, middleware, routes)
    ├── services/      # Service-oriented architecture
    │   ├── user/      # User service (controllers, models, tests)
    │   ├── student/   # Student service (controllers, models, tests)
    │   ├── invoice/   # Invoice service (controllers, models, tests)
    │   ├── payment/   # Payment service (controllers, models, tests)
    │   └── fee/       # Fee service (controllers, models, tests)
    ├── shared/        # Shared components (repositories)
    ├── migrations/    # Database migrations
    └── README.md      # Backend documentation
```

## Technology Stack

### Backend
- **Language**: Go (Golang)
- **Database**: MySQL 8.0+
- **Architecture**: Service-Oriented Architecture (SOA)
- **Driver**: go-sql-driver/mysql

### Frontend
- To be determined (React, Vue, or Angular)

## Getting Started

### Backend Setup

See [backend/README.md](./backend/README.md) for detailed backend setup instructions.

1. Navigate to backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   make install-deps
   ```

3. Configure environment:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. Create database and run migrations:
   ```bash
   mysql -u root -p -e "CREATE DATABASE university_invoices;"
   make migrate
   ```

5. Run the server:
   ```bash
   make run
   ```

The API will be available at `http://localhost:9000`

### Frontend Setup

Frontend implementation is pending. This folder is reserved for future frontend development.

## API Documentation

The backend provides RESTful API endpoints for:

- **Authentication**: Login and registration
- **Users**: User management
- **Students**: Student records
- **Fee Types**: Fee structure management
- **Invoices**: Invoice creation and management
- **Payments**: Payment processing and tracking

See [backend/README.md](./backend/README.md) for complete API endpoint documentation.

## Features

### Core Features
- ✅ User authentication and role-based access
- ✅ Student profile management
- ✅ Configurable fee structure
- ✅ Invoice generation and management
- ✅ Payment processing (cash, bank transfer, mobile money, card, check)
- ✅ Automatic receipt generation
- ✅ Audit trail for all operations
- ✅ Outstanding balance tracking
- ✅ Service-oriented architecture with isolated controllers, models, and tests

### Planned Features
- ⏳ JWT authentication
- ⏳ Financial reports
- ⏳ Email notifications
- ⏳ Payment gateway integration
- ⏳ Frontend dashboard
- ⏳ Export to PDF/Excel

## Database Schema

The system uses MySQL with the following main tables:
- `users` - User accounts
- `students` - Student information
- `fee_types` - Fee definitions
- `invoices` - Invoice records
- `invoice_items` - Invoice line items
- `payments` - Payment records
- `receipts` - Receipt records
- `audit_logs` - Audit trail

## User Roles

- **student** - View own invoices and payments
- **staff** - Manage invoices and payments
- **admin** - Full system access
- **finance** - Financial operations
- **registrar** - Academic records

## Development

### Prerequisites
- Go 1.21+
- MySQL 8.0+
- Git

### Backend Commands
```bash
cd backend

make help          # Show available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make clean         # Clean build artifacts
make migrate       # Run database migrations
make deps          # Download dependencies
```

### Service Testing
```bash
cd backend

# Test specific service
go test ./services/user/controllers/...
go test ./services/student/controllers/...
go test ./services/invoice/controllers/...
go test ./services/payment/controllers/...
go test ./services/fee/controllers/...

# Test all services
go test ./services/.../...
```

## Architecture

### Service-Oriented Design

The backend follows a service-oriented architecture where each service is self-contained:

- **User Service**: Handles user authentication and management
- **Student Service**: Manages student records
- **Invoice Service**: Handles invoice creation and management
- **Payment Service**: Processes payments and generates receipts
- **Fee Service**: Manages fee structure

Each service has:
- Controllers for HTTP handling
- Models for data structures
- Tests for endpoint validation

### Shared Components

- **Repositories**: Database operations shared across services
- **Middleware**: CORS and authentication middleware
- **Routes**: Centralized route registration
- **Config**: Configuration management

## Security Considerations

- Passwords are hashed using bcrypt
- SQL injection prevention through parameterized queries
- Role-based access control
- Audit trail for all financial operations
- CORS enabled for cross-origin requests

## License

This project is for University of Kigali internal use.

## Contact

For questions or support, please contact the development team.
# TEAM3
