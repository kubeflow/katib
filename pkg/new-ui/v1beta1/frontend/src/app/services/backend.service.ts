import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError, map } from 'rxjs/operators';
import { Observable, throwError } from 'rxjs';
import {
  BackendService,
  K8sObject,
  SnackBarService,
  SnackType,
} from 'kubeflow';

import { Experiments } from '../models/experiment.model';
import { ExperimentK8s } from '../models/experiment.k8s.model';
import { TrialTemplateResponse } from '../models/trial-templates.model';

@Injectable({
  providedIn: 'root',
})
export class KWABackendService extends BackendService {
  constructor(public http: HttpClient, public snack: SnackBarService) {
    super(http, snack);
  }

  private parseError(error) {
    let msg = 'An error occured while talking to the backend';

    if (error instanceof HttpErrorResponse) {
      if (!error.status) {
        msg = error.message;
      } else {
        msg = `[${error.status}] ${error.error}`;
      }
    }

    this.snack.open(msg, SnackType.Error, 8000);

    return throwError(msg);
  }

  getExperiments(): Observable<Experiments> {
    // If the route doesn't end in a "/"" then the backend will return a 301 to
    // the url ending with "/".
    const url = '/katib/fetch_experiments/';

    return this.http.get<any>(url).pipe(
      catchError(error => this.parseError(error)),
      map((resp: any) => {
        return resp;
      }),
    );
  }

  getExperimentTrialsInfo(name: string, namespace: string): Observable<any> {
    const url = `/katib/fetch_hp_job_info/?experimentName=${name}&namespace=${namespace}`;

    return this.http.get(url).pipe(catchError(error => this.parseError(error)));
  }

  getExperiment(name: string, namespace: string): Observable<ExperimentK8s> {
    const url = `/katib/fetch_experiment/?experimentName=${name}&namespace=${namespace}`;

    return this.http
      .get(url)
      .pipe(
        catchError(error => this.parseError(error)),
      ) as Observable<K8sObject>;
  }

  deleteExperiment(name: string, namespace: string): Observable<ExperimentK8s> {
    const url = `/katib/delete_experiment/?experimentName=${name}&namespace=${namespace}`;

    return this.http
      .delete(url)
      .pipe(
        catchError(error => this.handleError(error)),
      ) as Observable<ExperimentK8s>;
  }

  getTrial(name: string, namespace: string): Observable<any> {
    const url = `/katib/fetch_hp_job_trial_info/?trialName=${name}&namespace=${namespace}`;
    return this.http.get(url).pipe(catchError(error => this.parseError(error)));
  }

  getTrialTemplates(namespace: string): Observable<TrialTemplateResponse> {
    const url = `/katib/fetch_trial_templates/`;

    return this.http
      .get<TrialTemplateResponse>(url)
      .pipe(catchError(error => this.parseError(error)));
  }

  createExperiment(exp: ExperimentK8s): Observable<any> {
    const url = `/katib/create_experiment/`;

    return this.http
      .post(url, { postData: exp })
      .pipe(catchError(error => this.parseError(error)));
  }
}
