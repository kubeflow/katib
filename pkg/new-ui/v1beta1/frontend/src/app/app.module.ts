import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { KubeflowModule } from 'kubeflow';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ExperimentsModule } from './pages/experiments/experiments.module';
import { ExperimentDetailsModule } from './pages/experiment-details/experiment-details.module';
import { ExperimentCreationModule } from './pages/experiment-creation/experiment-creation.module';

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
  ],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
