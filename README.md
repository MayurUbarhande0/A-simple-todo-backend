# GoTodo — REST API in Go

A simple Todo REST API built with Go, MySQL, JWT authentication, and Gorilla Mux. Built as a first backend project to learn Go, REST API design, JWT auth, and clean folder structure.

---

## Tech Stack

- **Language:** Go
- **Router:** Gorilla Mux
- **Database:** MySQL
- **Auth:** JWT (golang-jwt/jwt)
- **Password Hashing:** bcrypt
- **Config:** godotenv

---

## Folder Structure

```
gotodo/
├── db/
│   └── mysql.go          # database connection
├── internal/
│   ├── auth/
│   │   ├── handler.go    # register & login handlers
│   │   ├── jwt.go        # generate & validate JWT
│   │   └── model.go      # User struct, RegisterInput, LoginInput
│   ├── todo/
│   │   ├── handler.go    # CRUD handlers
│   │   ├── repository.go # SQL queries
│   │   └── model.go      # Task struct
│   └── middleware/
│       └── auth.go       # JWT auth middleware
├── main.go               # entry point, route registration
├── .env                  # environment variables (never commit this)
├── go.mod
└── go.sum
```

---

## API Routes

### Auth (Public)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Login and receive JWT token |

### Todos (Protected — requires JWT)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/todos` | Get all todos for logged-in user |
| POST | `/todos` | Create a new todo |
| PUT | `/todos/{id}` | Update a todo (mark done/undone) |
| DELETE | `/todos/{id}` | Delete a todo |

---

## Getting Started

### Prerequisites

- Go 1.21+
- MySQL

### Setup

**1. Clone the repo**
```bash
git clone https://github.com/yourusername/gotodo.git
cd gotodo
```

**2. Create the database and tables**

Run this in MySQL Workbench or any MySQL client:

```sql
CREATE DATABASE gotodo;
USE gotodo;

CREATE TABLE users (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    email      VARCHAR(100) UNIQUE NOT NULL,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    id         INT AUTO_INCREMENT PRIMARY KEY,
    user_id    INT NOT NULL,
    title      VARCHAR(255) NOT NULL,
    done       BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
```

**3. Create `.env` file**
```env
DB_USER=root
DB_PASS=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=gotodo
JWT_SECRET=your_jwt_secret_here
PORT=8080
```

**4. Install dependencies**
```bash
go mod tidy
```

**5. Run the server**
```bash
go run main.go
```

Server starts at `http://localhost:8080`

---

## Testing with Postman

### Register
```
POST /auth/register
Content-Type: application/json

{
    "name": "Mayur",
    "email": "mayur@test.com",
    "password": "password123"
}
```

### Login
```
POST /auth/login
Content-Type: application/json

{
    "email": "mayur@test.com",
    "password": "password123"
}
```

Copy the token from the response and use it in the Authorization header for all todo requests:
```
Authorization: Bearer <your_token_here>
```

### Create Todo
```
POST /todos
Authorization: Bearer <token>
Content-Type: application/json

{
    "title": "Learn JWT"
}
```

### Get All Todos
```
GET /todos
Authorization: Bearer <token>
```

### Update Todo
```
PUT /todos/1
Authorization: Bearer <token>
Content-Type: application/json

{
    "done": true
}
```

### Delete Todo
```
DELETE /todos/1
Authorization: Bearer <token>
```

---

## Key Concepts Learned

- Handler/Repository separation — handlers talk HTTP, repositories talk SQL
- JWT authentication flow — generate on login, validate on protected routes
- Middleware — intercepts requests before they reach handlers
- Context — passing user_id from middleware to handlers via `r.Context()`
- bcrypt — never store plain text passwords
- Foreign key constraints — data integrity between users and tasks


## Author

Mayur Ubarhande