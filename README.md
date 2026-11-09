# University of Kigali Invoices & Payments System

A comprehensive backend API for managing student invoices, payments, and financial records for the University of Kigali.

## Quick Start

This repository contains the backend API. To get started:

```bash
cd backend
./setup.sh
make run
```

The API will be available at `http://localhost:9107`

## Documentation

For detailed documentation including API endpoints, architecture, and setup instructions, see [backend/README.md](./backend/README.md).

## Project Structure

```
TEAM3/
└── backend/            # Go backend API
    ├── cmd/           # Application entry point
    ├── internal/      # Internal packages (config, middleware, routes)
    ├── services/      # Service-oriented architecture
    ├── shared/        # Shared components (repositories)
    ├── migrations/    # Database migrations
    ├── tools/         # Database utility scripts
    └── README.md      # Detailed documentation
```

## Technology Stack

- **Language**: Go (Golang)
- **Database**: MySQL 8.0+
- **Architecture**: Service-Oriented Architecture (SOA)
- **Authentication**: JWT

## License

This project is for University of Kigali internal use.
