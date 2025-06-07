// Terminal Alpine.js component
export function createTerminalComponent() {
  return (containerId) => ({
    terminal: null,
    socket: null,
    isConnected: false,
    containerId: containerId,
    fitAddon: null,

    initTerminal() {
      // Initialize xterm.js
      this.terminal = new Terminal({
        cursorBlink: true,
        fontFamily: "JetBrains Mono, Fira Code, Courier New, monospace",
        fontSize: 14,
        lineHeight: 1.2,
        theme: {
          background: "#1e2939",
          foreground: "#ffffff",
          cursor: "#ffffff",
          selection: "#ffffff40",
          black: "#000000",
          red: "#ff5555",
          green: "#50fa7b",
          yellow: "#f1fa8c",
          blue: "#bd93f9",
          magenta: "#ff79c6",
          cyan: "#8be9fd",
          white: "#bfbfbf",
        },
      })

      this.fitAddon = new FitAddon.FitAddon()
      this.terminal.loadAddon(this.fitAddon)
      this.terminal.open(this.$refs.terminalElement)
      this.fitAddon.fit()

      // Send terminal input to WebSocket
      this.terminal.onData((data) => {
        if (this.isConnected && this.socket.readyState === WebSocket.OPEN) {
          this.socket.send(data)
        }
      })

      // Handle window resize
      this.resizeHandler = () => {
        this.fitAddon.fit()
      }
      window.addEventListener("resize", this.resizeHandler)

      // Handle page visibility for reconnection
      this.visibilityHandler = () => {
        if (!document.hidden && !this.isConnected) {
          this.connectWebSocket()
        }
      }
      document.addEventListener("visibilitychange", this.visibilityHandler)

      // Start connection
      this.connectWebSocket()
    },

    connectWebSocket() {
      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"

      const locationHost = window.location.host.includes("8082")
        ? window.location.host.replace("8082", "8080")
        : window.location.host

      const wsUrl = `${protocol}//${locationHost}/terminal/${this.containerId}`

      this.socket = new WebSocket(wsUrl)

      this.socket.onopen = (event) => {
        console.log("WebSocket connected")
        this.isConnected = true
        this.terminal.writeln(
          "\x1b[32mConnected to container terminal...\x1b[0m\r\n",
        )
      }

      this.socket.onmessage = (event) => {
        this.terminal.write(event.data)
      }

      this.socket.onclose = (event) => {
        console.log("WebSocket disconnected")
        this.isConnected = false
        this.terminal.writeln(
          "\r\n\x1b[31mConnection lost. Attempting to reconnect...\x1b[0m",
        )

        // Attempt to reconnect after 3 seconds
        setTimeout(() => this.connectWebSocket(), 3000)
      }

      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error)
        this.terminal.writeln("\r\n\x1b[31mConnection error occurred.\x1b[0m")
      }
    },

    // Cleanup when component is destroyed
    destroy() {
      if (this.socket) {
        this.socket.close()
      }
      if (this.resizeHandler) {
        window.removeEventListener("resize", this.resizeHandler)
      }
      if (this.visibilityHandler) {
        document.removeEventListener("visibilitychange", this.visibilityHandler)
      }
    },
  })
}
