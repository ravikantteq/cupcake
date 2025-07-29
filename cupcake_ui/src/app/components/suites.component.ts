import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-suites',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div class="suites-container">
      <div class="suites-header">
        <h1>📋 Test Suites</h1>
        <p>Organize and execute multiple test flows together</p>
      </div>

      <div class="placeholder">
        <div class="icon">📦</div>
        <h2>Test Suites Management</h2>
        <p>Suite creation and execution coming soon!</p>
      </div>
    </div>
  `,
  styles: [`
    .suites-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .suites-header {
      margin-bottom: 30px;
      text-align: center;
    }

    .suites-header h1 {
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
export class SuitesComponent {}
