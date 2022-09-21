import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ExperimentsComponent } from './pages/experiments/experiments.component';
import { ExperimentDetailsComponent } from './pages/experiment-details/experiment-details.component';
import { ExperimentCreationComponent } from './pages/experiment-creation/experiment-creation.component';
import { TrialModalComponent } from './pages/experiment-details/trials-table/trial-modal/trial-modal.component';

const routes: Routes = [
  { path: '', component: ExperimentsComponent },
  { path: 'experiment/:experimentName', component: ExperimentDetailsComponent },
  { path: 'new', component: ExperimentCreationComponent },
  {
    path: 'experiment/:experimentName/trial/:trialName',
    component: TrialModalComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { relativeLinkResolution: 'legacy' })],
  exports: [RouterModule],
})
export class AppRoutingModule {}
