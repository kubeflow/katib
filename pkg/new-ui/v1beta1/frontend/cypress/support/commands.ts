Cypress.Commands.add('mockDashboardRequest', () => {
  cy.intercept('GET', '/dashboard_lib.bundle.js', { body: [] }).as(
    'mockDashboardRequest',
  );
});

Cypress.Commands.add('mockNamespacesRequest', () => {
  cy.intercept('GET', '/katib/fetch_namespaces', {
    fixture: 'namespaces',
  }).as('mockNamespacesRequest');
});

Cypress.Commands.add('mockExperimentsRequest', namespace => {
  cy.intercept(
    'GET',
    `katib/fetch_experiments/?namespace=${namespace}`,
    {
      fixture: 'experiments',
    },
  ).as('mockExperimentsRequest');
});

Cypress.Commands.add('mockTrialTemplate', () => {
  cy.intercept('GET', `/katib/fetch_trial_templates/`, {
    fixture: 'trial-template',
  }).as('mockTrialTemplate');
});
