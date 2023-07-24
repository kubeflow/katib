describe('New Experiment form page', () => {
  beforeEach(() => {
    cy.mockDashboardRequest();
    cy.mockTrialTemplate();
  });

  it('New Experiment form page loads template without errors', () => {
    cy.visit('/new');
    cy.wait(['@mockDashboardRequest', '@mockTrialTemplate']);

    // after fetching the data the page should not have an error snackbar
    cy.get('[data-cy-snack-status=ERROR]').should('not.exist');
  });
});
