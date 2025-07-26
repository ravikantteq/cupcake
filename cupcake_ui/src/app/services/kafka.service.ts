import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface KafkaMessage {
  broker: string;
  topic: string;
  key: string;
  value: string;
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
}
