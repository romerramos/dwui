# Docker Web UI (DWUI)

A web interface for managing Docker containers.

### Running with Auto-Generated Password

```bash
go run main.go
```

The application will generate a random password and display it in the terminal:

```
🔐 DWUI Authentication Password: cf674507f571
🌐 Server will be available at :8080
💡 Use this password to sign in to the web interface
```

### Running with Custom Password

```bash
go run main.go --password yourpassword
```

# Development

Run in separate terminals:

- `air`
- `npx @tailwindcss/cli -i ./tailwind.css -o ./assets/stylesheets/output.css --watch`

# Build

Linux: `GOOS=linux GOARCH=amd64 go build -o dwui-linux`

# Deploy

The main branch should contain the latest build in the root of the repo, e.g: `dwui-linux`, scp/clone/ftp/move that file to a server running docker, execute it and that's it!

```bash
./dwui-linux --password yourpassword
```

## License

This project is licensed under the AGPL v3.0 license. See the [LICENSE](LICENSE) file for more details.

Copyright (C) 2025 Romer Ramos
