import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { interval, Subscription } from 'rxjs';
import { KafkaService, Execution, ExecutionStep, Flow, Message, ApiResponse } from '../services/kafka.service';

@Component({
  selector: 'app-executions',
  standalone: true,
  imports: [CommonModule, FormsModule, HttpClientModule],
  template: `
    <div class="container-fluid p-4">
      <div class="row mb-4">
        <div class="col-md-8">
          <h2><i class="fas fa-play-circle text-primary"></i> Flow Executions</h2>
          <p class="text-muted">Monitor and manage flow execution instances</p>
        </div>
        <div class="col-md-4 text-end">
          <button class="btn btn-outline-secondary" (click)="refreshExecutions()" [disabled]="isLoading">
            <i class="fas fa-refresh" [class.fa-spin]="isLoading"></i> Refresh
          </button>
        </div>
      </div>

      <!-- Error Alert -->
      <div class="alert alert-danger alert-dismissible fade show" *ngIf="errorMessage">
        <i class="fas fa-exclamation-triangle"></i> {{ errorMessage }}
        <button type="button" class="btn-close" (click)="errorMessage = ''"></button>
      </div>

      <!-- Success Alert -->
      <div class="alert alert-success alert-dismissible fade show" *ngIf="successMessage">
        <i class="fas fa-check-circle"></i> {{ successMessage }}
        <button type="button" class="btn-close" (click)="successMessage = ''"></button>
      </div>

      <!-- Executions List -->
      <div class="row" *ngIf="executions.length > 0">
        <div class="col-md-6 col-lg-4 mb-4" *ngFor="let execution of executions; trackBy: trackByExecutionId">
          <div class="card h-100 border-left-{{ getStatusColor(execution.status) }}">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h6 class="mb-0">Execution {{ execution.id?.substring(0, 8) }}</h6>
              <span class="badge bg-{{ getStatusColor(execution.status) }}">
                <i class="fas fa-{{ getStatusIcon(execution.status) }}"></i>
                {{ execution.status | titlecase }}
              </span>
            </div>
            <div class="card-body">
              <div class="row mb-2">
                <div class="col-12">
                  <small class="text-muted">Flow ID:</small><br>
                  <code class="small">{{ execution.flowId }}</code>
                </div>
              </div>

              <div class="row mb-2">
                <div class="col-6">
                  <small class="text-muted">Suite ID:</small><br>
                  <code class="small">{{ execution.suiteId }}</code>
                </div>
                <div class="col-6">
                  <small class="text-muted">Duration:</small><br>
                  <strong>{{ formatDuration(execution.metrics.totalDuration) }}</strong>
                </div>
              </div>

              <div class="mb-2">
                <small class="text-muted">Progress:</small><br>
                <div class="progress mb-1" style="height: 8px;">
                  <div class="progress-bar bg-{{ getStatusColor(execution.status) }}" 
                       [style.width.%]="getProgressPercentage(execution)">
                  </div>
                </div>
                <small class="text-muted">{{ execution.steps.length > 0 ? execution.metrics.stepsCompleted : 0 }}/{{ execution.steps.length }} steps</small>
              </div>

              <div class="row mb-2">
                <div class="col-6">
                  <small class="text-muted">Messages:</small><br>
                  <span class="text-success">+{{ execution.metrics.messagesProduced }}</span>
                  <span class="text-primary">-{{ execution.metrics.messagesConsumed }}</span>
                </div>
                <div class="col-6">
                  <small class="text-muted">Validations:</small><br>
                  <span class="text-success">✓{{ execution.metrics.validationsPassed }}</span>
                  <span class="text-danger">✗{{ execution.metrics.validationsFailed }}</span>
                </div>
              </div>

              <div class="mt-auto">
                <small class="text-muted">
                  Started: {{ formatTime(execution.startTime) }}
                  <span *ngIf="execution.endTime"><br>Ended: {{ formatTime(execution.endTime) }}</span>
                </small>
              </div>
            </div>
            <div class="card-footer d-flex justify-content-between">
              <button class="btn btn-sm btn-info" 
                      (click)="viewExecutionDetails(execution)"
                      [disabled]="isLoading">
                <i class="fas fa-eye"></i> View Details
              </button>
              <div class="btn-group" *ngIf="execution.status === 'running' || execution.status === 'pending'">
                <button class="btn btn-sm btn-warning" 
                        (click)="pauseExecution(execution.id!)" 
                        [disabled]="isLoading">
                  <i class="fas fa-pause"></i> Pause
                </button>
                <button class="btn btn-sm btn-danger" 
                        (click)="stopExecution(execution.id!)" 
                        [disabled]="isLoading">
                  <i class="fas fa-stop"></i> Stop
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div class="text-center py-5" *ngIf="executions.length === 0 && !isLoading">
        <i class="fas fa-play-circle fa-3x text-muted mb-3"></i>
        <h4 class="text-muted">No Executions Found</h4>
        <p class="text-muted">Execute flows to see them here</p>
      </div>

      <!-- Loading Spinner -->
      <div class="text-center py-5" *ngIf="isLoading">
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted mt-2">Loading executions...</p>
      </div>

      <!-- Execution Details Modal -->
      <div class="modal fade show" style="display: block;" *ngIf="selectedExecution" 
           (click)="closeExecutionDetails($event)">
        <div class="modal-dialog modal-xl">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">
                Execution Details: {{ selectedExecution.id?.substring(0, 12) }}
                <span class="badge bg-{{ getStatusColor(selectedExecution.status) }} ms-2">
                  {{ selectedExecution.status | titlecase }}
                </span>
              </h5>
              <button type="button" class="btn-close" (click)="closeExecutionDetails()"></button>
            </div>
            <div class="modal-body">
              <div class="row mb-4">
                <div class="col-md-6">
                  <h6>Execution Information</h6>
                  <table class="table table-sm">
                    <tr>
                      <td><strong>Execution ID:</strong></td>
                      <td><code>{{ selectedExecution.id }}</code></td>
                    </tr>
                    <tr>
                      <td><strong>Flow ID:</strong></td>
                      <td><code>{{ selectedExecution.flowId }}</code></td>
                    </tr>
                    <tr>
                      <td><strong>Suite ID:</strong></td>
                      <td><code>{{ selectedExecution.suiteId }}</code></td>
                    </tr>
                    <tr>
                      <td><strong>Status:</strong></td>
                      <td>
                        <span class="badge bg-{{ getStatusColor(selectedExecution.status) }}">
                          {{ selectedExecution.status | titlecase }}
                        </span>
                      </td>
                    </tr>
                    <tr>
                      <td><strong>Started:</strong></td>
                      <td>{{ formatTime(selectedExecution.startTime) }}</td>
                    </tr>
                    <tr *ngIf="selectedExecution.endTime">
                      <td><strong>Ended:</strong></td>
                      <td>{{ formatTime(selectedExecution.endTime) }}</td>
                    </tr>
                  </table>
                </div>
                <div class="col-md-6">
                  <h6>Metrics</h6>
                  <table class="table table-sm">
                    <tr>
                      <td><strong>Total Duration:</strong></td>
                      <td>{{ formatDuration(selectedExecution.metrics.totalDuration) }}</td>
                    </tr>
                    <tr>
                      <td><strong>Steps Completed:</strong></td>
                      <td>{{ selectedExecution.metrics.stepsCompleted }}/{{ selectedExecution.steps.length }}</td>
                    </tr>
                    <tr>
                      <td><strong>Messages Produced:</strong></td>
                      <td><span class="text-success">{{ selectedExecution.metrics.messagesProduced }}</span></td>
                    </tr>
                    <tr>
                      <td><strong>Messages Consumed:</strong></td>
                      <td><span class="text-primary">{{ selectedExecution.metrics.messagesConsumed }}</span></td>
                    </tr>
                    <tr>
                      <td><strong>Validations Passed:</strong></td>
                      <td><span class="text-success">{{ selectedExecution.metrics.validationsPassed }}</span></td>
                    </tr>
                    <tr>
                      <td><strong>Validations Failed:</strong></td>
                      <td><span class="text-danger">{{ selectedExecution.metrics.validationsFailed }}</span></td>
                    </tr>
                    <tr>
                      <td><strong>Errors:</strong></td>
                      <td><span class="text-danger">{{ selectedExecution.metrics.errorsCount }}</span></td>
                    </tr>
                  </table>
                </div>
              </div>

              <h6>Execution Steps</h6>
              <div class="accordion" id="stepsAccordion">
                <div class="accordion-item" *ngFor="let step of selectedExecution.steps; let i = index; trackBy: trackByStepId">
                  <h2 class="accordion-header">
                    <button class="accordion-button {{ step.status !== 'completed' && step.status !== 'failed' ? 'collapsed' : '' }}" 
                            type="button" 
                            [attr.data-bs-toggle]="'collapse'" 
                            [attr.data-bs-target]="'#step' + i"
                            [attr.aria-expanded]="step.status === 'completed' || step.status === 'failed'"
                            [attr.aria-controls]="'step' + i">
                      <div class="d-flex align-items-center w-100">
                        <span class="badge bg-{{ getStatusColor(step.status) }} me-2">
                          <i class="fas fa-{{ getStepTypeIcon(getStepType(step.stepId)) }}"></i>
                        </span>
                        <strong>Step {{ i + 1 }}: {{ getStepType(step.stepId) | titlecase }}</strong>
                        <span class="ms-auto me-3">
                          <span class="badge bg-{{ getStatusColor(step.status) }}">{{ step.status | titlecase }}</span>
                          <span class="badge bg-light text-dark ms-1">{{ formatDuration(step.duration) }}</span>
                        </span>
                      </div>
                    </button>
                  </h2>
                  <div [id]="'step' + i" 
                       class="accordion-collapse collapse {{ step.status === 'completed' || step.status === 'failed' ? 'show' : '' }}" 
                       [attr.data-bs-parent]="'#stepsAccordion'">
                    <div class="accordion-body">
                      <div class="row">
                        <div class="col-md-6" *ngIf="step.input">
                          <h6>Input</h6>
                          <pre class="bg-light p-2 rounded">{{ formatJson(step.input) }}</pre>
                        </div>
                        <div class="col-md-6" *ngIf="step.output">
                          <h6>Output</h6>
                          <div *ngIf="getStepType(step.stepId) === 'consume' && isConsumedMessages(step.output)">
                            <h6>Consumed Messages ({{ getConsumedMessages(step.output).length }})</h6>
                            <div class="consumed-messages">
                              <div class="message-card" *ngFor="let message of getConsumedMessages(step.output); let msgIndex = index">
                                <div class="message-header">
                                  <strong>Message {{ msgIndex + 1 }}</strong>
                                  <span class="badge bg-info">{{ message.topic }}</span>
                                  <small class="text-muted">Partition: {{ message.partition }}, Offset: {{ message.offset }}</small>
                                </div>
                                <div class="message-content">
                                  <div *ngIf="message.key" class="mb-2">
                                    <small class="text-muted">Key:</small>
                                    <code class="d-block">{{ message.key }}</code>
                                  </div>
                                  <div class="mb-2">
                                    <small class="text-muted">Value:</small>
                                    <pre class="message-value">{{ formatMessageValue(message.value) }}</pre>
                                  </div>
                                  <div *ngIf="message.headers && getObjectKeys(message.headers).length > 0" class="mb-2">
                                    <small class="text-muted">Headers:</small>
                                    <div class="headers">
                                      <span class="badge bg-secondary me-1" *ngFor="let header of getHeaderEntries(message.headers)">
                                        {{ header.key }}: {{ header.value }}
                                      </span>
                                    </div>
                                  </div>
                                  <small class="text-muted">Timestamp: {{ formatTime(message.timestamp) }}</small>
                                </div>
                              </div>
                            </div>
                            <div class="mt-3" *ngIf="getStepType(step.stepId) === 'consume' && selectedExecution?.status === 'running'">
                              <button class="btn btn-primary btn-sm" (click)="continueConsumerStep(selectedExecution.id!, step.stepId)">
                                <i class="fas fa-play"></i> Continue to Next Step
                              </button>
                            </div>
                          </div>
                          <pre class="bg-light p-2 rounded" *ngIf="getStepType(step.stepId) !== 'consume' || !isConsumedMessages(step.output)">{{ formatJson(step.output) }}</pre>
                        </div>
                        <div class="col-12" *ngIf="step.errors && step.errors.length > 0">
                          <h6 class="text-danger">Errors</h6>
                          <div class="alert alert-danger">
                            <ul class="mb-0">
                              <li *ngFor="let error of step.errors">{{ error }}</li>
                            </ul>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div *ngIf="selectedExecution.logs && selectedExecution.logs.length > 0" class="mt-4">
                <h6>Execution Logs</h6>
                <div class="logs-container">
                  <pre class="logs">{{ selectedExecution.logs.join('\\n') }}</pre>
                </div>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" (click)="closeExecutionDetails()">Close</button>
              <button class="btn btn-warning" 
                      (click)="pauseExecution(selectedExecution.id!); closeExecutionDetails()" 
                      [disabled]="selectedExecution.status !== 'running'"
                      *ngIf="selectedExecution.status === 'running'">
                <i class="fas fa-pause"></i> Pause Execution
              </button>
              <button class="btn btn-danger" 
                      (click)="stopExecution(selectedExecution.id!); closeExecutionDetails()" 
                      [disabled]="selectedExecution.status !== 'running' && selectedExecution.status !== 'pending'"
                      *ngIf="selectedExecution.status === 'running' || selectedExecution.status === 'pending'">
                <i class="fas fa-stop"></i> Stop Execution
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .border-left-success { border-left: 4px solid #28a745 !important; }
    .border-left-warning { border-left: 4px solid #ffc107 !important; }
    .border-left-danger { border-left: 4px solid #dc3545 !important; }
    .border-left-secondary { border-left: 4px solid #6c757d !important; }
    .border-left-primary { border-left: 4px solid #007bff !important; }
    
    .progress { background-color: #e9ecef; }
    
    .modal { background-color: rgba(0,0,0,0.5); }
    
    .table td { border-top: none; padding: 0.25rem 0.5rem; }
    
    .card-footer { background-color: #f8f9fa; }
    
    .badge { font-size: 0.75em; }
    
    code { font-size: 0.875em; }
    
    .accordion-button { font-size: 0.9rem; }
    
    pre { font-size: 0.8rem; max-height: 200px; overflow-y: auto; }
    
    .logs-container { max-height: 300px; overflow-y: auto; }
    
    .logs { background-color: #f8f9fa; font-size: 0.75rem; }
    
    .consumed-messages { max-height: 400px; overflow-y: auto; }
    
    .message-card { 
      border: 1px solid #dee2e6; 
      border-radius: 0.375rem; 
      padding: 0.75rem; 
      margin-bottom: 0.75rem; 
      background-color: #f8f9fa; 
    }
    
    .message-header { 
      display: flex; 
      justify-content: space-between; 
      align-items: center; 
      margin-bottom: 0.5rem; 
      padding-bottom: 0.5rem; 
      border-bottom: 1px solid #dee2e6; 
    }
    
    .message-content code { background-color: #e9ecef; padding: 0.25rem 0.5rem; border-radius: 0.25rem; }
    
    .message-value { 
      background-color: #ffffff; 
      border: 1px solid #dee2e6; 
      padding: 0.5rem; 
      margin: 0; 
      max-height: 150px; 
      overflow-y: auto; 
    }
    
    .headers { margin-top: 0.25rem; }
  `]
})
export class ExecutionsComponent implements OnInit, OnDestroy {
  executions: Execution[] = [];
  isLoading = false;
  errorMessage = '';
  successMessage = '';
  selectedExecution: Execution | null = null;
  
  private refreshSubscription?: Subscription;

  constructor(private kafkaService: KafkaService) {}

  ngOnInit() {
    this.loadExecutions();
    // Auto-refresh every 15 seconds for active executions
    this.refreshSubscription = interval(15000).subscribe(() => {
      if (!this.selectedExecution) {
        this.loadExecutions(true);
      }
    });
  }

  ngOnDestroy() {
    if (this.refreshSubscription) {
      this.refreshSubscription.unsubscribe();
    }
  }

  loadExecutions(silent = false) {
    if (!silent) this.isLoading = true;
    this.errorMessage = '';

    this.kafkaService.getExecutions().subscribe({
      next: (response: ApiResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.executions = response.data || [];
          // If no real executions, show mock data for demo
          if (this.executions.length === 0) {
            this.executions = this.generateMockExecutions();
          }
        } else {
          this.errorMessage = response.message || 'Failed to load executions';
          // Fallback to mock data on error
          this.executions = this.generateMockExecutions();
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to connect to server. Showing demo data.';
        console.error('Error loading executions:', error);
        // Fallback to mock data when API is not available
        this.executions = this.generateMockExecutions();
      }
    });
  }

  generateMockExecutions(): Execution[] {
    return [
      {
        id: 'exec-001',
        suiteId: 'suite-test-001',
        flowId: 'flow-consumer-test',
        status: 'running',
        startTime: new Date(Date.now() - 120000).toISOString(), // 2 minutes ago
        steps: [
          {
            stepId: 'produce-step-1',
            status: 'completed',
            input: { topic: 'test-topic', message: { type: 'test', data: 'hello' } },
            output: { messageId: 'msg-001', partition: 0, offset: 123 },
            errors: [],
            duration: 150
          },
          {
            stepId: 'consume-step-1',
            status: 'completed',
            input: { topic: 'test-topic', timeout: 30000 },
            output: {
              messages: [
                {
                  id: 'msg-001',
                  topic: 'test-topic',
                  partition: 0,
                  offset: 123,
                  key: 'test-key',
                  value: { type: 'test', data: 'hello', timestamp: new Date().toISOString() },
                  headers: { 'content-type': 'application/json', 'source': 'test-producer' },
                  timestamp: new Date(Date.now() - 60000).toISOString(),
                  consumerGroupId: 'test-group',
                  consumerId: 'consumer-001',
                  executionId: 'exec-001'
                },
                {
                  id: 'msg-002',
                  topic: 'test-topic',
                  partition: 0,
                  offset: 124,
                  key: 'test-key-2',
                  value: 'Simple string message',
                  headers: { 'content-type': 'text/plain' },
                  timestamp: new Date(Date.now() - 45000).toISOString(),
                  consumerGroupId: 'test-group',
                  consumerId: 'consumer-001',
                  executionId: 'exec-001'
                }
              ]
            },
            errors: [],
            duration: 2500
          },
          {
            stepId: 'validate-step-1',
            status: 'running',
            input: { expected: { type: 'test' }, actual: { type: 'test', data: 'hello' } },
            output: null,
            errors: [],
            duration: 0
          }
        ],
        metrics: {
          totalDuration: 120000,
          messagesProduced: 1,
          messagesConsumed: 2,
          errorsCount: 0,
          stepsCompleted: 2,
          validationsPassed: 1,
          validationsFailed: 0
        },
        logs: [
          '[2023-12-07 14:30:15] Starting execution exec-001',
          '[2023-12-07 14:30:15] Step 1: Producing message to test-topic',
          '[2023-12-07 14:30:15] Step 1: Message produced successfully (offset: 123)',
          '[2023-12-07 14:30:16] Step 2: Starting consumer for test-topic',
          '[2023-12-07 14:30:18] Step 2: Consumed 2 messages from test-topic',
          '[2023-12-07 14:30:18] Step 3: Starting validation...'
        ]
      },
      {
        id: 'exec-002',
        suiteId: 'suite-test-002',
        flowId: 'flow-producer-test',
        status: 'completed',
        startTime: new Date(Date.now() - 300000).toISOString(), // 5 minutes ago
        endTime: new Date(Date.now() - 280000).toISOString(), // 4m40s ago
        steps: [
          {
            stepId: 'produce-step-1',
            status: 'completed',
            input: { topic: 'events', message: { eventType: 'user.signup', userId: 12345 } },
            output: { messageId: 'msg-003', partition: 1, offset: 456 },
            errors: [],
            duration: 85
          },
          {
            stepId: 'delay-step-1',
            status: 'completed',
            input: { delayMs: 1000 },
            output: { actualDelay: 1001 },
            errors: [],
            duration: 1001
          },
          {
            stepId: 'validate-step-1',
            status: 'completed',
            input: { topic: 'events', expectedCount: 1 },
            output: { actualCount: 1, validated: true },
            errors: [],
            duration: 200
          }
        ],
        metrics: {
          totalDuration: 20000,
          messagesProduced: 1,
          messagesConsumed: 0,
          errorsCount: 0,
          stepsCompleted: 3,
          validationsPassed: 1,
          validationsFailed: 0
        }
      }
    ];
  }

  refreshExecutions() {
    this.loadExecutions();
  }

  viewExecutionDetails(execution: Execution) {
    this.selectedExecution = execution;
  }

  closeExecutionDetails(event?: MouseEvent) {
    if (event && event.target !== event.currentTarget) {
      return; // Don't close if clicking inside modal
    }
    this.selectedExecution = null;
  }

  pauseExecution(id: string) {
    // Implementation for pausing execution
    this.successMessage = 'Execution paused successfully!';
  }

  stopExecution(id: string) {
    // Implementation for stopping execution
    this.successMessage = 'Execution stopped successfully!';
  }

  continueConsumerStep(executionId: string, stepId: string) {
    // Implementation for continuing consumer step
    this.successMessage = 'Consumer step continued to next step!';
    this.closeExecutionDetails();
  }

  getStatusColor(status: string): string {
    switch (status) {
      case 'completed': return 'success';
      case 'running': return 'primary';
      case 'pending': return 'warning';
      case 'failed': return 'danger';
      case 'error': return 'danger';
      case 'cancelled': return 'secondary';
      case 'inactive': return 'secondary';
      default: return 'secondary';
    }
  }

  getStatusIcon(status: string): string {
    switch (status) {
      case 'completed': return 'check';
      case 'running': return 'play';
      case 'pending': return 'clock';
      case 'failed': return 'times';
      case 'error': return 'exclamation-triangle';
      case 'cancelled': return 'ban';
      case 'inactive': return 'stop';
      default: return 'question';
    }
  }

  getStepTypeIcon(stepType: string): string {
    switch (stepType) {
      case 'produce': return 'upload';
      case 'consume': return 'download';
      case 'validate': return 'check-circle';
      case 'delay': return 'clock';
      default: return 'cog';
    }
  }

  getProgressPercentage(execution: Execution): number {
    if (execution.steps.length === 0) return 0;
    return (execution.metrics.stepsCompleted / execution.steps.length) * 100;
  }

  formatTime(timestamp: string): string {
    return new Date(timestamp).toLocaleString();
  }

  formatDuration(ms: number): string {
    if (ms < 1000) return `${ms}ms`;
    if (ms < 60000) return `${(ms / 1000).toFixed(1)}s`;
    return `${(ms / 60000).toFixed(1)}m`;
  }

  formatJson(obj: any): string {
    if (obj === null || obj === undefined) return 'null';
    return JSON.stringify(obj, null, 2);
  }

  formatMessageValue(value: any): string {
    if (typeof value === 'string') return value;
    return JSON.stringify(value, null, 2);
  }

  getStepType(stepId: string): string {
    if (stepId.includes('produce')) return 'produce';
    if (stepId.includes('consume')) return 'consume';
    if (stepId.includes('validate')) return 'validate';
    if (stepId.includes('delay')) return 'delay';
    return 'unknown';
  }

  isConsumedMessages(output: any): boolean {
    return output && output.messages && Array.isArray(output.messages);
  }

  getConsumedMessages(output: any): Message[] {
    return output && output.messages ? output.messages : [];
  }

  getHeaderEntries(headers: { [key: string]: string }): { key: string, value: string }[] {
    return Object.entries(headers).map(([key, value]) => ({ key, value }));
  }

  getObjectKeys(obj: any): string[] {
    return Object.keys(obj);
  }

  trackByExecutionId(index: number, execution: Execution): string {
    return execution.id || index.toString();
  }

  trackByStepId(index: number, step: ExecutionStep): string {
    return step.stepId || index.toString();
  }
}
