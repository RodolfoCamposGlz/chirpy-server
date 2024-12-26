# Chirpy API

A simple social media API built with Go that allows users to post short messages called "chirps". The stored data is persisted in a PostgreSQL database.

## Features

- User authentication with JWT
- Create, read, and delete chirps
- User registration and login
- Sort chirps by creation date
- Filter chirps by author
- Chirpy Red premium user status

## API Endpoints

### Users

- `POST /api/users` - Register a new user
  - Body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
- `POST /api/login` - Login and receive JWT token
  - Body:
    ```json
    {
      "email": "user@example.com",
      "password": "password123"
    }
    ```
- `PUT /api/users` - Update user details
  - Body:
    ```json
    {
      "email": "newemail@example.com",
      "password": "newpassword123"
    }
    ```

### Chirps

- `POST /api/chirps` - Create a new chirp
  - Body:
    ```json
    {
      "body": "Hello, world!"
    }
    ```
- `GET /api/chirps` - Get all chirps
  - Query params:
    - `author_id` - Filter by author
    - `sort` - Sort order ("asc" or "desc")
- `GET /api/chirps/{chirpID}` - Get a specific chirp
- `DELETE /api/chirps/{chirpID}` - Delete a chirp (auth required)

### Premium

- `POST /api/polka/webhook` - Handle Polka webhook

## Setup

1. Clone the repository
2. Create a `.env` file with:
   ```
   PLATFORM=dev
   POLKA_KEY=your_polka_key
   JWT_SECRET=your_jwt_secret (Got from openssl rand -base64 64)
   DB_URL=postgresql://user:password@localhost:5432/dbname
   ```
3. Install dependencies:
   ```
   go mod download
   ```
4. Run database migrations:
   ```
   goose postgres "your_db_connection_string" up
   ```
5. Start the server:
   ```
   go run main.go
   ```

## Authentication

Most endpoints require JWT authentication. Include the token in requests:

```json
{
  "Authorization": "Bearer <token>"
}
```
