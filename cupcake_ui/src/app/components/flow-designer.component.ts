import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-flow-designer',
  standalone: true,
  imports: [CommonModule, RouterModule],
  template: `
    <div class="designer-container">
      <div class="designer-header">
        <h1>🎨 Flow Designer</h1>
        <p>Create a new test flow with visual designer</p>
        <button class="btn btn-secondary" routerLink="/flows">← Back to Flows</button>
      </div>

      <div class="designer-content">
        <div class="coming-soon">
          <div class="icon">🚧</div>
          <h2>Visual Flow Designer</h2>
          <p>The drag-and-drop flow designer is coming soon!</p>
          
          <div class="preview-features">
            <h3>Planned Features:</h3>
            <ul>
              <li>🔄 Drag-and-drop step placement</li>
              <li>📝 Visual step configuration</li>
              <li>🔗 Auto-connection between steps</li>
              <li>✅ Real-time validation</li>
              <li>📋 Flow templates library</li>
            </ul>
          </div>

          <div class="current-options">
            <h3>Current Options:</h3>
            <div class="options-grid">
              <div class="option-card">
                <h4>🔧 Use API</h4>
                <p>Create flows programmatically using our REST API</p>
                <a href="http://localhost:8080/swagger/index.html" target="_blank" class="btn btn-primary">
                  View API Docs
                </a>
              </div>
              <div class="option-card">
                <h4>⚡ Quick Test</h4>
                <p>Use the legacy producer for simple message testing</p>
                <button class="btn btn-secondary" routerLink="/legacy">
                  Go to Producer
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .designer-container {
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
    }

    .designer-header {
      margin-bottom: 30px;
      text-align: center;
    }

    .designer-header h1 {
      margin: 0 0 10px 0;
      color: #2d3748;
      font-size: 2.5rem;
      font-weight: 700;
    }

    .designer-header p {
      margin: 0 0 20px 0;
      color: #718096;
      font-size: 1.1rem;
    }

    .coming-soon {
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

    .coming-soon h2 {
      margin: 0 0 15px 0;
      color: #2d3748;
      font-size: 2rem;
    }

    .coming-soon > p {
      margin: 0 0 40px 0;
      color: #718096;
      font-size: 1.2rem;
    }

    .preview-features, .current-options {
      margin: 40px 0;
      text-align: left;
    }

    .preview-features h3, .current-options h3 {
      margin: 0 0 20px 0;
      color: #2d3748;
      font-size: 1.25rem;
      text-align: center;
    }

    .preview-features ul {
      max-width: 400px;
      margin: 0 auto;
      text-align: left;
      color: #718096;
      line-height: 1.8;
    }

    .options-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 20px;
      margin-top: 20px;
    }

    .option-card {
      background: #f7fafc;
      padding: 25px;
      border-radius: 8px;
      border: 1px solid #e2e8f0;
      text-align: center;
    }

    .option-card h4 {
      margin: 0 0 10px 0;
      color: #2d3748;
      font-size: 1.125rem;
    }

    .option-card p {
      margin: 0 0 20px 0;
      color: #718096;
      line-height: 1.5;
    }

    .btn {
      padding: 10px 20px;
      border: none;
      border-radius: 6px;
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      text-decoration: none;
      transition: all 0.2s ease;
      display: inline-block;
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
  `]
})
export class FlowDesignerComponent {}
