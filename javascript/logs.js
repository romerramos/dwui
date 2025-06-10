// Logs Alpine.js component for streaming container logs
export function createLogsComponent() {
  return (containerId) => ({
    socket: null,
    isConnected: false,
    containerId: containerId,
    logs: "",
    logLines: [], // Array to store individual log lines
    maxLines: 50, // Maximum number of lines to display (matches WebSocket initial)

    initLogs() {
      // Initialize with empty logs - WebSocket will populate everything
      this.logLines = []
      this.updateDisplay()

      // Start WebSocket connection immediately
      this.connectWebSocket()

      // Handle page visibility for reconnection
      this.visibilityHandler = () => {
        if (!document.hidden && !this.isConnected) {
          this.connectWebSocket()
        }
      }
      document.addEventListener("visibilitychange", this.visibilityHandler)
    },

    // Helper method to update the display with current log lines
    updateDisplay() {
      this.$refs.logsElement.textContent = this.logLines.join("\n")
      // Auto-scroll to bottom
      this.$refs.logsElement.scrollTop = this.$refs.logsElement.scrollHeight
    },

    // Helper method to add new lines while maintaining the limit
    addLogLines(newContent) {
      const newLines = newContent
        .split("\n")
        .filter((line) => line.trim() !== "")

      // Add new lines to the array
      this.logLines.push(...newLines)

      // Keep only the last maxLines
      if (this.logLines.length > this.maxLines) {
        this.logLines = this.logLines.slice(-this.maxLines)
      }

      this.updateDisplay()
    },

    connectWebSocket() {
      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"

      const locationHost = window.location.host.includes("8082")
        ? window.location.host.replace("8082", "8080")
        : window.location.host

      const wsUrl = `${protocol}//${locationHost}/logs/stream/${this.containerId}`

      this.socket = new WebSocket(wsUrl)

      this.socket.onopen = (event) => {
        console.log("WebSocket connected for logs")
        this.isConnected = true
      }

      this.socket.onmessage = (event) => {
        const logData = event.data
        if (logData.trim()) {
          // All messages are new streaming logs, add them
          this.addLogLines(logData)
        }
      }

      this.socket.onclose = (event) => {
        console.log("WebSocket disconnected for logs")
        this.isConnected = false
        this.addLogLines(
          "\n--- Connection lost. Attempting to reconnect... ---",
        )

        // Attempt to reconnect after 3 seconds
        setTimeout(() => this.connectWebSocket(), 3000)
      }

      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error)
        this.addLogLines("\n--- Connection error occurred ---")
      }
    },

    // Cleanup when component is destroyed
    destroy() {
      if (this.socket) {
        this.socket.close()
      }
      if (this.visibilityHandler) {
        document.removeEventListener("visibilitychange", this.visibilityHandler)
      }
    },
  })
}
