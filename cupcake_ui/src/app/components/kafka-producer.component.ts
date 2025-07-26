import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { KafkaService, KafkaMessage } from '../services/kafka.service';

@Component({
  selector: 'app-kafka-producer',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="kafka-producer">
      <h2>Kafka Producer</h2>
      
      <form (ngSubmit)="onSubmit()" #kafkaForm="ngForm">
        <div class="form-group">
          <label for="broker">Kafka Broker:</label>
          <input 
            type="text" 
            id="broker" 
            name="broker"
            [(ngModel)]="message.broker" 
            required 
            placeholder="localhost:9093"
            class="form-control">
        </div>

        <div class="form-group">
          <label for="topic">Topic:</label>
          <input 
            type="text" 
            id="topic" 
            name="topic"
            [(ngModel)]="message.topic" 
            required 
            placeholder="test-topic"
            class="form-control">
        </div>

        <div class="form-group">
          <label for="key">Key (optional):</label>
          <input 
            type="text" 
            id="key" 
            name="key"
            [(ngModel)]="message.key" 
            placeholder="message-key"
            class="form-control">
        </div>

        <div class="form-group">
          <label for="value">Message Value:</label>
          <textarea 
            id="value" 
            name="value"
            [(ngModel)]="message.value" 
            required 
            placeholder='{"data": "your message here"}'
            rows="5"
            class="form-control"></textarea>
        </div>

        <button 
          type="submit" 
          [disabled]="!kafkaForm.form.valid || isLoading"
          class="btn btn-primary">
          {{isLoading ? 'Publishing...' : 'Publish Message'}}
        </button>
      </form>

      <div *ngIf="result" class="result">
        <h3>Result:</h3>
        <div [ngClass]="result.success ? 'success' : 'error'">
          <p><strong>{{result.success ? 'Success' : 'Error'}}:</strong> {{result.message}}</p>
          <pre *ngIf="result.data">{{result.data | json}}</pre>
        </div>
      </div>

      <div class="health-check">
        <button 
          (click)="checkHealth()" 
          [disabled]="isHealthChecking"
          class="btn btn-secondary">
          {{isHealthChecking ? 'Checking...' : 'Check Backend Health'}}
        </button>
        
        <div *ngIf="healthStatus" class="health-status">
          <span [ngClass]="healthStatus.success ? 'healthy' : 'unhealthy'">
            Backend: {{healthStatus.success ? 'Healthy' : 'Unhealthy'}}
          </span>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .kafka-producer {
      max-width: 600px;
      margin: 20px auto;
      padding: 20px;
      border: 1px solid #ddd;
      border-radius: 8px;
      background-color: #f9f9f9;
    }

    .form-group {
      margin-bottom: 15px;
    }

    label {
      display: block;
      margin-bottom: 5px;
      font-weight: bold;
    }

    .form-control {
      width: 100%;
      padding: 8px;
      border: 1px solid #ccc;
      border-radius: 4px;
      font-size: 14px;
    }

    .btn {
      padding: 10px 20px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 14px;
      margin-right: 10px;
    }

    .btn:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }

    .btn-primary {
      background-color: #007bff;
      color: white;
    }

    .btn-secondary {
      background-color: #6c757d;
      color: white;
    }

    .result {
      margin-top: 20px;
      padding: 15px;
      border-radius: 4px;
    }

    .success {
      background-color: #d4edda;
      border: 1px solid #c3e6cb;
      color: #155724;
    }

    .error {
      background-color: #f8d7da;
      border: 1px solid #f5c6cb;
      color: #721c24;
    }

    .health-check {
      margin-top: 20px;
      padding-top: 20px;
      border-top: 1px solid #ddd;
    }

    .health-status {
      margin-top: 10px;
    }

    .healthy {
      color: #28a745;
      font-weight: bold;
    }

    .unhealthy {
      color: #dc3545;
      font-weight: bold;
    }

    pre {
      background-color: #f8f9fa;
      padding: 10px;
      border-radius: 4px;
      overflow-x: auto;
    }
  `]
})
export class KafkaProducerComponent {
  message: KafkaMessage = {
    broker: 'localhost:9093',
    topic: '',
    key: '',
    value: ''
  };

  result: any = null;
  healthStatus: any = null;
  isLoading = false;
  isHealthChecking = false;

  constructor(private kafkaService: KafkaService) {}

  onSubmit() {
    if (!this.message.topic || !this.message.value) {
      return;
    }

    this.isLoading = true;
    this.result = null;

    this.kafkaService.publishMessage(this.message).subscribe({
      next: (response) => {
        this.result = response;
        this.isLoading = false;
      },
      error: (error) => {
        this.result = {
          success: false,
          message: error.error?.message || error.message || 'An error occurred'
        };
        this.isLoading = false;
      }
    });
  }

  checkHealth() {
    this.isHealthChecking = true;
    this.healthStatus = null;

    this.kafkaService.healthCheck().subscribe({
      next: (response) => {
        this.healthStatus = response;
        this.isHealthChecking = false;
      },
      error: (error) => {
        this.healthStatus = {
          success: false,
          message: 'Backend unavailable'
        };
        this.isHealthChecking = false;
      }
    });
  }
}
