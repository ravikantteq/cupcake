import { Routes } from '@angular/router';
import { DashboardComponent } from './components/dashboard.component';
import { KafkaProducerComponent } from './components/kafka-producer.component';
import { FlowsComponent } from './components/flows.component';
import { SuitesComponent } from './components/suites.component';
import { ConsumersComponent } from './components/consumers.component';
import { ExecutionsComponent } from './components/executions.component';

export const routes: Routes = [
  { path: '', component: DashboardComponent },
  { path: 'dashboard', component: DashboardComponent },
  { path: 'producer', component: KafkaProducerComponent },
  { path: 'legacy', redirectTo: 'producer' }, // Backward compatibility
  { path: 'flows', component: FlowsComponent },
  { path: 'suites', component: SuitesComponent },
  { path: 'consumers', component: ConsumersComponent },
  { path: 'executions', component: ExecutionsComponent },
  { path: '**', redirectTo: '' }
];
