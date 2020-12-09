import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/mode-yaml';

import { withStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import TextField from '@material-ui/core/TextField';
import DialogTitle from '@material-ui/core/DialogTitle';
import Typography from '@material-ui/core/Typography';
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Checkbox from '@material-ui/core/Checkbox';
import Grid from '@material-ui/core/Grid';

import { closeDialog, addTemplate, changeTemplate } from '../../../actions/templateActions';

import { TEMPLATE_MODULE, GENERAL_MODULE } from '../../../constants/constants';

const styles = () => ({
  header: {
    textAlign: 'center',
    width: 650,
  },
  headerTypography: {
    textAlign: 'center',
    marginTop: 5,
    fontSize: 19,
  },
  textField: {
    marginBottom: 10,
    width: '100%',
  },
  selectBox: {
    width: 200,
  },
  textFieldConfigMap: {
    marginLeft: 10,
    marginRight: 10,
  },
  selectForm: {
    margin: 10,
  },
  checkBox: {
    marginLeft: 'auto',
  },
});

class AddDialog extends React.Component {
  constructor(props) {
    super(props);
    this.state = { checkedNewName: false };
  }

  onConfigMapNamespaceChange = event => {
    let templateData = this.props.trialTemplatesData;
    let newConfigMapNamespace = event.target.value;
    let newConfigMapName = '';

    let namespaceIndex = templateData.findIndex(function (trialTemplate, i) {
      return trialTemplate.ConfigMapNamespace === newConfigMapNamespace;
    });

    // Assign new ConfigMap name only if namespace exists in Template data
    if (newConfigMapNamespace !== this.props.updatedConfigMapNamespace && namespaceIndex !== -1) {
      newConfigMapName = templateData[namespaceIndex].ConfigMaps[0].ConfigMapName;
    }

    this.props.changeTemplate(
      newConfigMapNamespace,
      newConfigMapName,
      this.props.updatedConfigMapPath,
      this.props.updatedTemplateYaml,
    );

    // Reset check box when changing namespace
    this.setState({ checkedNewName: false });
  };

  onConfigMapNameChange = event => {
    this.props.changeTemplate(
      this.props.updatedConfigMapNamespace,
      event.target.value,
      this.props.updatedConfigMapPath,
      this.props.updatedTemplateYaml,
    );
  };

  onConfigMapPathChange = event => {
    this.props.changeTemplate(
      this.props.updatedConfigMapNamespace,
      this.props.updatedConfigMapName,
      event.target.value,
      this.props.updatedTemplateYaml,
    );
  };

  onTemplateYamlChange = newTemplateYaml => {
    this.props.changeTemplate(
      this.props.updatedConfigMapNamespace,
      this.props.updatedConfigMapName,
      this.props.updatedConfigMapPath,
      newTemplateYaml,
    );
  };

  submitAddTemplate = () => {
    this.props.addTemplate(
      this.props.updatedConfigMapNamespace,
      this.props.updatedConfigMapName,
      this.props.updatedConfigMapPath,
      this.props.updatedTemplateYaml,
    );
    this.setState({ checkedNewName: false });
  };

  onCheckBoxChange = event => {
    this.setState({ checkedNewName: event.target.checked });
    if (event.target.checked) {
      this.props.changeTemplate(
        this.props.updatedConfigMapNamespace,
        '',
        this.props.updatedConfigMapPath,
        this.props.updatedTemplateYaml,
      );
    } else {
      this.props.changeTemplate(
        this.props.updatedConfigMapNamespace,
        this.props.trialTemplatesData[this.props.configMapNamespaceIndex].ConfigMaps[0]
          .ConfigMapName,
        this.props.updatedConfigMapPath,
        this.props.updatedTemplateYaml,
      );
    }
  };

  onDialogClose = () => {
    this.props.closeDialog();
    this.setState({ checkedNewName: false });
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        <Dialog open={this.props.addOpen} onClose={this.onDialogClose} maxWidth="md">
          <DialogTitle id="alert-dialog-title" className={classes.header}>
            {'Template Creator'}
            <Typography className={classes.headerTypography}>
              {'Select ConfigMap Namespace and Name'}
            </Typography>
          </DialogTitle>
          <DialogContent>
            <Grid container alignItems="center">
              <Grid item>
                <FormControl variant="outlined" className={classes.selectForm}>
                  <InputLabel>Namespace</InputLabel>
                  <Select
                    value={this.props.updatedConfigMapNamespace}
                    onChange={this.onConfigMapNamespaceChange}
                    className={classes.selectBox}
                    label="Namespace"
                  >
                    {this.props.namespaces
                      .filter(namespace => namespace !== 'All namespaces')
                      .map((namespace, i) => {
                        return (
                          <MenuItem value={namespace} key={i}>
                            {namespace}
                          </MenuItem>
                        );
                      })}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item>
                {!this.state.checkedNewName && this.props.configMapNamespaceIndex !== -1 ? (
                  <FormControl variant="outlined" className={classes.selectForm}>
                    <InputLabel>Name</InputLabel>
                    <Select
                      value={this.props.updatedConfigMapName}
                      onChange={this.onConfigMapNameChange}
                      className={classes.selectBox}
                      label="Name"
                    >
                      {this.props.trialTemplatesData[
                        this.props.configMapNamespaceIndex
                      ].ConfigMaps.map((configMap, i) => {
                        return (
                          <MenuItem value={configMap.ConfigMapName} key={i}>
                            {configMap.ConfigMapName}
                          </MenuItem>
                        );
                      })}
                    </Select>
                  </FormControl>
                ) : (
                  <TextField
                    variant="outlined"
                    label="Name"
                    className={classes.textFieldConfigMap}
                    value={this.props.updatedConfigMapName}
                    onChange={this.onConfigMapNameChange}
                  />
                )}
              </Grid>
              {this.props.configMapNamespaceIndex !== -1 && (
                <Grid item>
                  <FormControlLabel
                    className={classes.checkBox}
                    control={
                      <Checkbox
                        checked={this.state.checkedNewName}
                        onChange={this.onCheckBoxChange}
                        color="primary"
                      />
                    }
                    label="New ConfigMap"
                  />
                </Grid>
              )}
            </Grid>
            <TextField
              className={classes.textField}
              value={this.props.updatedConfigMapPath}
              onChange={this.onConfigMapPathChange}
              label="Template ConfigMap Path"
              placeholder="Template ConfigMap Path"
            />
            <br />
            <AceEditor
              mode="yaml"
              theme="sqlserver"
              value={this.props.updatedTemplateYaml}
              tabSize={2}
              fontSize={13}
              width={'100%'}
              showPrintMargin={false}
              autoScrollEditorIntoView={true}
              maxLines={30}
              minLines={10}
              onChange={this.onTemplateYamlChange}
            />
          </DialogContent>
          <DialogActions>
            <Button
              disabled={
                // Config Map name can't contain spaces and must exists
                !this.props.updatedConfigMapName ||
                this.props.updatedConfigMapName.indexOf(' ') !== -1 ||
                // ConfigMap name must be unique, when state.checkedNewName = true and configMapNamespaceIndex != -1
                (this.state.checkedNewName &&
                  this.props.configMapNamespaceIndex !== -1 &&
                  this.props.trialTemplatesData[this.props.configMapNamespaceIndex].ConfigMaps.some(
                    t => t.ConfigMapName === this.props.updatedConfigMapName,
                  )) ||
                // Path can't contain spaces and must exists
                !this.props.updatedConfigMapPath ||
                this.props.updatedConfigMapPath.indexOf(' ') !== -1 ||
                // Path in ConfigMap must be unique, when configMapNameIndex != -1
                (this.props.configMapNameIndex !== -1 &&
                  this.props.trialTemplatesData[this.props.configMapNamespaceIndex].ConfigMaps[
                    this.props.configMapNameIndex
                  ].Templates.some(t => t.Path === this.props.updatedConfigMapPath)) ||
                // Yaml must exists
                !this.props.updatedTemplateYaml
              }
              onClick={this.submitAddTemplate}
              color={'primary'}
            >
              Save
            </Button>
            <Button onClick={this.onDialogClose} color={'primary'}>
              Discard
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    );
  }
}

const mapStateToProps = state => {
  let templatesData = state[TEMPLATE_MODULE].trialTemplatesData;

  let nsIndex = templatesData.findIndex(function (trialTemplate, i) {
    return trialTemplate.ConfigMapNamespace === state[TEMPLATE_MODULE].updatedConfigMapNamespace;
  });

  let cmIndex = -1;
  if (nsIndex !== -1) {
    cmIndex = templatesData[nsIndex].ConfigMaps.findIndex(function (configMap, i) {
      return configMap.ConfigMapName === state[TEMPLATE_MODULE].updatedConfigMapName;
    });
  }

  return {
    addOpen: state[TEMPLATE_MODULE].addOpen,
    trialTemplatesData: templatesData,
    configMapNamespaceIndex: nsIndex,
    configMapNameIndex: cmIndex,
    updatedConfigMapNamespace: state[TEMPLATE_MODULE].updatedConfigMapNamespace,
    updatedConfigMapName: state[TEMPLATE_MODULE].updatedConfigMapName,
    updatedConfigMapPath: state[TEMPLATE_MODULE].updatedConfigMapPath,
    updatedTemplateYaml: state[TEMPLATE_MODULE].updatedTemplateYaml,
    namespaces: state[GENERAL_MODULE].namespaces,
  };
};

export default connect(mapStateToProps, {
  closeDialog,
  addTemplate,
  changeTemplate,
})(withStyles(styles)(AddDialog));
