# Inventory Management System

## Overview
A comprehensive web application designed to optimize the tracking, organization, and control of inventory products. The backend is built with Go and SQLite, while the frontend is developed using React and TypeScript. Features include automated backups, sales management with receipt generation, real-time dashboard analytics, and complete inventory control.

## Screenshot
<img width="1678" height="1069" alt="image" src="https://github.com/user-attachments/assets/d14c4362-a0f2-4625-a510-5450c73b1817" />
<img width="1004" height="532" alt="image" src="https://github.com/user-attachments/assets/d43f4d87-653c-46c2-8949-eacb8c1c4bc0" />
<img width="948" height="324" alt="image" src="https://github.com/user-attachments/assets/452c2cbf-00bc-4f1d-8cf2-5bd208c85886" />
<img width="937" height="619" alt="image" src="https://github.com/user-attachments/assets/7c051efc-b5af-4766-ad02-ad28569a96de" />
<img width="667" height="622" alt="sales" src="https://github.com/user-attachments/assets/09e4e9ee-5559-4210-82ef-96ff1f15c4c5" />


## Features

### Inventory Management
- Product management with CRUD operations
- SKU-based product search and tracking
- Category management and organization
- Low stock alerts with urgency-based sorting
- CSV export functionality for low stock reports
- Product filtering by category and status
- Batch and expiry date tracking
- Shelf location management

### Sales & Transactions
- Point of sale system with SKU scanning
- Sales transaction recording with multiple items
- Receipt PDF generation with customizable store details
- Date-based sales filtering and reporting
- Payment status tracking (paid, pending, partial, refunded)
- Customer and supplier management
- Real-time inventory updates on sales

### Dashboard & Analytics
- Real-time statistics (total products, low stock count, inventory value)
- Today's sales tracking
- Low stock product alerts
- Recent sales overview
- Expiring products monitoring

### User Management
- User authentication with JWT tokens
- Role-based access control (admin, manager, staff)
- User profile management
- Password change functionality
- Session management with secure cookies

### Settings & Configuration
- Store information management (name, address, phone)
- Store policy configuration
- Sales tax rate settings
- Receipt customization

### Backup System
- Automated database backups every 24 hours
- Backup on server startup and graceful shutdown
- 7-day backup retention with automatic cleanup
- Manual backup trigger via API
- Backup listing and management
- Safe backup process during active database usage

## Technologies

### Backend
- **Go**: High-performance, statically typed programming language
- **SQLite**: Lightweight, serverless SQL database engine
- **Echo**: Minimalist web framework for Go
- **github.com/golang-jwt/jwt/v5**: JSON Web Tokens for secure authentication
- **github.com/mattn/go-sqlite3**: Pure Go SQLite driver
- **golang.org/x/crypto/bcrypt**: Password hashing and verification
- **github.com/jung-kurt/gofpdf**: PDF generation for receipts

### Frontend
- **React**: Component-based UI library
- **TypeScript**: Type-safe JavaScript development
- **Tailwind CSS**: Utility-first CSS framework
- **Vite**: Fast build tool and development server
- **React Router**: Client-side routing
- **Heroicons**: Beautiful hand-crafted SVG icons

## Getting Started

### Prerequisites
- Go 1.23.0 or later
- Node.js 14.x or later
- SQLite 3.x (included with most systems)

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/tousif-izaz/Inventory-Management-System-Lightweight.git
   cd Inventory-Management-System-Lightweight
   ```

2. **Backend Setup:**
   ```sh
   cd server
   go mod download
   go build -o ../bin/imsapi ./cmd/imsapi
   ```

3. **Frontend Setup:**
   ```sh
   cd client
   npm install
   npm run build
   ```

4. **Environment Configuration:**
   Create a `.env` file in the server directory:
   ```env
   JWT_KEY=your-secret-key-here
   PORT=8080
   ```

### Running the Application

#### Development Mode

1. **Start the Backend:**
   ```sh
   cd server
   go run cmd/imsapi/main.go
   ```

2. **Start the Frontend:**
   ```sh
   cd client
   npm run dev
   ```

   The application will be available at `http://localhost:5173`

#### Production Mode

1. **Build both components** (see Installation steps above)

2. **Run the backend:**
   ```sh
   cd server
   go run cmd/imsapi/main.go
   ```

3. **Serve the frontend** using the built files in `client/dist`

#### Using Docker Compose

```sh
docker-compose up
```

The application will be available at `http://localhost:5173`

### Database Initialization

The SQLite database will be automatically created on first run at `server/data/inventory.db`. The schema includes:

- Users table with authentication
- Products with inventory tracking
- Categories for product organization
- Sales and transactions
- Settings for store configuration
- Automated backup system

### Default Credentials

Create a user through the signup page or use the API to create an admin user.

## API Documentation

Comprehensive API documentation is available in `server/API_SPECIFICATION.md`

### Key Endpoints

**Authentication:**
- `POST /login` - User login
- `POST /signup` - User registration
- `POST /logout` - User logout
- `GET /profile` - Get current user profile
- `PUT /profile/password` - Change password

**Products:**
- `GET /products` - List all products
- `GET /products/sku/:sku` - Get product by SKU
- `POST /products` - Create product
- `PUT /products/:id` - Update product
- `DELETE /products/:id` - Delete product

**Sales:**
- `GET /sales` - List all sales
- `POST /sales` - Create sale
- `GET /sales/:id/receipt` - Generate receipt PDF

**Backup:**
- `POST /backup` - Trigger manual backup
- `GET /backup/list` - List available backups

**Settings:**
- `GET /settings` - Get all settings
- `PUT /settings/:key` - Update setting

## Backup System

The application includes an automated backup system that:

- Creates a backup on server startup
- Runs scheduled backups every 24 hours
- Creates a final backup on graceful shutdown
- Keeps backups for 7 days (configurable)
- Stores backups in `server/data/backups/`

For detailed backup documentation, see `server/BACKUP_SYSTEM.md`

### Manual Backup

Trigger a manual backup via API:
```sh
curl -X POST http://localhost:8080/backup \
  -H "Cookie: token=your-jwt-token"
```

### Restoring from Backup

1. Stop the server
2. Copy the desired backup file:
   ```sh
   cp server/data/backups/inventory_backup_YYYY-MM-DD_HH-MM-SS.db server/data/inventory.db
   ```
3. Restart the server

## Project Structure

```
.
├── client/                 # Frontend React application
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── services/      # API services
│   │   ├── types/         # TypeScript types
│   │   └── App.tsx        # Main application component
│   └── package.json
│
├── server/                # Backend Go application
│   ├── cmd/imsapi/        # Application entry point
│   ├── pkg/
│   │   ├── backup/        # Backup system
│   │   ├── controller/    # HTTP handlers
│   │   ├── service/       # Business logic
│   │   ├── repository/    # Data access layer
│   │   ├── domain/        # Domain models
│   │   └── middleware/    # Authentication middleware
│   ├── data/              # SQLite database and backups
│   ├── schema_sqlite.sql  # Database schema
│   └── go.mod
│
└── docker-compose.yml     # Docker configuration
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Security Considerations

- JWT tokens are stored in HTTP-only cookies
- Passwords are hashed using bcrypt
- Authentication required for all API endpoints except login/signup
- Role-based access control for sensitive operations
- Database backups are stored locally (consider encrypting for production)

## Troubleshooting

### Database Locked Error
If you encounter database locked errors, ensure only one server instance is running.

### Backup Failures
Check write permissions for the `server/data/backups/` directory.

### Port Already in Use
Change the PORT in the `.env` file if 8080 is already in use.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
