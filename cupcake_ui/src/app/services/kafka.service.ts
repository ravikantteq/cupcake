import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface KafkaMessage {
  broker: string;
  topic: string;
  key: string;
  value: any; // Changed to 'any' to support both strings and JSON objects
  headers?: { [key: string]: string };
}

export interface Consumer {
  id?: string;
  name: string;
  description?: string;
  broker: string;
  groupId: string;
  topics: string[];
  status: 'inactive' | 'active' | 'error' | 'running' | 'stopped';
  config: ConsumerConfig;
  messageCount: number;
  lastHeartbeat?: string;
  startedAt?: string;
  stoppedAt?: string;
  errorMessage?: string;
  createdAt: string;
  updatedAt: string;
}

export interface ConsumerConfig {
  autoOffsetReset: string;
  enableAutoCommit: boolean;
  maxPollRecords: number;
  sessionTimeoutMs?: number;
  heartbeatIntervalMs?: number;
}

export interface Flow {
  id?: string;
  name: string;
  description: string;
  version: string;
  steps: FlowStep[];
  createdAt?: string;
  updatedAt?: string;
  createdBy?: string;
}

export interface FlowStep {
  stepId: string;
  type: 'produce' | 'consume' | 'validate' | 'delay';
  config: StepConfig;
}

export interface StepConfig {
  topic?: string;
  message?: any;
  expectedMessage?: any;
  timeout?: number;
  retries?: number;
  expectedCount?: number;
  delayMs?: number;
}

export interface Execution {
  id?: string;
  suiteId: string;
  flowId: string;
  status: 'inactive' | 'active' | 'error' | 'running' | 'pending' | 'completed' | 'failed' | 'cancelled';
  startTime: string;
  endTime?: string;
  steps: ExecutionStep[];
  metrics: ExecutionMetrics;
  logs?: string[];
}

export interface ExecutionStep {
  stepId: string;
  status: 'inactive' | 'active' | 'error' | 'running' | 'pending' | 'completed' | 'failed' | 'cancelled';
  input?: any;
  output?: any;
  errors?: string[];
  duration: number; // milliseconds
}

export interface ExecutionMetrics {
  totalDuration: number; // milliseconds
  messagesProduced: number;
  messagesConsumed: number;
  errorsCount: number;
  stepsCompleted: number;
  validationsPassed: number;
  validationsFailed: number;
}

export interface Message {
  id?: string;
  topic: string;
  partition: number;
  offset: number;
  key?: string;
  value: any;
  headers?: { [key: string]: string };
  timestamp: string;
  consumerGroupId: string;
  consumerId?: string;
  executionId?: string;
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

  // Kafka Message Publishing
  publishMessage(message: KafkaMessage): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/kafka/publish`, message);
  }

  // Health Check
  healthCheck(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/health`);
  }

  // Producer History
  getRecentProducerHistory(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/history/recent`);
  }

  getProducerHistory(limit: number = 10, offset: number = 0): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/history?limit=${limit}&offset=${offset}`);
  }

  // Consumer Management
  getConsumers(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/consumers`);
  }

  getConsumer(id: string): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/consumers/${id}`);
  }

  createConsumer(consumer: Partial<Consumer>): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/consumers`, consumer);
  }

  startConsumer(id: string): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/consumers/${id}/start`, {});
  }

  stopConsumer(id: string): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/consumers/${id}/stop`, {});
  }

  deleteConsumer(id: string): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(`${this.baseUrl}/api/v1/consumers/${id}`);
  }

  // Flow Management
  getFlows(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/flows`);
  }

  getFlow(id: string): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/flows/${id}`);
  }

  createFlow(flow: Partial<Flow>): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/flows`, flow);
  }

  updateFlow(id: string, flow: Partial<Flow>): Observable<ApiResponse> {
    return this.http.put<ApiResponse>(`${this.baseUrl}/api/v1/flows/${id}`, flow);
  }

  deleteFlow(id: string): Observable<ApiResponse> {
    return this.http.delete<ApiResponse>(`${this.baseUrl}/api/v1/flows/${id}`);
  }

  executeFlow(id: string, suiteId: string): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/flows/${id}/execute`, { suiteId });
  }

  // Execution Management
  getExecutions(): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/executions`);
  }

  getExecution(id: string): Observable<ApiResponse> {
    return this.http.get<ApiResponse>(`${this.baseUrl}/api/v1/executions/${id}`);
  }

  stopExecution(id: string): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/executions/${id}/stop`, {});
  }

  pauseExecution(id: string): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/executions/${id}/pause`, {});
  }

  continueExecution(id: string, stepId: string): Observable<ApiResponse> {
    return this.http.post<ApiResponse>(`${this.baseUrl}/api/v1/executions/${id}/continue/${stepId}`, {});
  }
}
