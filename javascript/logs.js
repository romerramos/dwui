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
export function createLogsComponent(containerId) {
  return {
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
    shiftKeyPressed: false,

    initLogs() {
      this.logLines = []

      this.updateDisplay()
      this.connectWebSocket()

      // For binding `this` into the Alpine.js component
      // and cleaning the event listener when the component is destroyed
      this.scrollHandler = () => this.handleScroll()
      this.$refs.logsElement.addEventListener("scroll", this.scrollHandler)

      this.visibilityHandler = () => {
        if (!document.hidden && !this.isConnected) {
          this.connectWebSocket()
        }
      }
      document.addEventListener("visibilitychange", this.visibilityHandler)

      this.keyDownHandler = (event) => {
        if (event.key === "f" && (event.ctrlKey || event.metaKey)) {
          event.preventDefault()
          event.stopPropagation()
          this.toggleSearch()
        }
      }

      document.addEventListener("keydown", this.keyDownHandler)
    },

    updateDisplay() {
      const content = this.logLines.join("\n")

      if (this.searchQuery && this.searchMatches.length > 0) {
        this.$refs.logsElement.innerHTML = this.highlightedContent
      } else {
        this.$refs.logsElement.textContent = content
      }

      this.updateFontSize()
      if (this.autoScroll && !this.userScrolledUp) {
        this.scrollToBottom()
      }
    },

    updateFontSize() {
      if (this.$refs.logsElement) {
        this.$refs.logsElement.style.fontSize = `${this.fontSize}px`
      }
    },

    scrollToBottom() {
      this.$refs.logsElement.scrollTop = this.$refs.logsElement.scrollHeight
      this.userScrolledUp = false
    },

    connectWebSocket() {
      const protocol = window.location.protocol === "https:" ? "wss:" : "ws:"

      // When running in dev mode the port 8082 is used for `air` live-reload,
      // but the sockets are running on 8300
      const locationHost = window.location.host.includes("8082")
        ? window.location.host.replace("8082", "8300")
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
          this.addLogLines(logData)
        }
      }

      this.socket.onclose = (event) => {
        console.log("WebSocket disconnected for logs")
        this.isConnected = false
        this.addLogLines(
          "\n--- Connection lost. Attempting to reconnect... ---",
        )
        setTimeout(() => this.connectWebSocket(), 3000)
      }

      this.socket.onerror = (error) => {
        console.error("WebSocket error:", error)
        this.addLogLines("\n--- Connection error occurred ---")
      }
    },

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

    performSearch() {
      if (!this.searchQuery) {
        this.clearSearch()
        return
      }

      const query = this.searchQuery.toLowerCase()
      const matches = []

      this.logLines.forEach((line, lineIndex) => {
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

    clearSearch() {
      this.searchMatches = []
      this.currentMatchIndex = 0
      this.highlightedContent = ""
      this.updateDisplay()
    },

    highlightMatches() {
      if (!this.searchQuery || this.searchMatches.length === 0) {
        return
      }

      const content = this.logLines.join("\n")
      const query = this.searchQuery
      const escapedQuery = query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")

      let safeContent = content
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")

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

    handleScroll() {
      const element = this.$refs.logsElement
      const isAtBottom =
        element.scrollTop + element.clientHeight >= element.scrollHeight - 5
      this.userScrolledUp = !isAtBottom

      if (this.userScrolledUp && this.autoScroll) {
        this.autoScroll = false
      }
    },

    toggleAutoScroll() {
      this.autoScroll = !this.autoScroll
      if (this.autoScroll) {
        this.scrollToBottom()
      }
    },

    increaseFontSize() {
      if (this.fontSize < this.maxFontSize) {
        this.fontSize += 2
        this.updateFontSize()
      }
    },

    decreaseFontSize() {
      if (this.fontSize > this.minFontSize) {
        this.fontSize -= 2
        this.updateFontSize()
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

    findNext() {
      if (this.shiftKeyPressed) return
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
      if (this.keyDownHandler) {
        document.removeEventListener("keydown", this.keyDownHandler)
      }
    },
  }
}
