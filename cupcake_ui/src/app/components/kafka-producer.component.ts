import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { KafkaService, KafkaMessage } from '../services/kafka.service';

interface MessageHistory {
  id: string;
  timestamp: Date;
  broker: string;
  topic: string;
  key: string;
  value: string;
  success: boolean;
  response?: any;
  expanded?: boolean;
}

@Component({
  selector: 'app-kafka-producer',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="kafka-producer">
      <div class="page-header">
        <h2>� Producer</h2>
        <p>Send messages to Kafka topics with history tracking</p>
      </div>
      
      <form (ngSubmit)="onSubmit()" #kafkaForm="ngForm">
        <div class="form-group">
          <label for="broker">Kafka Broker:</label>
          <input 
            type="text" 
            id="broker" 
            name="broker"
            [(ngModel)]="message.broker" 
            required 
            placeholder="192.168.65.254:9093"
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

      <!-- Message History Section -->
      <div class="message-history" *ngIf="messageHistory.length > 0">
        <div class="history-header">
          <h3>📋 Message History</h3>
          <button class="btn btn-secondary btn-sm" (click)="clearHistory()">Clear History</button>
        </div>
        
        <div class="history-list">
          <div class="history-item" *ngFor="let item of messageHistory; let i = index" 
               [class.expanded]="item.expanded" [class.success]="item.success" [class.error]="!item.success">
            
            <div class="history-summary" (click)="toggleHistoryItem(item)">
              <div class="history-info">
                <span class="history-status">{{item.success ? '✅' : '❌'}}</span>
                <span class="history-topic">{{item.topic}}</span>
                <span class="history-key" *ngIf="item.key">({{item.key}})</span>
                <span class="history-time">{{formatTime(item.timestamp)}}</span>
              </div>
              <div class="history-actions">
                <button class="btn-icon-retry" (click)="retryMessage(item); $event.stopPropagation()" 
                        title="Retry this message">🔄</button>
                <span class="expand-icon">{{item.expanded ? '▼' : '▶'}}</span>
              </div>
            </div>

            <div class="history-details" *ngIf="item.expanded">
              <div class="detail-section">
                <h5>Input:</h5>
                <div class="detail-grid">
                  <div><strong>Broker:</strong> {{item.broker}}</div>
                  <div><strong>Topic:</strong> {{item.topic}}</div>
                  <div><strong>Key:</strong> {{item.key || 'None'}}</div>
                </div>
                <div class="message-value">
                  <strong>Message Value:</strong>
                  <pre>{{item.value}}</pre>
                </div>
              </div>
              
              <div class="detail-section" *ngIf="item.response">
                <h5>Response:</h5>
                <pre>{{item.response | json}}</pre>
              </div>
              
              <div class="detail-section">
                <button class="btn btn-retry" (click)="retryMessage(item)" title="Fill form with this message">
                  🔄 Retry Message
                </button>
              </div>
            </div>
          </div>
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
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .page-header {
      text-align: center;
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 2px solid #f0f0f0;
    }

    .page-header h2 {
      margin: 0;
      color: #2d3748;
      font-size: 2rem;
      font-weight: 600;
    }

    .page-header p {
      margin: 10px 0 0 0;
      color: #718096;
      font-size: 1rem;
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
    
    .message-history {
      margin-top: 30px;
    }
    
    .history-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 2px solid #eee;
    }
    
    .history-header h3 {
      margin: 0;
      color: #333;
    }
    
    .history-list {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
    
    .history-item {
      border: 1px solid #ddd;
      border-radius: 6px;
      background: white;
      overflow: hidden;
    }
    
    .history-item.success {
      border-left: 4px solid #28a745;
    }
    
    .history-item.error {
      border-left: 4px solid #dc3545;
    }
    
    .history-summary {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      cursor: pointer;
      transition: background 0.2s ease;
    }
    
    .history-summary:hover {
      background: #f8f9fa;
    }
    
    .history-info {
      display: flex;
      align-items: center;
      gap: 12px;
    }

    .history-actions {
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .btn-icon-retry {
      background: none;
      border: none;
      padding: 4px 8px;
      border-radius: 4px;
      cursor: pointer;
      font-size: 14px;
      transition: background 0.2s ease;
    }

    .btn-icon-retry:hover {
      background: #e9ecef;
    }
    
    .history-status {
      font-size: 16px;
    }
    
    .history-topic {
      font-weight: 500;
      color: #333;
    }
    
    .history-key {
      color: #666;
      font-size: 14px;
    }
    
    .history-time {
      color: #999;
      font-size: 12px;
    }
    
    .expand-icon {
      color: #666;
      font-size: 14px;
    }
    
    .history-details {
      padding: 16px;
      border-top: 1px solid #eee;
      background: #f8f9fa;
    }
    
    .detail-section {
      margin-bottom: 15px;
    }
    
    .detail-section h5 {
      margin: 0 0 10px 0;
      color: #333;
      font-size: 14px;
      font-weight: 600;
    }
    
    .detail-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 8px;
      margin-bottom: 12px;
      font-size: 14px;
    }
    
    .message-value {
      margin-top: 12px;
    }
    
    .message-value pre {
      background: white;
      border: 1px solid #ddd;
      border-radius: 4px;
      padding: 12px;
      margin: 8px 0 0 0;
      font-size: 13px;
      line-height: 1.4;
      max-height: 200px;
      overflow-y: auto;
    }
    
    .btn-sm {
      padding: 6px 12px;
      font-size: 12px;
    }

    .btn-retry {
      background: #17a2b8;
      color: white;
      border: none;
      padding: 8px 16px;
      border-radius: 4px;
      font-size: 14px;
      cursor: pointer;
      transition: background 0.2s ease;
      margin-top: 10px;
    }

    .btn-retry:hover {
      background: #138496;
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
    broker: '192.168.65.254:9093',
    topic: '',
    key: '',
    value: ''
  };

  result: any = null;
  healthStatus: any = null;
  isLoading = false;
  isHealthChecking = false;
  messageHistory: MessageHistory[] = [];

  constructor(private kafkaService: KafkaService) {
    // Load history from localStorage
    this.loadHistory();
  }

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
        
        // Add to history
        this.addToHistory(this.message, response, true);
      },
      error: (error) => {
        this.result = {
          success: false,
          message: error.error?.message || error.message || 'An error occurred'
        };
        this.isLoading = false;
        
        // Add to history
        this.addToHistory(this.message, this.result, false);
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

  addToHistory(message: KafkaMessage, response: any, success: boolean) {
    const historyItem: MessageHistory = {
      id: Date.now().toString(),
      timestamp: new Date(),
      broker: message.broker,
      topic: message.topic,
      key: message.key,
      value: message.value,
      success: success,
      response: response,
      expanded: false
    };

    this.messageHistory.unshift(historyItem);
    
    // Keep only last 50 messages
    if (this.messageHistory.length > 50) {
      this.messageHistory = this.messageHistory.slice(0, 50);
    }
    
    this.saveHistory();
  }

  toggleHistoryItem(item: MessageHistory) {
    item.expanded = !item.expanded;
  }

  clearHistory() {
    if (confirm('Are you sure you want to clear the message history?')) {
      this.messageHistory = [];
      this.saveHistory();
    }
  }

  retryMessage(historyItem: MessageHistory) {
    // Fill the form with the selected message data
    this.message = {
      broker: historyItem.broker,
      topic: historyItem.topic,
      key: historyItem.key,
      value: historyItem.value
    };
    
    // Scroll to the top to show the form
    window.scrollTo({ top: 0, behavior: 'smooth' });
    
    // Optionally close the history item
    historyItem.expanded = false;
  }

  formatTime(timestamp: Date): string {
    return timestamp.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    });
  }

  private loadHistory() {
    try {
      const stored = localStorage.getItem('kafka-producer-history');
      if (stored) {
        const parsed = JSON.parse(stored);
        this.messageHistory = parsed.map((item: any) => ({
          ...item,
          timestamp: new Date(item.timestamp)
        }));
      }
    } catch (error) {
      console.warn('Failed to load message history:', error);
      this.messageHistory = [];
    }
  }

  private saveHistory() {
    try {
      localStorage.setItem('kafka-producer-history', JSON.stringify(this.messageHistory));
    } catch (error) {
      console.warn('Failed to save message history:', error);
    }
  }
}
