# Docker Web UI (DWUI)

A secure web interface for managing Docker containers with built-in authentication.

## 🔐 Authentication

DWUI includes a simple but effective authentication system to protect access to your Docker environment.

### Running with Auto-Generated Password

```bash
go run main.go
```

The application will generate a random password and display it in the terminal:

```
🔐 DWUI Authentication Password: cf674507f571
🌐 Server will be available at: http://localhost:8080
💡 Use this password to sign in to the web interface
```

### Running with Custom Password

```bash
go run main.go --password yourpassword
```

### Security Features

- 🔒 Password-based authentication
- 🍪 Secure session management with HTTP-only cookies
- ⏰ 24-hour session expiration with auto-renewal
- 🚫 Automatic redirection to sign-in for unauthorized access
- 🛡️ Protection for all Docker management endpoints

### 🏗️ **Architecture**

- **`cmd/auth/middleware.go`**: Core authentication logic and session management
- **`cmd/auth/handlers.go`**: Sign-in/sign-out request handlers
- **`cmd/auth/signin.gohtml`**: Beautiful sign-in page template
- **Updated `main.go`**: Integrated authentication system with proper routing

# Development

`air`
`npx @tailwindcss/cli -i ./tailwind.css -o ./assets/stylesheets/output.css --watch`

# Build

`GOOS=linux GOARCH=amd64 go build -o dwui-linux`
