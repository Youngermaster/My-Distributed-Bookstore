# Admin Service

## Overview

Provides admin dashboard backend, analytics, reporting, and system-wide management capabilities.

## Technology Stack

- **Language**: Go 1.21+
- **Framework**: Fiber v2
- **Database**: Aggregates data from all services
- **Ports**: HTTP: 8090, gRPC: 50060

## Responsibilities

- Admin dashboard backend
- Analytics and reporting
- User management (admin functions)
- Book management (admin functions)
- Order management (admin view)
- System health monitoring
- Bulk operations

## Dashboard Metrics

- Total orders (today/week/month)
- Revenue (today/week/month)
- Active users
- Low stock items
- Average order value
- Top selling books
- Order status breakdown

## REST API Endpoints

```
GET    /api/v1/admin/dashboard           # Dashboard stats
GET    /api/v1/admin/orders               # All orders
GET    /api/v1/admin/users                # All users
GET    /api/v1/admin/analytics/sales      # Sales analytics
GET    /api/v1/admin/analytics/inventory  # Inventory reports
POST   /api/v1/admin/books                # Bulk book import
PUT    /api/v1/admin/orders/:id/status    # Update order status
POST   /api/v1/admin/users/:id/roles      # Manage user roles
```

## Features

- Real-time dashboard
- Data aggregation from all services
- Export capabilities (CSV, PDF)
- Advanced filtering and search
- Role-based admin access
- Audit logging

## Next Steps

- [ ] Implement analytics aggregation
- [ ] Create dashboard endpoints
- [ ] Add bulk operations
- [ ] Integrate with all services via gRPC
- [ ] Add export functionality
- [ ] Write tests
