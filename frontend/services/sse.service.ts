import { Product } from "@/types/products.type";

const API_URI = process.env.NEXT_PUBLIC_API_URI || 'http://localhost:8080';

type SSECallback<T> = (data: T) => void;
type SSEErrorCallback = (error: Event) => void;
type SSEOpenCallback = () => void;

export class SSEService<T> {
  private eventSource: EventSource | null = null;
  private endpoint: string;
  private reconnectTimeout: number;
  private onMessageCallback: SSECallback<T> | null = null;
  private onErrorCallback: SSEErrorCallback | null = null;
  private onOpenCallback: SSEOpenCallback | null = null;

  constructor(endpoint: string, reconnectTimeout = 5000) {
    this.endpoint = `${API_URI}${endpoint}`;
    this.reconnectTimeout = reconnectTimeout;
  }

  connect(): void {
    if (this.eventSource) {
      this.disconnect();
    }

    this.eventSource = new EventSource(this.endpoint);
    
    this.eventSource.onopen = () => {
      console.log(`SSE connection opened for ${this.endpoint}`);
      if (this.onOpenCallback) {
        this.onOpenCallback();
      }
    };

    this.eventSource.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data) as T;
        if (this.onMessageCallback) {
          this.onMessageCallback(data);
        }
      } catch (error) {
        console.error('Error parsing SSE event:', error);
      }
    };

    this.eventSource.onerror = (error) => {
      console.error('SSE connection error:', error);
      if (this.onErrorCallback) {
        this.onErrorCallback(error);
      }
      
      this.disconnect();
      
      setTimeout(() => {
        console.log('Attempting to reconnect to SSE...');
        this.connect();
      }, this.reconnectTimeout);
    };
  }

  disconnect(): void {
    if (this.eventSource) {
      console.log('Closing SSE connection');
      this.eventSource.close();
      this.eventSource = null;
    }
  }

  onMessage(callback: SSECallback<T>): void {
    this.onMessageCallback = callback;
  }

  onError(callback: SSEErrorCallback): void {
    this.onErrorCallback = callback;
  }

  onOpen(callback: SSEOpenCallback): void {
    this.onOpenCallback = callback;
  }
}

export const createProductRestockSSE = () => {
  return new SSEService<Product>('/api/sse/restocks');
};