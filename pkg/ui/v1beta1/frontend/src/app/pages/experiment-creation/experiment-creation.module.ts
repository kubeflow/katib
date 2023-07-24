import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatStepperModule } from '@angular/material/stepper';

import { TitleActionsToolbarModule, FormModule } from 'kubeflow';

import { ExperimentCreationComponent } from './experiment-creation.component';
import { ExperimentFormService } from '../../services/experiment-form.service';
import { FormMetadataModule } from './metadata/metadata.module';
import { FormTrialThresholdsModule } from './trial-thresholds/trial-thresholds.module';
import { FormObjectiveModule } from './objective/objective.module';
import { FormAlgorithmModule } from './algorithm/algorithm.module';
import { FormHyperParametersModule } from './hyper-parameters/hyper-parameters.module';
import { FormNasGraphModule } from './nas-graph/nas-graph.module';
import { FormNasOperationsModule } from './nas-operations/nas-operations.module';
import { FormMetricsCollectorModule } from './metrics-collector/metrics-collector.module';
import { FormTrialTemplateModule } from './trial-template/trial-template.module';
import { YamlModalModule } from './yaml-modal/yaml-modal.module';
import { FormEarlyStoppingModule } from './early-stopping/early-stopping.module';

@NgModule({
  declarations: [ExperimentCreationComponent],
  imports: [
    CommonModule,
    MatIconModule,
    MatStepperModule,
    FormModule,
    TitleActionsToolbarModule,
    FormMetadataModule,
    FormTrialThresholdsModule,
    FormObjectiveModule,
    FormAlgorithmModule,
    FormEarlyStoppingModule,
    FormHyperParametersModule,
    FormNasGraphModule,
    FormNasOperationsModule,
    FormMetricsCollectorModule,
    FormTrialTemplateModule,
    YamlModalModule,
  ],
  providers: [ExperimentFormService],
  exports: [ExperimentCreationComponent],
})
export class ExperimentCreationModule {}
