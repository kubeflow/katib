import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ReactiveFormsModule } from '@angular/forms';
import { KubeflowModule } from 'kubeflow';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { ExperimentsModule } from './pages/experiments/experiments.module';
import { ExperimentDetailsModule } from './pages/experiment-details/experiment-details.module';
import { ExperimentCreationModule } from './pages/experiment-creation/experiment-creation.module';
import {
  MatSnackBarConfig,
  MAT_SNACK_BAR_DEFAULT_OPTIONS,
} from '@angular/material/snack-bar';

/**
 * MAT_SNACK_BAR_DEFAULT_OPTIONS values can be found
 * here
 * https://github.com/angular/components/blob/main/src/material/snack-bar/snack-bar-config.ts#L25-L58
 */
const KwaSnackBarConfig: MatSnackBarConfig = {
  duration: 4000,
};

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
  providers: [
    { provide: MAT_SNACK_BAR_DEFAULT_OPTIONS, useValue: KwaSnackBarConfig },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
