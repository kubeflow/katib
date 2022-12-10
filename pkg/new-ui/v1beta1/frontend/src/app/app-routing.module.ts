import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ExperimentsComponent } from './pages/experiments/experiments.component';
import { ExperimentDetailsComponent } from './pages/experiment-details/experiment-details.component';
import { ExperimentCreationComponent } from './pages/experiment-creation/experiment-creation.component';
import { TrialDetailsComponent } from './pages/experiment-details/trials-table/trial-details/trial-details.component';

const routes: Routes = [
  { path: '', component: ExperimentsComponent },
  {
    path: 'experiment/:namespace/:experimentName',
    component: ExperimentDetailsComponent,
  },
  { path: 'new', component: ExperimentCreationComponent },
  {
    path: 'experiment/:experimentName/trial/:trialName',
    component: TrialDetailsComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' })],
  exports: [RouterModule],
})
export class AppRoutingModule {}
