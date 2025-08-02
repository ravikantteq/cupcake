import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface KafkaMessage {
  broker: string;
  topic: string;
  key: string;
  value: string;
}

export interface ProducerHistory {
  id: string;
  broker: string;
  topic: string;
  key: string;
  value: string;
  success: boolean;
  response?: any;
  error?: string;
  timestamp: string;
  userId?: string;
}

export interface ApiResponse {
  success: boolean;
  message: string;
  data?: any;
}

export interface ErrorResponse {
  error: string;
  message: string;
}

@Injectable({
  providedIn: 'root'
})
export class KafkaService {
  private baseUrl = 'http://localhost:8080';

  constructor(private http: HttpClient) { }

  publishMessage(message: KafkaMessage): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/kafka/publish`, message);
  }

  healthCheck(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/health`);
  }

  // Get recent producer history (for initial load and caching)
  getRecentProducerHistory(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/history/recent`);
  }

  // Get producer history with pagination (for viewing more history)
  getProducerHistory(limit: number = 10, offset: number = 0): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/history?limit=${limit}&offset=${offset}`);
  }
}
