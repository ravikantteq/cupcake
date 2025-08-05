import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { interval, Subscription } from 'rxjs';
import { KafkaService, Consumer, ConsumerConfig, ApiResponse } from '../services/kafka.service';

@Component({
  selector: 'app-consumers',
  standalone: true,
  imports: [CommonModule, FormsModule, HttpClientModule],
  template: `
    <div class="container-fluid p-4">
      <div class="row mb-4">
        <div class="col-md-8">
          <h2><i class="fas fa-download text-primary"></i> Kafka Consumers</h2>
          <p class="text-muted">Manage and monitor Kafka consumer instances</p>
        </div>
        <div class="col-md-4 text-end">
          <button class="btn btn-primary" (click)="showCreateForm = true" [disabled]="showCreateForm">
            <i class="fas fa-plus"></i> Create Consumer
          </button>
          <button class="btn btn-outline-secondary ms-2" (click)="refreshConsumers()" [disabled]="isLoading">
            <i class="fas fa-refresh" [class.fa-spin]="isLoading"></i> Refresh
          </button>
        </div>
      </div>

      <!-- Create Consumer Form -->
      <div class="card mb-4" *ngIf="showCreateForm">
        <div class="card-header d-flex justify-content-between align-items-center">
          <h5 class="mb-0">Create New Consumer</h5>
          <button class="btn btn-sm btn-outline-secondary" (click)="cancelCreate()">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="card-body">
          <form (ngSubmit)="createConsumer()" #createForm="ngForm">
            <div class="row">
              <div class="col-md-6">
                <div class="mb-3">
                  <label class="form-label">Consumer Name *</label>
                  <input type="text" class="form-control" [(ngModel)]="newConsumer.name" name="name" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Description</label>
                  <textarea class="form-control" [(ngModel)]="newConsumer.description" name="description" rows="2"></textarea>
                </div>
                <div class="mb-3">
                  <label class="form-label">Broker *</label>
                  <input type="text" class="form-control" [(ngModel)]="newConsumer.broker" name="broker" 
                         placeholder="localhost:9092" required>
                </div>
                <div class="mb-3">
                  <label class="form-label">Group ID *</label>
                  <input type="text" class="form-control" [(ngModel)]="newConsumer.groupId" name="groupId" required>
                </div>
              </div>
              <div class="col-md-6">
                <div class="mb-3">
                  <label class="form-label">Topics *</label>
                  <input type="text" class="form-control" [(ngModel)]="topicsInput" name="topics" 
                         placeholder="topic1,topic2,topic3" required>
                  <small class="form-text text-muted">Comma-separated list of topics</small>
                </div>
                <div class="mb-3">
                  <label class="form-label">Auto Offset Reset</label>
                  <select class="form-control" [(ngModel)]="newConsumer.config!.autoOffsetReset" name="autoOffsetReset">
                    <option value="earliest">Earliest</option>
                    <option value="latest">Latest</option>
                  </select>
                </div>
                <div class="mb-3">
                  <label class="form-label">Max Poll Records</label>
                  <input type="number" class="form-control" [(ngModel)]="newConsumer.config!.maxPollRecords" 
                         name="maxPollRecords" min="1" max="10000">
                </div>
                <div class="form-check mb-3">
                  <input type="checkbox" class="form-check-input" [(ngModel)]="newConsumer.config!.enableAutoCommit" 
                         name="enableAutoCommit" id="enableAutoCommit">
                  <label class="form-check-label" for="enableAutoCommit">Enable Auto Commit</label>
                </div>
              </div>
            </div>
            <div class="d-flex justify-content-end">
              <button type="button" class="btn btn-secondary me-2" (click)="cancelCreate()">Cancel</button>
              <button type="submit" class="btn btn-primary" [disabled]="!createForm.form.valid || isLoading">
                <i class="fas fa-save"></i> Create Consumer
              </button>
            </div>
          </form>
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

      <!-- Consumers Grid -->
      <div class="row" *ngIf="consumers.length > 0">
        <div class="col-md-6 col-lg-4 mb-4" *ngFor="let consumer of consumers; trackBy: trackByConsumerId">
          <div class="card h-100 border-left-{{ getStatusColor(consumer.status) }}">
            <div class="card-header d-flex justify-content-between align-items-center">
              <h6 class="mb-0">{{ consumer.name }}</h6>
              <span class="badge bg-{{ getStatusColor(consumer.status) }}">
                <i class="fas fa-{{ getStatusIcon(consumer.status) }}"></i>
                {{ consumer.status | titlecase }}
              </span>
            </div>
            <div class="card-body">
              <p class="text-muted small mb-2" *ngIf="consumer.description">{{ consumer.description }}</p>
              
              <div class="row mb-2">
                <div class="col-6">
                  <small class="text-muted">Broker:</small><br>
                  <code class="small">{{ consumer.broker }}</code>
                </div>
                <div class="col-6">
                  <small class="text-muted">Group ID:</small><br>
                  <code class="small">{{ consumer.groupId }}</code>
                </div>
              </div>

              <div class="mb-2">
                <small class="text-muted">Topics:</small><br>
                <span class="badge bg-light text-dark me-1" *ngFor="let topic of consumer.topics">{{ topic }}</span>
              </div>

              <div class="row mb-2">
                <div class="col-6">
                  <small class="text-muted">Messages:</small><br>
                  <strong class="text-primary">{{ consumer.messageCount | number }}</strong>
                </div>
                <div class="col-6" *ngIf="isConsumerActive(consumer) && consumer.lastHeartbeat">
                  <small class="text-muted">Last Heartbeat:</small><br>
                  <small>{{ formatTime(consumer.lastHeartbeat) }}</small>
                </div>
              </div>

              <div *ngIf="consumer.errorMessage" class="alert alert-danger alert-sm p-2 mb-2">
                <small><i class="fas fa-exclamation-triangle"></i> {{ consumer.errorMessage }}</small>
              </div>

              <div class="mt-auto">
                <small class="text-muted">
                  Created: {{ formatTime(consumer.createdAt) }}
                  <span *ngIf="consumer.startedAt">, Started: {{ formatTime(consumer.startedAt) }}</span>
                </small>
              </div>
            </div>
            <div class="card-footer d-flex justify-content-between">
              <div class="btn-group" role="group">
                <button class="btn btn-sm btn-success" 
                        (click)="startConsumer(consumer.id!)" 
                        [disabled]="isConsumerActive(consumer) || isLoading"
                        *ngIf="isConsumerNotActive(consumer)">
                  <i class="fas fa-play"></i> Start
                </button>
                <button class="btn btn-sm btn-warning" 
                        (click)="stopConsumer(consumer.id!)" 
                        [disabled]="isConsumerNotActive(consumer) || isLoading"
                        *ngIf="isConsumerActive(consumer)">
                  <i class="fas fa-stop"></i> Stop
                </button>
                <button class="btn btn-sm btn-info" 
                        (click)="viewConsumerDetails(consumer)"
                        [disabled]="isLoading">
                  <i class="fas fa-eye"></i> Details
                </button>
              </div>
              <button class="btn btn-sm btn-danger" 
                      (click)="deleteConsumer(consumer.id!)" 
                      [disabled]="isConsumerActive(consumer) || isLoading"
                      onclick="return confirm('Are you sure you want to delete this consumer?')">>
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div class="text-center py-5" *ngIf="consumers.length === 0 && !isLoading">
        <i class="fas fa-download fa-3x text-muted mb-3"></i>
        <h4 class="text-muted">No Consumers Found</h4>
        <p class="text-muted">Create your first Kafka consumer to get started</p>
        <button class="btn btn-primary" (click)="showCreateForm = true">
          <i class="fas fa-plus"></i> Create Consumer
        </button>
      </div>

      <!-- Loading Spinner -->
      <div class="text-center py-5" *ngIf="isLoading">
        <div class="spinner-border text-primary" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted mt-2">Loading consumers...</p>
      </div>

      <!-- Consumer Details Modal -->
      <div class="modal fade show" style="display: block;" *ngIf="selectedConsumer" 
           (click)="closeConsumerDetails($event)">
        <div class="modal-dialog modal-lg">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">
                Consumer Details: {{ selectedConsumer.name }}
                <span class="badge bg-{{ getStatusColor(selectedConsumer.status) }} ms-2">
                  {{ selectedConsumer.status | titlecase }}
                </span>
              </h5>
              <button type="button" class="btn-close" (click)="closeConsumerDetails()"></button>
            </div>
            <div class="modal-body">
              <div class="row mb-3">
                <div class="col-md-6">
                  <h6>Basic Information</h6>
                  <table class="table table-sm">
                    <tr>
                      <td><strong>Name:</strong></td>
                      <td>{{ selectedConsumer.name }}</td>
                    </tr>
                    <tr *ngIf="selectedConsumer.description">
                      <td><strong>Description:</strong></td>
                      <td>{{ selectedConsumer.description }}</td>
                    </tr>
                    <tr>
                      <td><strong>Broker:</strong></td>
                      <td><code>{{ selectedConsumer.broker }}</code></td>
                    </tr>
                    <tr>
                      <td><strong>Group ID:</strong></td>
                      <td><code>{{ selectedConsumer.groupId }}</code></td>
                    </tr>
                    <tr>
                      <td><strong>Status:</strong></td>
                      <td>
                        <span class="badge bg-{{ getStatusColor(selectedConsumer.status) }}">
                          {{ selectedConsumer.status | titlecase }}
                        </span>
                      </td>
                    </tr>
                  </table>
                </div>
                <div class="col-md-6">
                  <h6>Statistics</h6>
                  <table class="table table-sm">
                    <tr>
                      <td><strong>Messages Consumed:</strong></td>
                      <td>{{ selectedConsumer.messageCount | number }}</td>
                    </tr>
                    <tr>
                      <td><strong>Created:</strong></td>
                      <td>{{ formatTime(selectedConsumer.createdAt) }}</td>
                    </tr>
                    <tr>
                      <td><strong>Updated:</strong></td>
                      <td>{{ formatTime(selectedConsumer.updatedAt) }}</td>
                    </tr>
                    <tr *ngIf="selectedConsumer.startedAt">
                      <td><strong>Started:</strong></td>
                      <td>{{ formatTime(selectedConsumer.startedAt) }}</td>
                    </tr>
                    <tr *ngIf="selectedConsumer.lastHeartbeat">
                      <td><strong>Last Heartbeat:</strong></td>
                      <td>{{ formatTime(selectedConsumer.lastHeartbeat) }}</td>
                    </tr>
                  </table>
                </div>
              </div>

              <h6>Topics</h6>
              <div class="mb-3">
                <span class="badge bg-primary me-1" *ngFor="let topic of selectedConsumer.topics">{{ topic }}</span>
              </div>

              <h6>Configuration</h6>
              <table class="table table-sm">
                <tr>
                  <td><strong>Auto Offset Reset:</strong></td>
                  <td><code>{{ selectedConsumer.config.autoOffsetReset }}</code></td>
                </tr>
                <tr>
                  <td><strong>Enable Auto Commit:</strong></td>
                  <td>
                    <span class="badge {{ selectedConsumer.config.enableAutoCommit ? 'bg-success' : 'bg-secondary' }}">
                      {{ selectedConsumer.config.enableAutoCommit ? 'Yes' : 'No' }}
                    </span>
                  </td>
                </tr>
                <tr>
                  <td><strong>Max Poll Records:</strong></td>
                  <td>{{ selectedConsumer.config.maxPollRecords }}</td>
                </tr>
                <tr *ngIf="selectedConsumer.config.sessionTimeoutMs">
                  <td><strong>Session Timeout:</strong></td>
                  <td>{{ selectedConsumer.config.sessionTimeoutMs }}ms</td>
                </tr>
                <tr *ngIf="selectedConsumer.config.heartbeatIntervalMs">
                  <td><strong>Heartbeat Interval:</strong></td>
                  <td>{{ selectedConsumer.config.heartbeatIntervalMs }}ms</td>
                </tr>
              </table>

              <div *ngIf="selectedConsumer.errorMessage" class="alert alert-danger">
                <h6><i class="fas fa-exclamation-triangle"></i> Error Details</h6>
                <p class="mb-0">{{ selectedConsumer.errorMessage }}</p>
              </div>
            </div>
            <div class="modal-footer">
              <button type="button" class="btn btn-secondary" (click)="closeConsumerDetails()">Close</button>
              <button class="btn btn-success" 
                      (click)="startConsumer(selectedConsumer.id!); closeConsumerDetails()" 
                      [disabled]="isConsumerActive(selectedConsumer)"
                      *ngIf="isConsumerNotActive(selectedConsumer)">
                <i class="fas fa-play"></i> Start Consumer
              </button>
              <button class="btn btn-warning" 
                      (click)="stopConsumer(selectedConsumer.id!); closeConsumerDetails()" 
                      [disabled]="isConsumerNotActive(selectedConsumer)"
                      *ngIf="isConsumerActive(selectedConsumer)">
                <i class="fas fa-stop"></i> Stop Consumer
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
    
    .alert-sm { font-size: 0.875rem; }
    
    .modal { background-color: rgba(0,0,0,0.5); }
    
    .table td { border-top: none; padding: 0.25rem 0.5rem; }
    
    .card-footer { background-color: #f8f9fa; }
    
    .badge { font-size: 0.75em; }
    
    code { font-size: 0.875em; }
  `]
})
export class ConsumersComponent implements OnInit, OnDestroy {
  consumers: Consumer[] = [];
  showCreateForm = false;
  isLoading = false;
  errorMessage = '';
  successMessage = '';
  selectedConsumer: Consumer | null = null;
  topicsInput = '';
  
  private refreshSubscription?: Subscription;

  newConsumer: Partial<Consumer> = {
    name: '',
    description: '',
    broker: 'localhost:9092',
    groupId: '',
    topics: [],
    status: 'inactive',
    config: {
      autoOffsetReset: 'latest',
      enableAutoCommit: true,
      maxPollRecords: 500
    } as ConsumerConfig,
    messageCount: 0
  };

  constructor(private kafkaService: KafkaService) {}

  ngOnInit() {
    this.loadConsumers();
    // Auto-refresh every 30 seconds
    this.refreshSubscription = interval(30000).subscribe(() => {
      if (!this.showCreateForm && !this.selectedConsumer) {
        this.loadConsumers(true);
      }
    });
  }

  ngOnDestroy() {
    if (this.refreshSubscription) {
      this.refreshSubscription.unsubscribe();
    }
  }

  loadConsumers(silent = false) {
    if (!silent) this.isLoading = true;
    this.errorMessage = '';

    this.kafkaService.getConsumers().subscribe({
      next: (response: ApiResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.consumers = response.data || [];
        } else {
          this.errorMessage = response.message || 'Failed to load consumers';
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to connect to server. Please check if the backend is running.';
        console.error('Error loading consumers:', error);
      }
    });
  }

  refreshConsumers() {
    this.loadConsumers();
  }

  createConsumer() {
    if (!this.newConsumer.name || !this.newConsumer.broker || !this.newConsumer.groupId || !this.topicsInput) {
      this.errorMessage = 'Please fill in all required fields';
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    // Parse topics from comma-separated input
    this.newConsumer.topics = this.topicsInput.split(',').map(t => t.trim()).filter(t => t);

    this.kafkaService.createConsumer(this.newConsumer).subscribe({
      next: (response: ApiResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.successMessage = 'Consumer created successfully!';
          this.cancelCreate();
          this.loadConsumers();
        } else {
          this.errorMessage = response.message || 'Failed to create consumer';
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to create consumer. Please check the details and try again.';
        console.error('Error creating consumer:', error);
      }
    });
  }

  cancelCreate() {
    this.showCreateForm = false;
    this.newConsumer = {
      name: '',
      description: '',
      broker: 'localhost:9092',
      groupId: '',
      topics: [],
      status: 'inactive',
      config: {
        autoOffsetReset: 'latest',
        enableAutoCommit: true,
        maxPollRecords: 500
      } as ConsumerConfig,
      messageCount: 0
    };
    this.topicsInput = '';
  }

  startConsumer(id: string) {
    this.isLoading = true;
    this.errorMessage = '';

    this.kafkaService.startConsumer(id).subscribe({
      next: (response: ApiResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.successMessage = 'Consumer started successfully!';
          this.loadConsumers();
        } else {
          this.errorMessage = response.message || 'Failed to start consumer';
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to start consumer. Please try again.';
        console.error('Error starting consumer:', error);
      }
    });
  }

  stopConsumer(id: string) {
    this.isLoading = true;
    this.errorMessage = '';

    this.kafkaService.stopConsumer(id).subscribe({
      next: (response: ApiResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.successMessage = 'Consumer stopped successfully!';
          this.loadConsumers();
        } else {
          this.errorMessage = response.message || 'Failed to stop consumer';
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to stop consumer. Please try again.';
        console.error('Error stopping consumer:', error);
      }
    });
  }

  deleteConsumer(id: string) {
    this.isLoading = true;
    this.errorMessage = '';

    this.kafkaService.deleteConsumer(id).subscribe({
      next: (response: ApiResponse) => {
        this.isLoading = false;
        if (response.success) {
          this.successMessage = 'Consumer deleted successfully!';
          this.loadConsumers();
        } else {
          this.errorMessage = response.message || 'Failed to delete consumer';
        }
      },
      error: (error) => {
        this.isLoading = false;
        this.errorMessage = 'Failed to delete consumer. Please try again.';
        console.error('Error deleting consumer:', error);
      }
    });
  }

  viewConsumerDetails(consumer: Consumer) {
    this.selectedConsumer = consumer;
  }

  closeConsumerDetails(event?: MouseEvent) {
    if (event && event.target !== event.currentTarget) {
      return; // Don't close if clicking inside modal
    }
    this.selectedConsumer = null;
  }

  getStatusColor(status: string): string {
    switch (status) {
      case 'active': return 'success';
      case 'error': return 'danger';
      case 'inactive': return 'secondary';
      default: return 'secondary';
    }
  }

  getStatusIcon(status: string): string {
    switch (status) {
      case 'active': return 'play';
      case 'error': return 'exclamation-triangle';
      case 'inactive': return 'stop';
      default: return 'question';
    }
  }

  formatTime(timestamp: string): string {
    return new Date(timestamp).toLocaleString();
  }

  trackByConsumerId(index: number, consumer: Consumer): string {
    return consumer.id || index.toString();
  }

  isConsumerActive(consumer: Consumer): boolean {
    return consumer.status === 'active';
  }

  isConsumerNotActive(consumer: Consumer): boolean {
    return consumer.status !== 'active';
  }
}
