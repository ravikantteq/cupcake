import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { KafkaProducerComponent } from './components/kafka-producer.component';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, KafkaProducerComponent],
  template: `
    <div class="app-container">
      <header>
        <h1>🧁 Cupcake | Backyard</h1>
        <p>Kafka Producer Testing Tool</p>
      </header>
      
      <main>
        <app-kafka-producer></app-kafka-producer>
      </main>
      
      <footer>
        <p>Built with Angular & Go | 
          <a href="http://localhost:8080/swagger/index.html" target="_blank">API Documentation</a>
        </p>
      </footer>
    </div>
  `,
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'cupcake | backyard';
}
