import { ref, onUnmounted } from 'vue'

interface WebSocketOptions {
  onMessage?: (data: string) => void
  onConnect?: () => void
  onDisconnect?: () => void
  onError?: (error: Event) => void
  reconnectInterval?: number
  maxReconnectAttempts?: number
}

export function useWebSocket() {
  const ws = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const reconnectAttempts = ref(0)
  
  let currentUrl = ''
  let options: WebSocketOptions = {}
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null

  function connect(url: string, opts: WebSocketOptions = {}) {
    currentUrl = url
    options = {
      reconnectInterval: 3000,
      maxReconnectAttempts: 5,
      ...opts
    }
    
    createConnection()
  }

  function createConnection() {
    if (ws.value) {
      ws.value.close()
    }

    try {
      ws.value = new WebSocket(currentUrl)
      
      ws.value.onopen = () => {
        isConnected.value = true
        reconnectAttempts.value = 0
        options.onConnect?.()
      }

      ws.value.onmessage = (event) => {
        options.onMessage?.(event.data)
      }

      ws.value.onclose = () => {
        isConnected.value = false
        options.onDisconnect?.()
        attemptReconnect()
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket error:', error)
        options.onError?.(error)
      }
    } catch (error) {
      console.error('WebSocket connection error:', error)
    }
  }

  function attemptReconnect() {
    if (reconnectAttempts.value >= (options.maxReconnectAttempts || 5)) {
      console.log('Max reconnect attempts reached')
      return
    }

    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
    }

    reconnectTimer = setTimeout(() => {
      reconnectAttempts.value++
      console.log(`Reconnecting... Attempt ${reconnectAttempts.value}`)
      createConnection()
    }, options.reconnectInterval || 3000)
  }

  function disconnect() {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    
    reconnectAttempts.value = (options.maxReconnectAttempts || 5) // Prevent auto-reconnect
    
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    
    isConnected.value = false
  }

  function send(data: string | object) {
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      const message = typeof data === 'string' ? data : JSON.stringify(data)
      ws.value.send(message)
      return true
    }
    return false
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    connect,
    disconnect,
    send,
    isConnected
  }
}
