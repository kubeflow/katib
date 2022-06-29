import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { KubeflowModule } from 'kubeflow';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ExperimentsModule } from './pages/experiments/experiments.module';
import { ExperimentDetailsModule } from './pages/experiment-details/experiment-details.module';
import { ExperimentCreationModule } from './pages/experiment-creation/experiment-creation.module';
import { TrialModalModule } from './pages/experiment-details/trials-table/trial-modal/trial-modal.module';

@NgModule({
  declarations: [AppComponent],
  imports: [
    BrowserModule,
    AppRoutingModule,
    KubeflowModule,
    ExperimentsModule,
    ExperimentDetailsModule,
    ReactiveFormsModule,
    ExperimentCreationModule,
    TrialModalModule,
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
