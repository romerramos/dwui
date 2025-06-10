// Logs Alpine.js component for streaming container logs
export function createLogsComponent() {
  return (containerId) => ({
    socket: null,
    isConnected: false,
    containerId: containerId,
    logs: "",
    logLines: [], // Array to store individual log lines
    autoScroll: true, // Auto-scroll toggle state
    userScrolledUp: false, // Track if user manually scrolled up

    initLogs() {
      // Initialize with empty logs - WebSocket will populate everything
      this.logLines = []
      this.updateDisplay()

      // Start WebSocket connection immediately
      this.connectWebSocket()

      // Add scroll event listener to detect user scrolling
      this.scrollHandler = () => {
        this.handleScroll()
      }

      // Handle page visibility for reconnection
      this.visibilityHandler = () => {
        if (!document.hidden && !this.isConnected) {
          this.connectWebSocket()
        }
      }
      document.addEventListener("visibilitychange", this.visibilityHandler)
    },

    // Handle scroll events to detect if user scrolled up
    handleScroll() {
      const element = this.$refs.logsElement
      const isAtBottom =
        element.scrollTop + element.clientHeight >= element.scrollHeight - 5
      this.userScrolledUp = !isAtBottom

      // Turn off auto-scroll when user manually scrolls up
      if (this.userScrolledUp && this.autoScroll) {
        this.autoScroll = false
      }
    },

    // Helper method to update the display with current log lines
    updateDisplay() {
      this.$refs.logsElement.textContent = this.logLines.join("\n")

      // Add scroll event listener after content is updated
      this.$nextTick(() => {
        this.$refs.logsElement.addEventListener("scroll", this.scrollHandler)

        // Auto-scroll to bottom only if auto-scroll is enabled and user hasn't scrolled up
        if (this.autoScroll && !this.userScrolledUp) {
          this.scrollToBottom()
        }
      })
    },

    // Scroll to bottom method
    scrollToBottom() {
      this.$refs.logsElement.scrollTop = this.$refs.logsElement.scrollHeight
      this.userScrolledUp = false
    },

    // Toggle auto-scroll and scroll to bottom if enabling
    toggleAutoScroll() {
      this.autoScroll = !this.autoScroll
      if (this.autoScroll) {
        this.scrollToBottom()
      }
    },

    // Helper method to add new lines
    addLogLines(newContent) {
      const newLines = newContent
        .split("\n")
        .filter((line) => line.trim() !== "")

      // Add new lines to the array (no limit)
      this.logLines.push(...newLines)

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
      if (this.scrollHandler && this.$refs.logsElement) {
        this.$refs.logsElement.removeEventListener("scroll", this.scrollHandler)
      }
    },
  })
}
