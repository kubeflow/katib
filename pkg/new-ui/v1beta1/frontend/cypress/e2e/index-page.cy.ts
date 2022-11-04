import { STATUS_TYPE } from 'kubeflow';
import { parseStatus } from 'src/app/pages/experiments/utils';

describe('Index page', () => {
  beforeEach(() => {
    cy.mockDashboardRequest();
    cy.mockNamespacesRequest();
    cy.fixture('settings').then(settings => {
      cy.mockExperimentsRequest(settings.namespace);
    });
    cy.fixture('experiments').as('experimentsRequest');
  });

  it('should have an "Experiments" title', () => {
    cy.visit('/');
    cy.get('[data-cy-toolbar-title]').contains('Experiments').should('exist');
  });

  it('should list Experiments without errors', () => {
    cy.visit('/');
    // wait for the requests to complete
    cy.wait(['@mockNamespacesRequest', '@mockExperimentsRequest']);

    // after fetching the data the page should not have an error snackbar
    cy.get('[data-cy-snack-status=ERROR]').should('not.exist');
  });

  // We use function () in order to be able to access aliases via this
  // tslint:disable-next-line: space-before-function-paren
  it('renders every Experiment name into the table', function () {
    cy.visit('/');
    cy.wait(['@mockNamespacesRequest', '@mockExperimentsRequest']);

    let i = 0;
    const experiments = this.experimentsRequest;
    // Table is sorted by Name in ascending order by default
    // and experiment objects are also sorted alphabetically by name
    cy.get(`[data-cy-resource-table-row="Name"]`).each(element => {
      expect(element).to.contain(experiments[i].name);
      i++;
    });
  });

  // tslint:disable-next-line: space-before-function-paren
  it('renders properly Status icon for all  experiments', function () {
    cy.visit('/');
    cy.wait(['@mockNamespacesRequest', '@mockExperimentsRequest']);

    let i = 0;
    const experiments = this.experimentsRequest;
    cy.get('[data-cy-resource-table-row="Status"]').each(element => {
      const status = parseStatus(experiments[i]);
      if (status.phase === STATUS_TYPE.READY) {
        cy.wrap(element)
          .find('lib-status>mat-icon')
          .should('contain', 'check_circle');
      } else if (status.phase === STATUS_TYPE.STOPPED) {
        cy.wrap(element)
          .find('lib-status>lib-icon')
          .should('have.attr', 'icon', 'custom:stoppedResource');
      } else if (status.phase === STATUS_TYPE.UNAVAILABLE) {
        cy.wrap(element)
          .find('lib-status>mat-icon')
          .should('contain', 'timelapse');
      } else if (status.phase === STATUS_TYPE.WARNING) {
        cy.wrap(element)
          .find('lib-status>mat-icon')
          .should('contain', 'warning');
      } else if (
        status.phase === STATUS_TYPE.WAITING ||
        status.phase === STATUS_TYPE.TERMINATING
      ) {
        cy.wrap(element).find('mat-spinner').should('exist');
      }
      i++;
    });
  });
});
