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
    fontSize: 12, // Font size in pixels
    minFontSize: 8, // Minimum font size
    maxFontSize: 24, // Maximum font size

    // Search functionality
    showSearch: false,
    searchQuery: "",
    searchMatches: [],
    currentMatchIndex: 0,
    highlightedContent: "",

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
      const content = this.logLines.join("\n")

      if (this.searchQuery && this.searchMatches.length > 0) {
        // Use highlighted content if search is active
        this.$refs.logsElement.innerHTML = this.highlightedContent
      } else {
        // Use plain text content
        this.$refs.logsElement.textContent = content
      }

      // Add scroll event listener after content is updated
      this.$nextTick(() => {
        this.$refs.logsElement.addEventListener("scroll", this.scrollHandler)

        // Apply current font size
        this.updateFontSize()

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

    // Increase font size
    increaseFontSize() {
      if (this.fontSize < this.maxFontSize) {
        this.fontSize += 2
        this.updateFontSize()
      }
    },

    // Decrease font size
    decreaseFontSize() {
      if (this.fontSize > this.minFontSize) {
        this.fontSize -= 2
        this.updateFontSize()
      }
    },

    // Update the font size of the logs element
    updateFontSize() {
      if (this.$refs.logsElement) {
        this.$refs.logsElement.style.fontSize = `${this.fontSize}px`
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

      // Refresh search if active
      if (this.searchQuery) {
        this.performSearch()
      }
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

    // Search functionality methods
    toggleSearch() {
      this.showSearch = !this.showSearch
      if (this.showSearch) {
        this.$nextTick(() => {
          this.$refs.searchInput.focus()
        })
      } else {
        this.hideSearch()
      }
    },

    hideSearch() {
      this.showSearch = false
      this.searchQuery = ""
      this.clearSearch()
    },

    performSearch() {
      if (!this.searchQuery) {
        this.clearSearch()
        return
      }

      const content = this.logLines.join("\n")
      const query = this.searchQuery.toLowerCase()
      const matches = []
      const lines = content.split("\n")

      // Find all matches and their positions
      lines.forEach((line, lineIndex) => {
        const lowerLine = line.toLowerCase()
        let index = 0
        while ((index = lowerLine.indexOf(query, index)) !== -1) {
          matches.push({
            lineIndex,
            charIndex: index,
            line: line,
          })
          index += query.length
        }
      })

      this.searchMatches = matches
      this.currentMatchIndex = 0

      if (matches.length > 0) {
        this.highlightMatches()
        this.scrollToCurrentMatch()
      } else {
        this.clearSearch()
      }
    },

    highlightMatches() {
      if (!this.searchQuery || this.searchMatches.length === 0) {
        return
      }

      const content = this.logLines.join("\n")
      const query = this.searchQuery
      const escapedQuery = query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")

      // First escape HTML entities
      let safeContent = content
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")

      // Then highlight matches
      let matchCount = 0
      const highlightedContent = safeContent.replace(
        new RegExp(escapedQuery, "gi"),
        (match) => {
          const isCurrentMatch = matchCount === this.currentMatchIndex
          const className = isCurrentMatch
            ? "bg-yellow-400 text-black"
            : "bg-yellow-200 text-black"
          matchCount++
          return `<span class="${className}">${match}</span>`
        },
      )

      this.highlightedContent = highlightedContent
    },

    findNext() {
      if (this.searchMatches.length === 0) return

      this.currentMatchIndex =
        (this.currentMatchIndex + 1) % this.searchMatches.length
      this.highlightMatches()
      this.updateDisplay()
      this.scrollToCurrentMatch()
    },

    findPrevious() {
      if (this.searchMatches.length === 0) return

      this.currentMatchIndex =
        this.currentMatchIndex === 0
          ? this.searchMatches.length - 1
          : this.currentMatchIndex - 1
      this.highlightMatches()
      this.updateDisplay()
      this.scrollToCurrentMatch()
    },

    scrollToCurrentMatch() {
      if (this.searchMatches.length === 0) return

      this.$nextTick(() => {
        const spans =
          this.$refs.logsElement.querySelectorAll("span.bg-yellow-400")
        if (spans.length > 0) {
          spans[0].scrollIntoView({
            behavior: "smooth",
            block: "center",
            inline: "nearest",
          })
        }
      })
    },

    clearSearch() {
      this.searchMatches = []
      this.currentMatchIndex = 0
      this.highlightedContent = ""
      this.updateDisplay()
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
