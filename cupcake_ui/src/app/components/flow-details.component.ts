import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

interface TestFlow {
  id?: string;
  name: string;
  description: string;
  version: string;
  steps: FlowStep[];
  createdAt?: string;
  createdBy: string;
}

interface FlowStep {
  stepId: string;
  type: 'produce' | 'consume' | 'validate' | 'delay';
  config: StepConfig;
}

interface StepConfig {
  topic?: string;
  message?: { [key: string]: any };
  expectedMessage?: { [key: string]: any };
  timeout?: number;
  retries?: number;
  expectedCount?: number;
  delayMs?: number;
}

@Component({
  selector: 'app-flow-details',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  template: `
    <div class="flow-details-container">
      <!-- Header -->
      <div class="flow-header">
        <div class="header-content">
          <button class="btn btn-secondary" routerLink="/flows">
            ← Back to Flows
          </button>
          <div class="flow-title">
            <h1>{{ isNewFlow ? '➕ Create New Flow' : '📋 Edit Flow' }}</h1>
            <p>{{ isNewFlow ? 'Design a new Kafka test flow' : 'Modify your test flow' }}</p>
          </div>
        </div>
        <div class="header-actions">
          <button class="btn btn-secondary" (click)="executeFlow()" *ngIf="!isNewFlow" [disabled]="flow.steps.length === 0">
            ▶️ Execute Flow
          </button>
          <button class="btn btn-primary" (click)="saveFlow()" [disabled]="!isFormValid()">
            � {{ isNewFlow ? 'Create Flow' : 'Save Changes' }}
          </button>
        </div>
      </div>

      <!-- Flow Basic Info -->
      <div class="flow-info-section">
        <h2>📝 Flow Information</h2>
        <div class="form-grid">
          <div class="form-group">
            <label for="name">Flow Name *</label>
            <input
              id="name"
              type="text"
              [(ngModel)]="flow.name"
              placeholder="Enter flow name"
              class="form-input"
              required
            >
          </div>
          <div class="form-group">
            <label for="version">Version</label>
            <input
              id="version"
              type="text"
              [(ngModel)]="flow.version"
              placeholder="1.0.0"
              class="form-input"
            >
          </div>
          <div class="form-group full-width">
            <label for="description">Description</label>
            <textarea
              id="description"
              [(ngModel)]="flow.description"
              placeholder="Describe what this flow tests"
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>
          <div class="form-group">
            <label for="createdBy">Created By</label>
            <input
              id="createdBy"
              type="text"
              [(ngModel)]="flow.createdBy"
              placeholder="Your name"
              class="form-input"
            >
          </div>
        </div>
      </div>

      <!-- Steps Section -->
      <div class="steps-section">
        <div class="steps-header">
          <h2>🔧 Flow Steps</h2>
          <div class="steps-actions">
            <button class="btn btn-primary" (click)="addProduceStep()">
              ➕ Add Producer Step
            </button>
            <button class="btn btn-secondary" (click)="addDelayStep()">
              ⏱️ Add Delay Step
            </button>
          </div>
        </div>

        <div class="steps-list" *ngIf="flow.steps.length > 0; else noStepsTemplate">
          <div class="step-card" *ngFor="let step of flow.steps; let i = index">
            <div class="step-header">
              <div class="step-info">
                <div class="step-number">{{ i + 1 }}</div>
                <div class="step-title">
                  <h4>{{ getStepTypeDisplay(step.type) }}</h4>
                  <span class="step-id">{{ step.stepId }}</span>
                </div>
              </div>
              <div class="step-actions">
                <button class="btn btn-sm btn-danger" (click)="removeStep(i)">
                  🗑️ Remove
                </button>
                <button class="btn btn-sm btn-secondary" (click)="moveStepUp(i)" [disabled]="i === 0">
                  ↑
                </button>
                <button class="btn btn-sm btn-secondary" (click)="moveStepDown(i)" [disabled]="i === flow.steps.length - 1">
                  ↓
                </button>
              </div>
            </div>

            <!-- Producer Step Configuration -->
            <div class="step-config" *ngIf="step.type === 'produce'">
              <div class="config-grid">
                <div class="form-group">
                  <label>Step ID</label>
                  <input
                    type="text"
                    [(ngModel)]="step.stepId"
                    placeholder="step-1"
                    class="form-input"
                  >
                </div>
                <div class="form-group">
                  <label>Topic</label>
                  <input
                    type="text"
                    [(ngModel)]="getProduceConfig(step).topic"
                    placeholder="test-topic"
                    class="form-input"
                  >
                </div>
                <div class="form-group">
                  <label>Message Key</label>
                  <input
                    type="text"
                    [(ngModel)]="getProduceConfig(step).message!['key']"
                    placeholder="order-123"
                    class="form-input"
                  >
                </div>
                <div class="form-group">
                  <label>Message Value</label>
                  <input
                    type="text"
                    [(ngModel)]="getProduceConfig(step).message!['value']"
                    placeholder="Hello World"
                    class="form-input"
                  >
                </div>
                <div class="form-group">
                  <label>Timeout (ms)</label>
                  <input
                    type="number"
                    [(ngModel)]="getProduceConfig(step).timeout"
                    placeholder="5000"
                    class="form-input"
                  >
                </div>
              </div>
            </div>

            <!-- Delay Step Configuration -->
            <div class="step-config" *ngIf="step.type === 'delay'">
              <div class="config-grid">
                <div class="form-group">
                  <label>Step ID</label>
                  <input
                    type="text"
                    [(ngModel)]="step.stepId"
                    placeholder="delay-1"
                    class="form-input"
                  >
                </div>
                <div class="form-group">
                  <label>Duration (milliseconds)</label>
                  <input
                    type="number"
                    [(ngModel)]="getDelayConfig(step).delayMs"
                    placeholder="5000"
                    class="form-input"
                    min="100"
                  >
                </div>
              </div>
            </div>
          </div>
        </div>

        <ng-template #noStepsTemplate>
          <div class="no-steps">
            <div class="no-steps-icon">🔧</div>
            <h3>No Steps Yet</h3>
            <p>Add producer steps to define what your flow should do</p>
            <button class="btn btn-primary" (click)="addProduceStep()">
              Add Your First Producer Step
            </button>
          </div>
        </ng-template>
      </div>
    </div>
  `,
  styles: [`
    .flow-details-container {
      max-width: 1000px;
      margin: 0 auto;
      padding: 20px;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .flow-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 2px solid #f0f0f0;
    }

    .header-content {
      display: flex;
      align-items: center;
      gap: 20px;
    }

    .flow-title h1 {
      margin: 0;
      color: #2d3748;
      font-size: 2rem;
      font-weight: 700;
    }

    .flow-title p {
      margin: 5px 0 0 0;
      color: #718096;
      font-size: 1rem;
    }

    .header-actions {
      display: flex;
      gap: 15px;
    }

    .flow-info-section, .steps-section {
      background: white;
      padding: 25px;
      border-radius: 12px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
      margin-bottom: 25px;
    }

    .flow-info-section h2, .steps-section h2 {
      margin: 0 0 20px 0;
      color: #2d3748;
      font-size: 1.5rem;
      font-weight: 600;
    }

    .form-grid, .config-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 20px;
    }

    .form-group {
      display: flex;
      flex-direction: column;
      gap: 5px;
    }

    .form-group.full-width {
      grid-column: 1 / -1;
    }

    .form-group label {
      font-weight: 500;
      color: #2d3748;
      font-size: 0.875rem;
    }

    .form-input, .form-textarea {
      padding: 10px 12px;
      border: 1px solid #e2e8f0;
      border-radius: 6px;
      font-size: 0.875rem;
      transition: border-color 0.2s ease;
    }

    .form-input:focus, .form-textarea:focus {
      outline: none;
      border-color: #4299e1;
      box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
    }

    .steps-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 25px;
    }

    .steps-actions {
      display: flex;
      gap: 10px;
    }

    .steps-list {
      display: flex;
      flex-direction: column;
      gap: 20px;
    }

    .step-card {
      background: #f7fafc;
      border: 1px solid #e2e8f0;
      border-radius: 8px;
      padding: 20px;
    }

    .step-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 15px;
    }

    .step-info {
      display: flex;
      align-items: center;
      gap: 15px;
    }

    .step-number {
      background: #4299e1;
      color: white;
      width: 32px;
      height: 32px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      font-size: 0.875rem;
    }

    .step-title h4 {
      margin: 0;
      color: #2d3748;
      font-size: 1.125rem;
      font-weight: 600;
    }

    .step-id {
      color: #718096;
      font-size: 0.75rem;
    }

    .step-actions {
      display: flex;
      gap: 8px;
    }

    .step-config {
      background: white;
      padding: 15px;
      border-radius: 6px;
      border: 1px solid #e2e8f0;
    }

    .no-steps {
      text-align: center;
      padding: 40px 20px;
      color: #718096;
    }

    .no-steps-icon {
      font-size: 3rem;
      margin-bottom: 15px;
    }

    .no-steps h3 {
      margin: 0 0 10px 0;
      color: #2d3748;
      font-size: 1.25rem;
    }

    .no-steps p {
      margin: 0 0 20px 0;
      font-size: 1rem;
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

    .btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }

    .btn-sm {
      padding: 4px 8px;
      font-size: 0.75rem;
    }

    .btn-primary {
      background: #4299e1;
      color: white;
    }

    .btn-primary:hover:not(:disabled) {
      background: #3182ce;
    }

    .btn-secondary {
      background: #e2e8f0;
      color: #4a5568;
    }

    .btn-secondary:hover:not(:disabled) {
      background: #cbd5e0;
    }

    .btn-danger {
      background: #f56565;
      color: white;
    }

    .btn-danger:hover:not(:disabled) {
      background: #e53e3e;
    }

    @media (max-width: 768px) {
      .flow-header {
        flex-direction: column;
        gap: 20px;
      }

      .header-content {
        flex-direction: column;
        align-items: flex-start;
        gap: 10px;
      }

      .header-actions {
        width: 100%;
        justify-content: center;
      }

      .steps-header {
        flex-direction: column;
        gap: 15px;
      }

      .steps-actions {
        width: 100%;
        justify-content: center;
      }

      .step-header {
        flex-direction: column;
        gap: 15px;
        align-items: flex-start;
      }

      .step-actions {
        width: 100%;
        justify-content: flex-end;
      }
    }
  `]
})
export class FlowDetailsComponent implements OnInit {
  flow: TestFlow = {
    name: '',
    description: '',
    version: '1.0.0',
    steps: [],
    createdBy: ''
  };

  isNewFlow = true;
  loading = false;
  lastExecution: any = null;

  constructor(
    private http: HttpClient,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit() {
    const flowId = this.route.snapshot.params['id'];
    if (flowId && flowId !== 'new') {
      this.isNewFlow = false;
      this.loadFlow(flowId);
    }
  }

  loadFlow(flowId: string) {
    this.loading = true;
    this.http.get<any>(`${environment.apiBaseUrl}/api/v1/flows/${flowId}`).subscribe({
      next: (response: any) => {
        this.flow = response.data || response;
        this.loading = false;
      },
      error: (error: any) => {
        console.error('Failed to load flow:', error);
        this.loading = false;
        alert('Failed to load flow');
        this.router.navigate(['/flows']);
      }
    });
  }

  saveFlow() {
    if (!this.isFormValid()) {
      alert('Please fill in all required fields: Flow Name, Created By, and at least one step with complete data');
      return;
    }

    // Clean and prepare the flow data to match backend expectations
    const flowData = {
      name: this.flow.name.trim(),
      description: this.flow.description.trim(),
      version: this.flow.version.trim(),
      createdBy: this.flow.createdBy.trim(),
      steps: this.flow.steps.map(step => {
        if (step.type === 'produce') {
          return {
            stepId: step.stepId,
            type: step.type,
            config: {
              topic: step.config.topic,
              message: {
                key: step.config.message!['key'],
                value: step.config.message!['value']
              },
              timeout: step.config.timeout || 5000
            }
          };
        } else if (step.type === 'delay') {
          return {
            stepId: step.stepId,
            type: step.type,
            config: {
              delayMs: step.config.delayMs
            }
          };
        }
        return step;
      })
    };

    console.log('Saving flow data:', JSON.stringify(flowData, null, 2));

    const url = this.isNewFlow 
      ? `${environment.apiBaseUrl}/api/v1/flows`
      : `${environment.apiBaseUrl}/api/v1/flows/${this.flow.id}`;
    
    const request = this.isNewFlow 
      ? this.http.post<any>(url, flowData)
      : this.http.put<any>(url, flowData);

    request.subscribe({
      next: (response: any) => {
        console.log('Flow saved successfully:', response);
        alert(this.isNewFlow ? 'Flow created successfully!' : 'Flow updated successfully!');
        if (this.isNewFlow) {
          const flowId = response.data?.id || response.id;
          this.router.navigate(['/flows', flowId]);
        }
      },
      error: (error: any) => {
        console.error('Failed to save flow:', error);
        console.error('Error details:', error.error);
        const errorMessage = error.error?.message || error.message || 'Unknown error';
        alert(`Failed to save flow: ${errorMessage}`);
      }
    });
  }

  executeFlow() {
    if (!this.flow.id) return;

    this.http.post<any>(`${environment.apiBaseUrl}/api/v1/flows/${this.flow.id}/execute`, {}).subscribe({
      next: (response: any) => {
        this.lastExecution = response.data || response;
        alert(`Flow execution completed! Status: ${this.lastExecution.status}`);
      },
      error: (error: any) => {
        console.error('Failed to execute flow:', error);
        alert('Failed to execute flow. Please check the console for details.');
      }
    });
  }

  addProduceStep() {
    const stepNumber = this.flow.steps.length + 1;
    const newStep: FlowStep = {
      stepId: `step-${stepNumber}`,
      type: 'produce',
      config: {
        topic: '',
        message: {
          key: '',
          value: ''
        },
        timeout: 5000
      }
    };
    this.flow.steps.push(newStep);
    console.log('Added producer step:', newStep);
  }

  addDelayStep() {
    const stepNumber = this.flow.steps.length + 1;
    const newStep: FlowStep = {
      stepId: `delay-${stepNumber}`,
      type: 'delay',
      config: {
        delayMs: 5000
      }
    };
    this.flow.steps.push(newStep);
  }

  removeStep(index: number) {
    if (confirm('Are you sure you want to remove this step?')) {
      this.flow.steps.splice(index, 1);
    }
  }

  moveStepUp(index: number) {
    if (index > 0) {
      const step = this.flow.steps[index];
      this.flow.steps[index] = this.flow.steps[index - 1];
      this.flow.steps[index - 1] = step;
    }
  }

  moveStepDown(index: number) {
    if (index < this.flow.steps.length - 1) {
      const step = this.flow.steps[index];
      this.flow.steps[index] = this.flow.steps[index + 1];
      this.flow.steps[index + 1] = step;
    }
  }

  getProduceConfig(step: FlowStep): StepConfig {
    return step.config;
  }

  getDelayConfig(step: FlowStep): StepConfig {
    return step.config;
  }

  isFormValid(): boolean {
    if (this.flow.name.trim() === '' || 
        this.flow.createdBy.trim() === '' ||
        this.flow.steps.length === 0) {
      return false;
    }

    // Check if all steps have required data
    for (const step of this.flow.steps) {
      if (step.type === 'produce') {
        if (!step.config.topic || 
            !step.config.message ||
            !step.config.message['key'] || 
            !step.config.message['value']) {
          return false;
        }
      } else if (step.type === 'delay') {
        if (!step.config.delayMs || step.config.delayMs < 100) {
          return false;
        }
      }
    }

    return true;
  }

  getStepTypeDisplay(type: string): string {
    const types: { [key: string]: string } = {
      produce: 'Producer Step',
      consume: 'Consumer Step',
      validate: 'Validation Step',
      delay: 'Delay Step'
    };
    return types[type] || type;
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
}
