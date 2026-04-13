import type { WsEventPayloadMap, WsEventType, WsMessage } from './ws-events'

type EventHandler<T extends WsEventType> = (data: WsEventPayloadMap[T]) => void
type Handlers = Partial<{ [K in WsEventType]: EventHandler<K>[] }>

const RECONNECT_DELAY_MS = 3000

export class WsClient {
  private ws: WebSocket | null = null
  private handlers: Handlers = {}
  private url: string = ''
  private shouldReconnect = true
  private reconnectTimeout: ReturnType<typeof setTimeout> | null = null

  connect(url: string): void {
    this.url = url
    this.shouldReconnect = true
    this.open()
  }

  disconnect(): void {
    this.shouldReconnect = false
    if (this.reconnectTimeout) {
      clearTimeout(this.reconnectTimeout)
      this.reconnectTimeout = null
    }
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.close()
    }
    this.ws = null
  }

  on<T extends WsEventType>(event: T, handler: EventHandler<T>): () => void {
    if (!this.handlers[event]) {
      this.handlers[event] = []
    }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (this.handlers[event] as any[]).push(handler)

    // Return unsubscribe function
    return () => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      this.handlers[event] = (this.handlers[event] as any[]).filter((h) => h !== handler) as any
    }
  }

  send<T extends WsEventType>(type: T, data: WsEventPayloadMap[T]): void {
    this.sendRaw(type, data)
  }

  // For client→server messages that don't have a typed server→client counterpart
  sendRaw(type: string, data: unknown): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      console.warn('[WsClient] Cannot send — not connected')
      return
    }
    this.ws.send(JSON.stringify({ type, data }))
  }

  private open(): void {
    this.ws = new WebSocket(this.url)

    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data as string) as WsMessage
      const eventHandlers = this.handlers[message.type]
      if (eventHandlers) {
        for (const handler of eventHandlers) {
          (handler as EventHandler<typeof message.type>)(message.data)
        }
      } else {
        console.warn(`[WsClient] No handler for event: "${message.type}"`)
      }
    }

    this.ws.onclose = () => {
      if (this.shouldReconnect) {
        console.log(`[WsClient] Disconnected. Reconnecting in ${RECONNECT_DELAY_MS}ms...`)
        this.reconnectTimeout = setTimeout(() => this.open(), RECONNECT_DELAY_MS)
      }
    }

    this.ws.onerror = (event) => {
      console.error('[WsClient] Error:', event)
    }
  }
}
