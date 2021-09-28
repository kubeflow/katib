import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ExperimentsComponent } from './pages/experiments/experiments.component';
import { ExperimentDetailsComponent } from './pages/experiment-details/experiment-details.component';
import { ExperimentCreationComponent } from './pages/experiment-creation/experiment-creation.component';

const routes: Routes = [
  { path: '', component: ExperimentsComponent },
  { path: 'experiment/:experimentName', component: ExperimentDetailsComponent },
  { path: 'new', component: ExperimentCreationComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' })],
  exports: [RouterModule],
})
export class AppRoutingModule {}
