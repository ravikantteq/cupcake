import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

interface TestFlow {
  id: string;
  name: string;
  description: string;
  version: string;
  steps: FlowStep[];
  createdAt: string;
  createdBy: string;
  expanded?: boolean;
  isNew?: boolean;
  hasChanges?: boolean;
  currentExecution?: ExecutionData;
  isExecuting?: boolean;
  executionCompleted?: boolean;
}

interface FlowStep {
  stepId: string;
  type: 'produce' | 'consume' | 'validate' | 'delay';
  config: any;
  isNew?: boolean;
  isEditing?: boolean;
}

interface ExecutionData {
  id: string;
  flowId: string;
  status: 'running' | 'completed' | 'failed' | 'cancelled';
  startTime: string;
  endTime?: string;
  steps: ExecutionStep[];
  metrics: {
    totalDuration: number;
    messagesProduced: number;
    messagesConsumed: number;
    errorsCount: number;
    stepsCompleted: number;
  };
}

interface ExecutionStep {
  stepId: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  input?: any;
  output?: any;
  errors?: string[];
  duration: number;
}

interface NewStep {
  type: 'produce' | 'consume' | 'validate' | 'delay';
  stepId: string;
  config: any;
}

@Component({
  selector: 'app-flows',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="flows-container">
      <div class="flows-header">
        <h1>🔄 Test Flows</h1>
        <div class="header-actions">
          <button class="btn btn-primary" (click)="createNewFlow()">
            ➕ Create New Flow
          </button>
          <button class="btn btn-secondary" (click)="loadFlows()">
            🔄 Refresh
          </button>
        </div>
      </div>

      <div class="flows-content" *ngIf="!loading; else loadingTemplate">
        <div class="flows-list" *ngIf="flows.length > 0; else emptyTemplate">
          <div class="flow-item" *ngFor="let flow of flows" [class.expanded]="flow.expanded">
            
            <!-- Flow Header -->
            <div class="flow-header" (click)="toggleFlow(flow)">
              <div class="flow-info">
                <h3 *ngIf="!flow.isNew">{{ flow.name }}</h3>
                <input 
                  *ngIf="flow.isNew" 
                  [(ngModel)]="flow.name" 
                  placeholder="Enter flow name"
                  class="flow-name-input"
                  (click)="$event.stopPropagation()">
                <span class="flow-description">{{ flow.description || 'No description' }}</span>
                <span class="flow-meta">
                  {{ flow.steps.length }} steps • v{{ flow.version }} • {{ formatDate(flow.createdAt) }}
                </span>
              </div>
              <div class="flow-actions" (click)="$event.stopPropagation()">
                <button class="btn btn-sm btn-success" (click)="executeFlow(flow.id)" *ngIf="!flow.isNew && !flow.isExecuting && !flow.executionCompleted">
                  ▶️ Execute
                </button>
                <button class="btn btn-sm btn-warning" (click)="clearExecution(flow)" *ngIf="!flow.isNew && flow.executionCompleted">
                  🔄 Clear Results
                </button>
                <span class="execution-status" *ngIf="flow.isExecuting">
                  ⏳ Executing...
                </span>
                <button class="btn btn-sm btn-primary" (click)="saveFlow(flow)" *ngIf="flow.hasChanges || flow.isNew">
                  💾 Save
                </button>
                <span class="expand-icon">{{ flow.expanded ? '▼' : '▶' }}</span>
              </div>
            </div>

            <!-- Flow Details (Expanded) -->
            <div class="flow-details" *ngIf="flow.expanded">
              
              <!-- Flow Metadata -->
              <div class="flow-metadata" *ngIf="!flow.isNew">
                <div class="metadata-row">
                  <label>Description:</label>
                  <input [(ngModel)]="flow.description" (ngModelChange)="markAsChanged(flow)" placeholder="Enter description">
                </div>
                <div class="metadata-row">
                  <label>Version:</label>
                  <input [(ngModel)]="flow.version" (ngModelChange)="markAsChanged(flow)" placeholder="1.0.0">
                </div>
              </div>

              <div class="flow-metadata" *ngIf="flow.isNew">
                <div class="metadata-row">
                  <label>Description:</label>
                  <input [(ngModel)]="flow.description" placeholder="Enter description">
                </div>
                <div class="metadata-row">
                  <label>Version:</label>
                  <input [(ngModel)]="flow.version" placeholder="1.0.0" value="1.0.0">
                </div>
              </div>

              <!-- Steps List -->
              <div class="steps-section">
                <div class="steps-header">
                  <h4>Flow Steps</h4>
                  <div class="step-actions">
                    <select [(ngModel)]="newStepType" class="step-type-select">
                      <option value="">Select step type...</option>
                      <option value="produce">📤 Produce Message</option>
                      <option value="consume">📥 Consume Message</option>
                      <option value="validate">✅ Validate Response</option>
                      <option value="delay">⏱️ Delay</option>
                    </select>
                    <button class="btn btn-sm btn-primary" (click)="addStep(flow)" [disabled]="!newStepType">
                      ➕ Add Step
                    </button>
                  </div>
                </div>

                <div class="steps-list">
                  <div class="step-item" *ngFor="let step of flow.steps; let i = index" 
                       [class.editing]="step.isEditing" [class.is-new]="step.isNew">
                    
                    <div class="step-row">
                      <!-- Step Display -->
                      <div class="step-display" *ngIf="!step.isEditing">
                        <div class="step-number">{{ i + 1 }}</div>
                        <div class="step-content">
                          <div class="step-header">
                            <span class="step-type">{{ getStepTypeDisplay(step.type) }}</span>
                            <span class="step-id">{{ step.stepId }}</span>
                          </div>
                          <div class="step-config">{{ getStepConfigSummary(step) }}</div>
                        </div>
                        <div class="step-actions">
                          <button class="btn btn-xs btn-secondary" (click)="moveStepUp(flow, i)" [disabled]="i === 0">⬆️</button>
                          <button class="btn btn-xs btn-secondary" (click)="moveStepDown(flow, i)" [disabled]="i === flow.steps.length - 1">⬇️</button>
                          <button class="btn btn-xs btn-secondary" (click)="editStep(step)">✏️</button>
                          <button class="btn btn-xs btn-danger" (click)="removeStep(flow, i)">🗑️</button>
                        </div>
                      </div>

                      <!-- Execution Timeline Column -->
                      <div class="execution-timeline-column" *ngIf="flow.isExecuting || flow.executionCompleted">
                        <div class="timeline-connector" *ngIf="i > 0"></div>
                        <div class="timeline-status" [ngClass]="getStepStatusClass(flow, step.stepId)">
                          <span>{{ getStepStatusIcon(flow, step.stepId) }}</span>
                        </div>
                        <div class="timeline-connector" *ngIf="i < flow.steps.length - 1"></div>
                      </div>
                    </div>

                    <!-- Step Execution Response -->
                    <div class="step-response" *ngIf="(flow.isExecuting || flow.executionCompleted) && hasStepData(flow.currentExecution, step.stepId)">
                      <div class="response-header">
                        <span class="response-title">📊 Step Response</span>
                        <span class="response-duration">{{ getStepDuration(flow.currentExecution, step.stepId) }}ms</span>
                      </div>
                      <div class="response-content">
                        <pre>{{ getStepResponse(flow.currentExecution, step.stepId) }}</pre>
                      </div>
                    </div>

                    <!-- Step Edit Form -->
                    <div class="step-edit-form" *ngIf="step.isEditing">
                      <div class="step-edit-header">
                        <h5>{{ getStepTypeDisplay(step.type) }} - {{ step.stepId }}</h5>
                        <div class="edit-actions">
                          <button class="btn btn-xs btn-primary" (click)="saveStep(flow, step)">💾 Save</button>
                          <button class="btn btn-xs btn-secondary" (click)="cancelEditStep(step)">❌ Cancel</button>
                        </div>
                      </div>
                      
                      <!-- Produce Step Form -->
                      <div *ngIf="step.type === 'produce'" class="step-config-form">
                        <div class="form-row">
                          <label>Kafka Broker:</label>
                          <input [(ngModel)]="step.config.broker" placeholder="192.168.65.254:9093">
                        </div>
                        <div class="form-row">
                          <label>Topic:</label>
                          <input [(ngModel)]="step.config.topic" placeholder="kafka-topic-name">
                        </div>
                        <div class="form-row">
                          <label>Message Key:</label>
                          <input [(ngModel)]="step.config.message.key" placeholder="message-key">
                        </div>
                        <div class="form-row">
                          <label>Message Value:</label>
                          <textarea [(ngModel)]="step.config.message.value" placeholder='{"data": "your message content", "timestamp": "2025-07-28"}' rows="5"></textarea>
                        </div>
                        <div class="form-row">
                          <label>Timeout (ms):</label>
                          <input type="number" [(ngModel)]="step.config.timeout" placeholder="5000">
                        </div>
                      </div>

                      <!-- Consume Step Form -->
                      <div *ngIf="step.type === 'consume'" class="step-config-form">
                        <div class="form-row">
                          <label>Kafka Broker:</label>
                          <input [(ngModel)]="step.config.broker" placeholder="192.168.65.254:9093">
                        </div>
                        <div class="form-row">
                          <label>Topic:</label>
                          <input [(ngModel)]="step.config.topic" placeholder="kafka-topic-name">
                        </div>
                        <div class="form-row">
                          <label>Expected Count:</label>
                          <input type="number" [(ngModel)]="step.config.expectedCount" placeholder="1">
                        </div>
                        <div class="form-row">
                          <label>Timeout (ms):</label>
                          <input type="number" [(ngModel)]="step.config.timeout" placeholder="10000">
                        </div>
                      </div>

                      <!-- Delay Step Form -->
                      <div *ngIf="step.type === 'delay'" class="step-config-form">
                        <div class="form-row">
                          <label>Delay (ms):</label>
                          <input type="number" [(ngModel)]="step.config.delayMs" placeholder="1000">
                        </div>
                      </div>

                      <!-- Validate Step Form -->
                      <div *ngIf="step.type === 'validate'" class="step-config-form">
                        <div class="form-row">
                          <label>Expected Message (JSON):</label>
                          <textarea [(ngModel)]="step.config.expectedMessage" placeholder='{"status": "processed", "orderId": "12345", "amount": 100}' rows="6"></textarea>
                        </div>
                      </div>
                    </div>
                  </div>

                  <div class="empty-steps" *ngIf="flow.steps.length === 0">
                    <p>No steps yet. Add your first step using the dropdown above.</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <ng-template #emptyTemplate>
          <div class="empty-state">
            <div class="empty-icon">🔄</div>
            <h2>No Test Flows Yet</h2>
            <p>Create your first test flow to get started</p>
            <button class="btn btn-primary" (click)="createNewFlow()">
              Create Your First Flow
            </button>
          </div>
        </ng-template>
      </div>

      <ng-template #loadingTemplate>
        <div class="loading-state">
          <div class="loading-spinner"></div>
          <p>Loading test flows...</p>
        </div>
      </ng-template>
    </div>
  `,
  styles: [`
    .flows-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .flows-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 2px solid #f0f0f0;
    }

    .flows-header h1 {
      margin: 0;
      color: #2d3748;
      font-size: 2.5rem;
      font-weight: 700;
    }

    .header-actions {
      display: flex;
      gap: 15px;
    }

    .flows-list {
      display: flex;
      flex-direction: column;
      gap: 15px;
    }

    .flow-item {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 12px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
      transition: all 0.2s ease;
    }

    .flow-item.expanded {
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
    }

    .flow-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 20px 25px;
      cursor: pointer;
      transition: background 0.2s ease;
    }

    .flow-header:hover {
      background: #f7fafc;
    }

    .flow-info {
      flex: 1;
    }

    .flow-info h3 {
      margin: 0 0 5px 0;
      color: #2d3748;
      font-size: 1.25rem;
      font-weight: 600;
    }

    .flow-name-input {
      font-size: 1.25rem;
      font-weight: 600;
      border: 1px solid #e2e8f0;
      border-radius: 4px;
      padding: 5px 10px;
      margin: 0 0 5px 0;
      width: 300px;
    }

    .flow-description {
      display: block;
      color: #718096;
      font-size: 0.9rem;
      margin-bottom: 5px;
    }

    .flow-meta {
      font-size: 0.8rem;
      color: #a0aec0;
    }

    .flow-actions {
      display: flex;
      align-items: center;
      gap: 10px;
    }

    .expand-icon {
      font-size: 1.2rem;
      color: #718096;
      margin-left: 10px;
    }

    .flow-details {
      padding: 0 25px 25px 25px;
      border-top: 1px solid #f0f0f0;
      background: #fafafa;
    }

    .flow-metadata {
      margin-bottom: 25px;
      padding-top: 20px;
    }

    .metadata-row {
      display: flex;
      align-items: center;
      gap: 15px;
      margin-bottom: 15px;
    }

    .metadata-row label {
      font-weight: 500;
      color: #4a5568;
      min-width: 100px;
    }

    .metadata-row input {
      flex: 1;
      max-width: 300px;
      padding: 8px 12px;
      border: 1px solid #e2e8f0;
      border-radius: 4px;
      font-size: 0.9rem;
    }

    .steps-section {
      background: white;
      border-radius: 8px;
      padding: 20px;
      border: 1px solid #e2e8f0;
    }

    .steps-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
    }

    .steps-header h4 {
      margin: 0;
      color: #2d3748;
      font-size: 1.1rem;
      font-weight: 600;
    }

    .step-actions {
      display: flex;
      gap: 10px;
      align-items: center;
    }

    .step-type-select {
      padding: 6px 10px;
      border: 1px solid #e2e8f0;
      border-radius: 4px;
      font-size: 0.9rem;
      min-width: 200px;
    }

    .steps-list {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }

    .step-item {
      border: 1px solid #e2e8f0;
      border-radius: 6px;
      background: white;
      margin-bottom: 10px;
    }

    .step-item.editing {
      border-color: #4299e1;
      box-shadow: 0 0 0 1px #4299e1;
    }

    .step-item.is-new {
      border-color: #48bb78;
      background: #f0fff4;
    }

    .step-row {
      display: flex;
      align-items: stretch;
    }

    .step-display {
      display: flex;
      align-items: center;
      padding: 15px;
      flex: 1;
    }

    /* Execution Timeline Column */
    .execution-timeline-column {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      width: 60px;
      background: #f8f9fa;
      border-left: 1px solid #e2e8f0;
      position: relative;
    }

    .timeline-connector {
      width: 2px;
      height: 20px;
      background: #4299e1;
      flex: 1;
    }

    .timeline-status {
      background: white;
      border: 2px solid #e2e8f0;
      border-radius: 50%;
      width: 32px;
      height: 32px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 14px;
      z-index: 1;
      margin: 8px 0;
      transition: all 0.3s ease;
    }

    .timeline-status.status-running {
      border-color: #f6ad55;
      background: #fef5e7;
      animation: pulse 1.5s infinite;
    }

    .timeline-status.status-completed {
      border-color: #48bb78;
      background: #f0fff4;
    }

    .timeline-status.status-failed {
      border-color: #f56565;
      background: #fed7d7;
    }

    .timeline-status.status-pending {
      border-color: #e2e8f0;
      background: white;
    }

    @keyframes pulse {
      0%, 100% { transform: scale(1); }
      50% { transform: scale(1.1); }
    }

    /* Step Response Styles */
    .step-response {
      margin: 0 15px 15px 15px;
      padding: 15px;
      background: #f8f9fa;
      border: 1px solid #e9ecef;
      border-radius: 6px;
      border-left: 4px solid #4299e1;
    }

    .response-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;
      padding-bottom: 8px;
      border-bottom: 1px solid #e9ecef;
    }

    .response-title {
      font-weight: 600;
      color: #2d3748;
      font-size: 0.9rem;
    }

    .response-duration {
      font-size: 0.8rem;
      color: #718096;
      background: #e2e8f0;
      padding: 2px 8px;
      border-radius: 10px;
    }

    .response-content {
      font-size: 0.85rem;
    }

    .response-content pre {
      background: white;
      border: 1px solid #e2e8f0;
      border-radius: 4px;
      padding: 10px;
      margin: 0;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 0.8rem;
      line-height: 1.4;
      overflow-x: auto;
      white-space: pre-wrap;
      word-wrap: break-word;
    }

    .execution-status {
      color: #f6ad55;
      font-weight: 600;
      font-size: 0.9rem;
      display: flex;
      align-items: center;
      gap: 5px;
    }

    .btn-warning {
      background: #f6ad55;
      color: white;
    }

    .btn-warning:hover {
      background: #ed8936;
    }

    .step-number {
      background: #4299e1;
      color: white;
      width: 30px;
      height: 30px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.8rem;
      font-weight: 600;
      margin-right: 15px;
    }

    .step-content {
      flex: 1;
    }

    .step-header {
      display: flex;
      gap: 15px;
      align-items: center;
      margin-bottom: 5px;
    }

    .step-type {
      font-weight: 500;
      color: #2d3748;
      font-size: 0.9rem;
    }

    .step-id {
      color: #718096;
      font-size: 0.8rem;
      background: #f7fafc;
      padding: 2px 6px;
      border-radius: 3px;
    }

    .step-config {
      color: #718096;
      font-size: 0.8rem;
      line-height: 1.4;
    }

    .step-actions {
      display: flex;
      gap: 5px;
    }

    .step-edit-form {
      padding: 20px;
      background: #f7fafc;
      border-top: 1px solid #e2e8f0;
    }

    .step-edit-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 20px;
    }

    .step-edit-header h5 {
      margin: 0;
      color: #2d3748;
      font-size: 1rem;
      font-weight: 600;
    }

    .edit-actions {
      display: flex;
      gap: 8px;
    }

    .step-config-form {
      display: flex;
      flex-direction: column;
      gap: 15px;
    }

    .form-row {
      display: flex;
      flex-direction: column;
      gap: 5px;
    }

    .form-row label {
      font-weight: 500;
      color: #4a5568;
      font-size: 0.9rem;
    }

    .form-row input,
    .form-row textarea {
      padding: 8px 12px;
      border: 1px solid #e2e8f0;
      border-radius: 4px;
      font-size: 0.9rem;
      font-family: inherit;
    }

    .form-row textarea {
      resize: vertical;
      min-height: 60px;
      font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
      font-size: 0.85rem;
      line-height: 1.4;
      white-space: pre;
      overflow-wrap: break-word;
    }

    .empty-steps {
      text-align: center;
      padding: 40px 20px;
      color: #718096;
      font-style: italic;
    }

    .btn {
      padding: 8px 16px;
      border: none;
      border-radius: 6px;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      text-decoration: none;
      transition: all 0.2s ease;
      display: inline-flex;
      align-items: center;
      gap: 5px;
    }

    .btn-xs {
      padding: 4px 6px;
      font-size: 0.75rem;
    }

    .btn-sm {
      padding: 6px 12px;
      font-size: 0.8rem;
    }

    .btn-primary {
      background: #4299e1;
      color: white;
    }

    .btn-primary:hover {
      background: #3182ce;
    }

    .btn-secondary {
      background: #e2e8f0;
      color: #4a5568;
    }

    .btn-secondary:hover {
      background: #cbd5e0;
    }

    .btn-success {
      background: #48bb78;
      color: white;
    }

    .btn-success:hover {
      background: #38a169;
    }

    .btn-danger {
      background: #f56565;
      color: white;
    }

    .btn-danger:hover {
      background: #e53e3e;
    }

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .empty-state, .loading-state {
      text-align: center;
      padding: 60px 20px;
      background: white;
      border-radius: 12px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
    }

    .empty-icon {
      font-size: 4rem;
      margin-bottom: 20px;
    }

    .empty-state h2 {
      margin: 0 0 10px 0;
      color: #2d3748;
      font-size: 1.5rem;
    }

    .empty-state p {
      margin: 0 0 25px 0;
      color: #718096;
      font-size: 1.1rem;
    }

    .loading-spinner {
      width: 40px;
      height: 40px;
      border: 4px solid #e2e8f0;
      border-top: 4px solid #4299e1;
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin: 0 auto 20px;
    }

    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }

    @media (max-width: 768px) {
      .flows-header {
        flex-direction: column;
        gap: 15px;
      }

      .header-actions {
        width: 100%;
        justify-content: center;
      }

      .flow-header {
        flex-direction: column;
        gap: 15px;
        align-items: flex-start;
      }

      .flow-actions {
        width: 100%;
        justify-content: flex-end;
      }

      .steps-header {
        flex-direction: column;
        gap: 15px;
        align-items: stretch;
      }

      .step-actions {
        flex-direction: column;
      }

      .step-type-select {
        min-width: auto;
      }
    }
  `]
})
export class FlowsComponent implements OnInit, OnDestroy {
  flows: TestFlow[] = [];
  loading = true;
  newStepType: string = '';
  executionPollingInterval: any;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadFlows();
  }

  ngOnDestroy() {
    if (this.executionPollingInterval) {
      clearInterval(this.executionPollingInterval);
    }
  }

  loadFlows() {
    this.loading = true;
    this.http.get<any>(`${environment.apiBaseUrl}/api/v1/flows`).subscribe({
      next: (response: any) => {
        this.flows = (response.data || response || []).map((flow: any) => ({
          ...flow,
          expanded: false,
          isNew: false,
          hasChanges: false
        }));
        this.loading = false;
      },
      error: (error: any) => {
        console.error('Failed to load flows:', error);
        this.flows = [];
        this.loading = false;
      }
    });
  }

  createNewFlow() {
    const newFlow: TestFlow = {
      id: '',
      name: '',
      description: '',
      version: '1.0.0',
      steps: [],
      createdAt: new Date().toISOString(),
      createdBy: 'system',
      expanded: true,
      isNew: true,
      hasChanges: false
    };
    this.flows.unshift(newFlow);
  }

  toggleFlow(flow: TestFlow) {
    if (!flow.isNew) {
      flow.expanded = !flow.expanded;
    }
  }

  markAsChanged(flow: TestFlow) {
    flow.hasChanges = true;
  }

  addStep(flow: TestFlow) {
    if (!this.newStepType) return;

    const stepId = `step-${flow.steps.length + 1}`;
    const newStep: FlowStep = {
      stepId: stepId,
      type: this.newStepType as any,
      config: this.getDefaultStepConfig(this.newStepType),
      isNew: true,
      isEditing: true
    };

    flow.steps.push(newStep);
    this.markAsChanged(flow);
    this.newStepType = '';
  }

  getDefaultStepConfig(type: string): any {
    const defaults: { [key: string]: any } = {
      produce: {
        broker: '192.168.65.254:9093',
        topic: '',
        message: {
          key: '',
          value: '{"data": "your message content", "timestamp": "2025-07-28"}'
        },
        timeout: 5000
      },
      consume: {
        broker: '192.168.65.254:9093',
        topic: '',
        expectedCount: 1,
        timeout: 10000
      },
      validate: {
        expectedMessage: '{"status": "processed", "orderId": "12345", "amount": 100}'
      },
      delay: {
        delayMs: 1000
      }
    };
    return defaults[type] || {};
  }

  editStep(step: FlowStep) {
    step.isEditing = true;
  }

  saveStep(flow: TestFlow, step: FlowStep) {
    step.isEditing = false;
    step.isNew = false;
    this.markAsChanged(flow);
  }

  cancelEditStep(step: FlowStep) {
    step.isEditing = false;
    if (step.isNew) {
      // Remove the step if it was new and cancelled
      const flow = this.flows.find(f => f.steps.includes(step));
      if (flow) {
        const index = flow.steps.indexOf(step);
        flow.steps.splice(index, 1);
      }
    }
  }

  removeStep(flow: TestFlow, index: number) {
    if (confirm('Are you sure you want to remove this step?')) {
      flow.steps.splice(index, 1);
      this.markAsChanged(flow);
      
      // Regenerate step IDs
      flow.steps.forEach((step, i) => {
        step.stepId = `step-${i + 1}`;
      });
    }
  }

  moveStepUp(flow: TestFlow, index: number) {
    if (index > 0) {
      [flow.steps[index], flow.steps[index - 1]] = [flow.steps[index - 1], flow.steps[index]];
      this.markAsChanged(flow);
      this.regenerateStepIds(flow);
    }
  }

  moveStepDown(flow: TestFlow, index: number) {
    if (index < flow.steps.length - 1) {
      [flow.steps[index], flow.steps[index + 1]] = [flow.steps[index + 1], flow.steps[index]];
      this.markAsChanged(flow);
      this.regenerateStepIds(flow);
    }
  }

  regenerateStepIds(flow: TestFlow) {
    flow.steps.forEach((step, i) => {
      step.stepId = `step-${i + 1}`;
    });
  }

  saveFlow(flow: TestFlow) {
    if (!flow.name) {
      alert('Please enter a flow name');
      return;
    }

    const flowData = {
      name: flow.name,
      description: flow.description,
      version: flow.version,
      steps: flow.steps.map(step => ({
        stepId: step.stepId,
        type: step.type,
        config: step.config
      }))
    };

    const request = flow.isNew 
      ? this.http.post<any>(`${environment.apiBaseUrl}/api/v1/flows`, flowData)
      : this.http.put<any>(`${environment.apiBaseUrl}/api/v1/flows/${flow.id}`, flowData);

    request.subscribe({
      next: (response: any) => {
        if (flow.isNew) {
          flow.id = response.data?.id || response.id;
          flow.isNew = false;
        }
        flow.hasChanges = false;
        alert(`Flow ${flow.isNew ? 'created' : 'updated'} successfully!`);
      },
      error: (error: any) => {
        console.error('Failed to save flow:', error);
        alert('Failed to save flow. Please check the console for details.');
      }
    });
  }

  executeFlow(flowId: string) {
    if (!flowId) return;
    
    // Find the flow and set it as executing
    const flow = this.flows.find(f => f.id === flowId);
    if (!flow) return;
    
    flow.isExecuting = true;
    flow.executionCompleted = false;
    flow.expanded = true; // Auto-expand to show execution progress
    
    this.http.post<any>(`${environment.apiBaseUrl}/api/v1/flows/${flowId}/execute`, {}).subscribe({
      next: (response: any) => {
        const executionId = response.data?.id || response.id;
        if (executionId) {
          this.startExecutionPolling(flow, executionId);
        }
      },
      error: (error: any) => {
        console.error('Failed to execute flow:', error);
        alert('Failed to execute flow. Please check the console for details.');
        flow.isExecuting = false;
      }
    });
  }

  clearExecution(flow: TestFlow) {
    flow.isExecuting = false;
    flow.executionCompleted = false;
    flow.currentExecution = undefined;
    
    // Stop polling if still running
    if (this.executionPollingInterval) {
      clearInterval(this.executionPollingInterval);
      this.executionPollingInterval = null;
    }
  }

  startExecutionPolling(flow: TestFlow, executionId: string) {
    // Poll for execution updates every 1 second
    this.executionPollingInterval = setInterval(() => {
      this.http.get<any>(`${environment.apiBaseUrl}/api/v1/executions/${executionId}`).subscribe({
        next: (response: any) => {
          const execution = response.data || response;
          flow.currentExecution = execution;
          
          // Stop polling if execution is completed
          if (execution.status === 'completed' || execution.status === 'failed' || execution.status === 'cancelled') {
            flow.isExecuting = false;
            flow.executionCompleted = true;
            if (this.executionPollingInterval) {
              clearInterval(this.executionPollingInterval);
              this.executionPollingInterval = null;
            }
          }
        },
        error: (error: any) => {
          console.error('Failed to fetch execution status:', error);
        }
      });
    }, 1000);
  }

  getStepExecutionStatus(execution: ExecutionData, stepId: string): string {
    const executionStep = execution.steps.find(s => s.stepId === stepId);
    return executionStep?.status || 'pending';
  }

  getExecutionStepData(execution: ExecutionData, stepId: string): ExecutionStep | null {
    return execution.steps.find(s => s.stepId === stepId) || null;
  }

  hasStepData(execution: ExecutionData | undefined, stepId: string): boolean {
    if (!execution) return false;
    return execution.steps.some(s => s.stepId === stepId);
  }

  getStepDuration(execution: ExecutionData | undefined, stepId: string): number {
    if (!execution) return 0;
    const step = execution.steps.find(s => s.stepId === stepId);
    return step?.duration || 0;
  }

  getStepResponse(execution: ExecutionData | undefined, stepId: string): string {
    if (!execution) return '';
    const step = execution.steps.find(s => s.stepId === stepId);
    if (!step) return '';
    
    if (step.errors && step.errors.length > 0) {
      return `❌ ERRORS:\n${step.errors.join('\n')}`;
    }
    
    if (step.output) {
      return JSON.stringify(step.output, null, 2);
    }
    
    return 'No response data';
  }

  getStepStatusClass(flow: TestFlow, stepId: string): string {
    if (!flow.currentExecution) return 'status-pending';
    const step = flow.currentExecution.steps.find(s => s.stepId === stepId);
    return `status-${step?.status || 'pending'}`;
  }

  getStepStatusIcon(flow: TestFlow, stepId: string): string {
    if (!flow.currentExecution) return '⚪';
    const step = flow.currentExecution.steps.find(s => s.stepId === stepId);
    switch (step?.status) {
      case 'running': return '⏳';
      case 'completed': return '✅';
      case 'failed': return '❌';
      default: return '⚪';
    }
  }

  getStepTypeDisplay(type: string): string {
    const types: { [key: string]: string } = {
      produce: 'Produce Message',
      consume: 'Consume Message',
      validate: 'Validate Response',
      delay: 'Delay Execution'
    };
    return types[type] || type;
  }

  getStepConfigSummary(step: FlowStep): string {
    switch (step.type) {
      case 'produce':
        return `Topic: ${step.config.topic || 'Not set'} | Key: ${step.config.message?.key || 'None'}`;
      case 'consume':
        return `Topic: ${step.config.topic || 'Not set'} | Expected: ${step.config.expectedCount || 1} messages`;
      case 'delay':
        return `Delay: ${step.config.delayMs || 0}ms`;
      case 'validate':
        return `Validation rules configured`;
      default:
        return 'Configuration set';
    }
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
}
