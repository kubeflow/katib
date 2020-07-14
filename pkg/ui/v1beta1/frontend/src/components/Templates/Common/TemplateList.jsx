import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import Divider from '@material-ui/core/Divider';
import Button from '@material-ui/core/Button';
import LinearProgress from '@material-ui/core/LinearProgress';

import TemplatePanel from './TemplatePanel';
import FilterPanel from './FilterPanel';
import AddDialog from './AddDialog';
import EditDialog from './EditDialog';
import DeleteDialog from './DeleteDialog';

import { openDialog } from '../../../actions/templateActions';

import { TEMPLATE_MODULE, GENERAL_MODULE } from '../../../constants/constants';

const styles = theme => ({
  namespace: {
    marginTop: 25,
    marginRight: 15,
    fontSize: theme.typography.pxToRem(26),
  },
  configMap: {
    margin: 15,
    fontSize: theme.typography.pxToRem(23),
  },
  templatesBlock: {
    width: '96%',
    margin: '0 auto',
  },
  template: {
    fontSize: theme.typography.pxToRem(20),
    fontWeight: theme.typography.fontWeightRegular,
  },
  divider: {
    marginTop: 20,
  },
  buttonAdd: {
    textAlign: 'center',
  },
  noTemplates: {
    marginTop: 25,
    marginRight: 15,
    fontSize: theme.typography.pxToRem(50),
  },
  loading: {
    marginTop: 30,
  },
});

const dialogTypeAdd = 'add';

class TemplateList extends React.Component {
  openAddDialog = noTrialTemplates => () => {
    if (noTrialTemplates) {
      this.props.openDialog(dialogTypeAdd, this.props.namespaces[1]);
    } else {
      this.props.openDialog(
        dialogTypeAdd,
        this.props.trialTemplatesData[0].ConfigMapNamespace,
        this.props.trialTemplatesData[0].ConfigMaps[0].ConfigMapName,
      );
    }
  };

  render() {
    const { classes } = this.props;
    return (
      <div>
        {this.props.loading ? (
          <LinearProgress color={'primary'} className={classes.loading} />
        ) : (
          <div>
            {this.props.trialTemplatesData.length !== 0 ? (
              <div>
                <FilterPanel />
                <div className={classes.buttonAdd}>
                  <Button
                    variant={'contained'}
                    color={'primary'}
                    onClick={this.openAddDialog(false)}
                  >
                    Add Template
                  </Button>
                </div>
                {this.props.filteredTrialTemplatesData.map((trialTemplate, nsIndex) => {
                  return (
                    <div key={nsIndex}>
                      <Grid key={nsIndex} container>
                        <Grid item>
                          <Typography className={classes.namespace}>Namespace:</Typography>
                        </Grid>
                        <Grid item>
                          <Typography className={classes.namespace} style={{ fontWeight: 'bold' }}>
                            {trialTemplate.ConfigMapNamespace}
                          </Typography>
                        </Grid>
                        <Grid item xs={12}>
                          <hr />
                        </Grid>
                      </Grid>

                      {trialTemplate.ConfigMaps.map((configMap, cmIndex) => {
                        return (
                          <div key={cmIndex}>
                            <Grid container>
                              <Grid item>
                                <Typography className={classes.configMap}>ConfigMap:</Typography>
                              </Grid>
                              <Grid item>
                                <Typography
                                  className={classes.configMap}
                                  style={{ fontStyle: 'italic' }}
                                >
                                  {configMap.ConfigMapName}
                                </Typography>
                              </Grid>
                            </Grid>

                            {configMap.Templates.map((template, templateIndex) => {
                              return (
                                <div className={classes.templatesBlock} key={templateIndex}>
                                  <ExpansionPanel>
                                    <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                                      <Typography className={classes.template}>
                                        {template.Path}
                                      </Typography>
                                    </ExpansionPanelSummary>
                                    <ExpansionPanelDetails>
                                      <TemplatePanel
                                        configMapNamespace={trialTemplate.ConfigMapNamespace}
                                        configMapName={configMap.ConfigMapName}
                                        configMapPath={template.Path}
                                        templateYaml={template.Yaml}
                                      />
                                    </ExpansionPanelDetails>
                                  </ExpansionPanel>
                                </div>
                              );
                            })}
                            <Divider className={classes.divider} />
                          </div>
                        );
                      })}
                    </div>
                  );
                })}

                <EditDialog />
                <DeleteDialog />
              </div>
            ) : (
              <div>
                <Typography className={classes.namespace}>
                  No ConfigMaps with Katib Trial Templates
                </Typography>
                <div className={classes.buttonAdd}>
                  <Button
                    variant={'contained'}
                    color={'primary'}
                    onClick={this.openAddDialog(true)}
                  >
                    Add Template
                  </Button>
                </div>
              </div>
            )}
            <AddDialog />
          </div>
        )}
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    filteredTrialTemplatesData: state[TEMPLATE_MODULE].filteredTrialTemplatesData,
    trialTemplatesData: state[TEMPLATE_MODULE].trialTemplatesData,
    loading: state[TEMPLATE_MODULE].loading,
    namespaces: state[GENERAL_MODULE].namespaces,
  };
};

export default connect(mapStateToProps, { openDialog })(withStyles(styles)(TemplateList));
