# DWUI: Docker Web UI

![DWUI Logo](assets/images/dwui.png)

DWUI is a simple web interface for Docker that gives you a "Docker Desktop" like experience, right on your server. It's designed to be lightweight, easy to use, and a perfect companion for deployment tools like [Kamal](https://kamal-deploy.org/).

Born out of a need for a no-fuss container management tool and a passion for learning Go, DWUI aims to be simple, effective, and open to the community.

## Screenshot

_(A screenshot of the UI will be added here soon.)_

## Features

- **View Containers**: See all your running containers and their status at a glance.
- **Inspect Details**: Check environment variables and open ports.
- **Real-time Logs**: Stream container logs directly in your browser.
- **Web Terminal**: Open an interactive terminal into your containers.
- **Single Binary**: No dependencies or complex setup. Just one file to run.
- **Kamal-Friendly**: A great companion to your Kamal deployment workflow.
- **Responsive**: Access it from your desktop or on the go from your phone.

## Getting Started

### Quick Install (One-Liner)

You can download and run DWUI with a single command.

> **Note**: You will need to create a release on GitHub for this URL to be valid. Please replace `v0.1.0` with your desired release version and `romerramos/dwui` with your repository path if it's different.

```bash
# Download the binary for Linux
curl -L -o dwui-linux https://github.com/romerramos/dwui/releases/download/v0.1.0/dwui-linux

# Make it executable
chmod +x dwui-linux

# Run it!
./dwui-linux --password your-very-secure-password
```

The server will be available at `http://<your-server-ip>:8080`.

### Running as a Service (systemd)

For a more permanent setup on Linux servers using `systemd`, you can run DWUI as a service.

1.  **Download and move the binary:**

    ```bash
    # Download the binary
    curl -L -o dwui-linux https://github.com/romerramos/dwui/releases/download/v0.1.0/dwui-linux
    chmod +x dwui-linux

    # Move it to a directory in your PATH
    sudo mv dwui-linux /usr/local/bin/dwui
    ```

2.  **Create a service file:**

    Create a new file at `/etc/systemd/system/dwui.service`:

    ```bash
    sudo nano /etc/systemd/system/dwui.service
    ```

    Paste the following content into the file. Make sure to **change the password**!

    ```ini
    [Unit]
    Description=DWUI - Docker Web UI
    After=docker.service
    Requires=docker.service

    [Service]
    ExecStart=/usr/local/bin/dwui --password your-very-secure-password
    Restart=always
    User=root # Or another user that has access to the Docker socket

    [Install]
    WantedBy=multi-user.target
    ```

3.  **Enable and start the service:**

    ```bash
    # Reload systemd to recognize the new service
    sudo systemctl daemon-reload

    # Enable DWUI to start on boot
    sudo systemctl enable dwui

    # Start the service immediately
    sudo systemctl start dwui

    # Check the status to make sure it's running
    sudo systemctl status dwui
    ```

## Development

To run the project in a development environment, you'll need two separate terminal sessions:

1.  **Run the Go backend with live-reloading:**

    ```bash
    air
    ```

2.  **Run the TailwindCSS compiler in watch mode:**
    ```bash
    npx @tailwindcss/cli -i ./tailwind.css -o ./assets/stylesheets/output.css --watch
    ```

## Building from Source

You can build the binary from the source code.

- **For Linux:**
  ```bash
  GOOS=linux GOARCH=amd64 go build -o ./dist/dwui-linux
  ```

## Contributing

This project was built to solve a personal need and as a Go learning experience. Contributions, feedback, and suggestions from the community are highly welcome! Feel free to open an issue or a pull request.

## License

This project is licensed under the AGPL v3.0 license. See the [LICENSE](LICENSE) file for more details.

Copyright (C) 2025 Romer Ramos
