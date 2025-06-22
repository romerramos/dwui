/*
  DWUI (Docker Web UI)
  Copyright (C) 2025 Romer Ramos

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
export default (containerId) => {
  return {
    terminal: null,
    socket: null,
    isConnected: false,
    containerId: containerId,
    fitAddon: null,
    fontSize: 12,
    minFontSize: 6,
    maxFontSize: 24,

    handleResize() {
      this.fitAddon.fit()
    },

    handleVisibilityChange() {
      if (!document.hidden && !this.isConnected) {
        this.connectWebSocket()
      }
    },

    init() {
      // Initialize xterm.js
      this.terminal = new Terminal({
        cursorBlink: true,
        fontFamily: "JetBrains Mono, Fira Code, Courier New, monospace",
        fontSize: this.fontSize,
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

      // Start connection
      this.connectWebSocket()
    },

    increaseFontSize() {
      if (this.fontSize < this.maxFontSize) {
        this.fontSize += 2
        this.terminal.options.fontSize = this.fontSize
        this.fitAddon.fit()
      }
    },

    decreaseFontSize() {
      if (this.fontSize > this.minFontSize) {
        this.fontSize -= 2
        this.terminal.options.fontSize = this.fontSize
        this.fitAddon.fit()
      }
    },

    connectWebSocket() {
      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"

      const locationHost = window.location.host.includes("8082")
        ? window.location.host.replace("8082", "8300")
        : window.location.host

      const wsUrl = `${protocol}//${locationHost}/terminal/stream/${this.containerId}`

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
    },
  }
}
