import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

interface TestFlow {
  id: string;
  name: string;
  description: string;
  version: string;
  steps: any[];
  createdAt: string;
  createdBy: string;
}

interface HealthStatus {
  status: string;
  timestamp: string;
  services: { [key: string]: string };
  version: string;
}

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div class="dashboard-container">
      <!-- Quick Stats -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-number">{{ flows.length }}</div>
          <div class="stat-label">Test Flows</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ getTotalSteps() }}</div>
          <div class="stat-label">Total Steps</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">0</div>
          <div class="stat-label">Active Executions</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ getServicesCount() }}</div>
          <div class="stat-label">Healthy Services</div>
        </div>
      </div>

      <!-- Main Navigation -->
      <div class="nav-grid">
        <div class="nav-card" routerLink="/producer">
          <div class="nav-icon">⚡</div>
          <h3>Producer</h3>
          <p>Kafka message producer with history tracking</p>
          <div class="nav-action">Quick Test →</div>
        </div>
        <div class="nav-card" routerLink="/flows">
          <div class="nav-icon">🔄</div>
          <h3>Test Flows</h3>
          <p>Create and manage test flows with intelligent validation</p>
          <div class="nav-action">View Flows →</div>
        </div>

        <div class="nav-card" routerLink="/suites">
          <div class="nav-icon">📋</div>
          <h3>Test Suites</h3>
          <p>Organize flows into comprehensive test suites</p>
          <div class="nav-action">Manage Suites →</div>
        </div>

        <div class="nav-card" routerLink="/consumers">
          <div class="nav-icon">👂</div>
          <h3>Consumers</h3>
          <p>Setup and monitor Kafka consumers automatically</p>
          <div class="nav-action">View Consumers →</div>
        </div>

        <div class="nav-card" routerLink="/executions">
          <div class="nav-icon">📊</div>
          <h3>Executions</h3>
          <p>Monitor test execution results and metrics</p>
          <div class="nav-action">View Results →</div>
        </div>

        <div class="nav-card" (click)="openApiDocs()">
          <div class="nav-icon">📖</div>
          <h3>API Documentation</h3>
          <p>Interactive API documentation and testing</p>
          <div class="nav-action">View Docs →</div>
        </div>
      </div>

      <!-- Recent Flows -->
      <div class="recent-section" *ngIf="flows.length > 0">
        <h2>Recent Test Flows</h2>
        <div class="flows-list">
          <div class="flow-card" *ngFor="let flow of flows.slice(0, 3)">
            <div class="flow-header">
              <h4>{{ flow.name }}</h4>
              <span class="flow-version">v{{ flow.version }}</span>
            </div>
            <p class="flow-description">{{ flow.description }}</p>
            <div class="flow-meta">
              <span class="flow-steps">{{ flow.steps.length }} steps</span>
              <span class="flow-created">{{ formatDate(flow.createdAt) }}</span>
            </div>
            <div class="flow-actions">
              <button class="btn btn-secondary" (click)="executeFlow(flow.id)">
                ▶️ Execute
              </button>
              <button class="btn btn-primary" (click)="viewFlowDetails(flow.id)">
                View Details
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Getting Started -->
      <div class="getting-started" *ngIf="flows.length === 0">
        <h2>🚀 Getting Started</h2>
        <div class="steps-list">
          <div class="step">
            <div class="step-number">1</div>
            <div class="step-content">
              <h4>Create Your First Test Flow</h4>
              <p>Design a multi-step test scenario with produce, consume, and validate steps</p>
              <button class="btn btn-primary" routerLink="/flows/new">Create Flow</button>
            </div>
          </div>
          <div class="step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h4>Set Up Consumers</h4>
              <p>Configure Kafka consumers to automatically listen for test responses</p>
              <button class="btn btn-secondary" routerLink="/consumers">Setup Consumers</button>
            </div>
          </div>
          <div class="step">
            <div class="step-number">3</div>
            <div class="step-content">
              <h4>Execute & Monitor</h4>
              <p>Run your tests and monitor execution results in real-time</p>
              <button class="btn btn-secondary" routerLink="/executions">View Executions</button>
            </div>
          </div>
        </div>
      </div>

      <!-- System Info -->
      <div class="system-info" *ngIf="healthStatus">
        <h3>System Status</h3>
        <div class="services-grid">
          <div class="service-item" *ngFor="let service of getServices()">
            <span class="service-name">{{ service.name }}</span>
            <span class="service-status" [class.healthy]="service.healthy">
              {{ service.status }}
            </span>
          </div>
        </div>
        <div class="system-meta">
          <span>Version: {{ healthStatus.version }}</span>
          <span>Last Check: {{ formatDate(healthStatus.timestamp) }}</span>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .dashboard-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 20px;
      margin-bottom: 40px;
    }

    .stat-card {
      background: white;
      padding: 20px;
      border-radius: 12px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
      text-align: center;
    }

    .stat-number {
      font-size: 2.5rem;
      font-weight: 700;
      color: #2d3748;
      margin-bottom: 5px;
    }

    .stat-label {
      color: #718096;
      font-size: 0.875rem;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }

    .nav-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 20px;
      margin-bottom: 40px;
    }

    .nav-card {
      background: white;
      padding: 30px;
      border-radius: 12px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
      cursor: pointer;
      transition: all 0.2s ease;
      text-decoration: none;
      color: inherit;
    }

    .nav-card:hover {
      transform: translateY(-2px);
      box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
      border-color: #4299e1;
    }

    .nav-icon {
      font-size: 2.5rem;
      margin-bottom: 15px;
    }

    .nav-card h3 {
      margin: 0 0 10px 0;
      color: #2d3748;
      font-size: 1.25rem;
      font-weight: 600;
    }

    .nav-card p {
      margin: 0 0 15px 0;
      color: #718096;
      line-height: 1.5;
    }

    .nav-action {
      color: #4299e1;
      font-weight: 500;
      font-size: 0.875rem;
    }

    .recent-section, .getting-started, .system-info {
      background: white;
      padding: 30px;
      border-radius: 12px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
      margin-bottom: 30px;
    }

    .recent-section h2, .getting-started h2, .system-info h3 {
      margin: 0 0 20px 0;
      color: #2d3748;
      font-size: 1.5rem;
      font-weight: 600;
    }

    .flows-list {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 20px;
    }

    .flow-card {
      background: #f7fafc;
      padding: 20px;
      border-radius: 8px;
      border: 1px solid #e2e8f0;
    }

    .flow-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 10px;
    }

    .flow-header h4 {
      margin: 0;
      color: #2d3748;
      font-size: 1.125rem;
    }

    .flow-version {
      background: #4299e1;
      color: white;
      padding: 2px 8px;
      border-radius: 12px;
      font-size: 0.75rem;
      font-weight: 500;
    }

    .flow-description {
      color: #718096;
      margin: 0 0 15px 0;
      line-height: 1.4;
    }

    .flow-meta {
      display: flex;
      gap: 15px;
      margin-bottom: 15px;
      font-size: 0.875rem;
      color: #718096;
    }

    .flow-actions {
      display: flex;
      gap: 10px;
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

    .steps-list {
      display: flex;
      flex-direction: column;
      gap: 20px;
    }

    .step {
      display: flex;
      gap: 20px;
      align-items: flex-start;
    }

    .step-number {
      background: #4299e1;
      color: white;
      width: 40px;
      height: 40px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      flex-shrink: 0;
    }

    .step-content h4 {
      margin: 0 0 8px 0;
      color: #2d3748;
      font-size: 1.125rem;
    }

    .step-content p {
      margin: 0 0 15px 0;
      color: #718096;
      line-height: 1.5;
    }

    .services-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 15px;
      margin-bottom: 20px;
    }

    .service-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px;
      background: #f7fafc;
      border-radius: 6px;
    }

    .service-name {
      font-weight: 500;
      color: #2d3748;
    }

    .service-status {
      font-size: 0.875rem;
      padding: 2px 8px;
      border-radius: 12px;
      background: #fed7d7;
      color: #742a2a;
    }

    .service-status.healthy {
      background: #c6f6d5;
      color: #22543d;
    }

    .system-meta {
      display: flex;
      gap: 20px;
      font-size: 0.875rem;
      color: #718096;
      padding-top: 15px;
      border-top: 1px solid #e2e8f0;
    }

    .services-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
      gap: 15px;
      margin-bottom: 20px;
    }

    .service-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px;
      background: #f7fafc;
      border-radius: 6px;
    }

    .service-name {
      font-weight: 500;
      color: #2d3748;
    }

    .service-status {
      font-size: 0.875rem;
      padding: 2px 8px;
      border-radius: 12px;
      background: #fed7d7;
      color: #742a2a;
    }

    .service-status.healthy {
      background: #c6f6d5;
      color: #22543d;
    }

    @media (max-width: 768px) {
      .dashboard-header {
        flex-direction: column;
        gap: 15px;
        text-align: center;
      }

      .stats-grid {
        grid-template-columns: repeat(2, 1fr);
      }

      .nav-grid {
        grid-template-columns: 1fr;
      }

      .step {
        flex-direction: column;
        text-align: center;
      }
    }
  `]
})
export class DashboardComponent implements OnInit {
  flows: TestFlow[] = [];
  healthStatus: HealthStatus | null = null;
  isHealthy = false;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadFlows();
    this.loadHealthStatus();
  }

  loadHealthStatus() {
    this.http.get<any>(`${environment.apiBaseUrl}/health`).subscribe({
      next: (response) => {
        this.healthStatus = response.data || response;
        this.isHealthy = this.healthStatus?.status === 'healthy';
      },
      error: (error) => {
        console.error('Failed to load health status:', error);
        this.isHealthy = false;
      }
    });
  }

  loadFlows() {
    this.http.get<any>(`${environment.apiBaseUrl}/api/v1/flows`).subscribe({
      next: (response) => {
        this.flows = response.data || response || [];
      },
      error: (error) => {
        console.error('Failed to load flows:', error);
        this.flows = [];
      }
    });
  }

  executeFlow(flowId: string) {
    this.http.post<any>(`${environment.apiBaseUrl}/api/v1/flows/${flowId}/execute`, {}).subscribe({
      next: (response) => {
        alert(`Flow execution started! Execution ID: ${response.data?.id || 'unknown'}`);
      },
      error: (error) => {
        console.error('Failed to execute flow:', error);
        alert('Failed to execute flow. Please check the console for details.');
      }
    });
  }

  openApiDocs() {
    // Open API documentation in a new window
    window.open(`${environment.apiBaseUrl}/swagger/index.html`, '_blank');
  }

  viewFlowDetails(flowId: string) {
    // Navigate to flows page and show the specific flow details
    window.location.href = `/flows#${flowId}`;
  }

  getTotalSteps(): number {
    return this.flows.reduce((sum, f) => sum + (f.steps?.length || 0), 0);
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  getServices() {
    if (!this.healthStatus?.services) return [];
    
    return Object.entries(this.healthStatus.services).map(([name, status]) => ({
      name: name.charAt(0).toUpperCase() + name.slice(1),
      status: status as string,
      healthy: !status.toString().includes('unhealthy') && !status.toString().includes('error')
    }));
  }

  getServicesCount(): number {
    const services = this.getServices();
    return services.filter(s => s.healthy).length;
  }
}
