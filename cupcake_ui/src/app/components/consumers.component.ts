import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-consumers',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div class="consumers-container">
      <div class="consumers-header">
        <h1>👂 Consumer Management</h1>
        <p>Setup and monitor Kafka consumers automatically</p>
      </div>

      <div class="placeholder">
        <div class="icon">🎧</div>
        <h2>Consumer Management</h2>
        <p>Automatic consumer setup and monitoring coming soon!</p>
      </div>
    </div>
  `,
  styles: [`
    .consumers-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .consumers-header {
      margin-bottom: 30px;
      text-align: center;
    }

    .consumers-header h1 {
      margin: 0 0 10px 0;
      color: #2d3748;
      font-size: 2.5rem;
      font-weight: 700;
    }

    .placeholder {
      background: white;
      padding: 60px 40px;
      border-radius: 12px;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
      border: 1px solid #e2e8f0;
      text-align: center;
    }

    .icon {
      font-size: 4rem;
      margin-bottom: 20px;
    }
  `]
})
export class ConsumersComponent {}
