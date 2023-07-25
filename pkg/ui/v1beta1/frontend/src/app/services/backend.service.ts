import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { catchError, map } from 'rxjs/operators';
import { Observable, throwError } from 'rxjs';
import {
  BackendService,
  K8sObject,
  SnackBarConfig,
  SnackBarService,
  SnackType,
} from 'kubeflow';

import { Experiments } from '../models/experiment.model';
import { ExperimentK8s } from '../models/experiment.k8s.model';
import { TrialK8s } from '../models/trial.k8s.model';
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

    const config: SnackBarConfig = {
      data: {
        msg,
        snackType: SnackType.Error,
      },
      duration: 8000,
    };
    this.snack.open(config);

    return throwError(msg);
  }

  getExperimentsSingleNamespace(namespace: string): Observable<Experiments> {
    // If the route doesn't end in a "/"" then the backend will return a 301 to
    // the url ending with "/".
    const url = `/katib/fetch_experiments/?namespace=${namespace}`;

    return this.http.get<any>(url).pipe(
      catchError(error => this.parseError(error)),
      map((resp: any) => {
        return resp;
      }),
    );
  }

  getExperimentsAllNamespaces(namespaces: string[]): Observable<Experiments> {
    return this.getObjectsAllNamespaces(
      this.getExperimentsSingleNamespace.bind(this),
      namespaces,
    );
  }

  getExperiments(ns: string | string[]): Observable<Experiments> {
    if (Array.isArray(ns)) {
      return this.getExperimentsAllNamespaces(ns);
    }

    return this.getExperimentsSingleNamespace(ns);
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

  getTrialInfo(name: string, namespace: string): Observable<TrialK8s> {
    const url = `/katib/fetch_trial/?trialName=${name}&namespace=${namespace}`;
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

  getTrialLogs(name: string, namespace: string): Observable<any> {
    const url = `/katib/fetch_trial_logs/?trialName=${name}&namespace=${namespace}`;

    return this.http
      .get(url)
      .pipe(catchError(error => this.handleError(error, false)));
  }

  // ---------------------------Error Handling---------------------------------

  // Override common service's getBackendErrorLog
  // in order to properly show the message the backend has sent
  public getBackendErrorLog(error: HttpErrorResponse): string {
    if (error.error === null) {
      return error.message;
    } else {
      // Show the message the backend has sent
      return error.error.log ? error.error.log : error.error;
    }
  }
}
