# CineVault

CineVault is a full-stack movie booking application built with Go and vanilla JavaScript that provides a seamless experience for users to browse movies and book tickets, while giving administrators tools to manage content.

## Features

- **User Authentication** - Secure login and registration with role-based access control
- **Movie Browsing** - Explore current and upcoming movie listings with details
- **Ticket Booking** - Simple, intuitive booking flow with seat selection
- **Admin Dashboard** - Add and manage movies, view booking analytics
- **JWT Authentication** - Secure API access with refresh token support
- **Database Migrations** - Versioned schema changes with Goose
- **Type-Safe SQL Queries** - Auto-generated Go code using SQLC
- **AWS S3 Integration** - Store and retrieve media assets securely
- **Responsive Design** - Optimized for both desktop and mobile experiences

## Tech Stack

### Backend
- **Go** - Core backend language
- **PostgreSQL** - Persistent data storage
- **Goose** - Database migrations
- **SQLC** - Type-safe SQL queries with Go code generation
- **JWT** - Authentication and authorization
- **AWS S3** - Cloud storage for media files

### Frontend
- **Vanilla JavaScript** - No frameworks for lightweight performance
- **HTML5/CSS3** - Modern, responsive UI
- **Fetch API** - Asynchronous API communication

## Project Structure

```
├── frontend/             # Client-side code
│   ├── booking/          # Ticket booking interface
│   ├── createMovie/      # Admin movie creation interface
│   ├── movies/           # Movie browsing interface
├── internal/             # Internal Go packages
│   ├── auth/             # Authentication logic
│   ├── database/         # Database models and queries
│   ├── storage/          # AWS S3 integration logic
├── sql/                  # Database definitions
│   ├── queries/          # SQLC query definitions
│   ├── schema/           # Database migration scripts (Goose)
├── migrations/           # Goose migration files
├── cmd/                  # Application entry points
├── main.go               # Main application server
```

## API Endpoints

| Endpoint      | Method | Description |
|--------------|--------|-------------|
| `/login`     | POST   | User authentication |
| `/refresh`   | POST   | Refresh access token |
| `/users`     | POST   | Create new user account |
| `/movies`    | GET    | List all movies |
| `/movies`    | POST   | Create new movie (admin only) |
| `/bookings`  | POST   | Create new booking |
| `/upload`    | POST   | Upload media files to AWS S3 |

## Getting Started

### Prerequisites
- Go 1.16+
- PostgreSQL 12+
- Goose (Database migrations)
- SQLC (Go code generation from SQL queries)
- AWS S3 Bucket for media storage

### Installation

1. **Clone the repository**
   ```sh
   git clone https://github.com/vannestbun/movie-booking-platform.git
   cd movie-booking-platform
   ```

2. **Set up the database**
   ```sh
   psql -U postgres -c "CREATE DATABASE cinevault"
   ```

3. **Run migrations using Goose**
   ```sh
   goose -dir sql/schema postgres "$(DB_URL)" up
   goose -dir sql/schema postgres "$(DB_URL)" down
   ```

4. **Generate Go code from SQL queries using SQLC**
   ```sh
   sqlc generate
   ```

5. **Configure AWS S3 credentials** (set environment variables)
   ```sh
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   export AWS_REGION=your_region
   export S3_BUCKET_NAME=your_bucket_name
   ```

6. **Start the server (from the root of the project)**
   ```sh
   go run .
   ```

7. **Access the application at** `http://localhost:8080`

## Contributing

We welcome contributions to CineVault! To contribute, follow these steps:

1. **Fork the repository** on GitHub.
2. **Clone your fork** to your local machine:
   ```sh
   git clone https://github.com/vannestbun/movie-booking-platform.git
   ```
3. **Create a new branch** for your feature or bug fix:
   ```sh
   git checkout -b feature-name
   ```
4. **Make your changes** and commit them with a meaningful message:
   ```sh
   git commit -m "Added new feature X"
   ```
5. **Push your changes** to your fork:
   ```sh
   git push origin feature-name
   ```
6. **Submit a pull request** from your fork to the main repository.

Before submitting a pull request, ensure:
- Your code follows the project's style guidelines.
- You have tested your changes locally.
- Your commit messages are descriptive.

## Future Enhancements

- Payment gateway integration
- Email notifications for bookings
- Enhanced seat selection visualization
- Movie recommendations based on booking history
- Performance optimization and caching

## License

This project is licensed under the MIT License - see the LICENSE file for details.

