import './commands';

// types of the custom commands
// Must be declared global to be detected by typescript (allows import/export)
declare global {
  namespace Cypress {
    interface Chainable {
      /**
       * Custom command to mock request at '/dashboard_lib.bundle.js'
       */
      mockDashboardRequest(): Chainable<void>;

      /**
       * Custom command to mock request at '/katib/fetch_namespaces'
       */
      mockNamespacesRequest(): Chainable<void>;

      /**
       * Custom command to mock request at 'katib/fetch_namespaced_experiments/?namespace=<namespace>'
       * and return an array with mock Experiments []
       */
      mockExperimentsRequest(namespace: string): Chainable<void>;

      /**
       * Custom command to mock request at '/katib/fetch_trial_templates'
       */
      mockTrialTemplate(): Chainable<void>;
    }
  }
}
