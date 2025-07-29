import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../environments/environment';

interface HealthStatus {
  status: string;
  timestamp: string;
  services: { [key: string]: string };
  version: string;
}

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <header class="app-header">
      <div class="header-content">
        <div class="header-left" (click)="goHome()">
          <h1 class="app-title">🧁 Cupcake</h1>
          <p class="app-subtitle">Kafka Testing Platform v2.0</p>
        </div>
        <div class="header-right">
          <div class="health-status" [class.healthy]="isHealthy" [class.unhealthy]="!isHealthy">
            <span class="status-dot"></span>
            {{ healthStatus?.status || 'checking...' }}
          </div>
        </div>
      </div>
    </header>
  `,
  styles: [`
    .app-header {
      background: white;
      border-bottom: 2px solid #f0f0f0;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      position: sticky;
      top: 0;
      z-index: 100;
      width: 100%;
    }

    .header-content {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .header-left {
      cursor: pointer;
      transition: opacity 0.2s ease;
    }

    .header-left:hover {
      opacity: 0.8;
    }

    .app-title {
      margin: 0;
      color: #2d3748;
      font-size: 2.5rem;
      font-weight: 700;
    }

    .app-subtitle {
      margin: 5px 0 0 0;
      color: #718096;
      font-size: 1.1rem;
    }

    .header-right {
      display: flex;
      align-items: center;
      gap: 20px;
    }

    .health-status {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 16px;
      border-radius: 20px;
      font-weight: 500;
      text-transform: uppercase;
      font-size: 0.875rem;
    }

    .health-status.healthy {
      background: #c6f6d5;
      color: #22543d;
    }

    .health-status.unhealthy {
      background: #fed7d7;
      color: #742a2a;
    }

    .status-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: currentColor;
    }

    @media (max-width: 768px) {
      .header-content {
        flex-direction: column;
        gap: 15px;
        text-align: center;
      }

      .app-title {
        font-size: 2rem;
      }

      .app-subtitle {
        font-size: 1rem;
      }
    }
  `]
})
export class HeaderComponent implements OnInit {
  healthStatus: HealthStatus | null = null;
  isHealthy = false;

  constructor(private http: HttpClient, private router: Router) {}

  ngOnInit() {
    this.loadHealthStatus();
    // Refresh health status every 30 seconds
    setInterval(() => this.loadHealthStatus(), 30000);
  }

  goHome() {
    this.router.navigate(['/']);
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
}
