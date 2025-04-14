# Square POS API Integration – Qlub Backend Developer Case Study

This project is a backend API that integrates with Square's POS system to manage restaurant orders, payments, and customer interactions. It demonstrates the use of Go (Golang), PostgreSQL, and secure RESTful API design patterns.

---

## 📦 Tech Stack

- **Language**: Go 1.21+
- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT
- **3rd Party API**: Square POS (Sandbox environment)



## ⚙️ Setup Instructions

### 1. Clone the repository

`git clone https://github.com/OKWijewardena/square-pos-api.git `

`cd square-pos-api`\


### 2. Configure `.env` file

Create a `.env` file in the root with your credentials:

`DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=yourpassword DB_NAME=square_pos SQUARE_ACCESS_TOKEN=your_square_sandbox_token SQUARE_LOCATION_ID=your_sandbox_location_id`



### 3. Run the project

`go run main.go`

The server will start at:\
📡 `http://localhost:8080`



## 🧪 API Endpoints

### 🔓 Public Routes

| Method                        | Endpoint | Description                         |
| ----------------------------- | -------- | ----------------------------------- |
| POST                          | `/login` | Mock restaurant login (returns JWT) |
| GET                           | `/ping`  | Health check route                  |

---

### 🔐 Protected Routes (JWT Required)

| Method                        | Endpoint                          | Description                 |
| ----------------------------- | --------------------------------- | --------------------------- |
| POST                          | `/api/orders`                     | Create a restaurant order   |
| GET                           | `/api/orders/:order_id`           | Get order by ID             |
| GET                           | `/api/orders/table/:table_number` | Get orders for a table      |
| POST                          | `/api/payments`                   | Submit payment for an order |

🔑 Include header:\
`Authorization: Bearer <your_token_here>`

---

### 🔁 Square Integration

| Method                        | Endpoint        | Description                       |
| ----------------------------- | --------------- | --------------------------------- |
| POST                          | `/square/order` | Sends order to Square sandbox API |

## 📁 Project Structure 

```bash
.
├── main.go                  # Entry point
├── models/                 # DB models
├── database/               # DB connection setup
├── handlers/               # API route handlers
├── auth/                   # JWT middleware
├── square/                 # Square API interaction
├── .env                    # Environment variables
└── README.md
```
