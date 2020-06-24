import React from 'react';
import { connect } from 'react-redux';
import withStyles from '@material-ui/styles/withStyles';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/mode-yaml';

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
import OutlinedInput from '@material-ui/core/OutlinedInput';

import { closeDialog, addTemplate, changeTemplate } from '../../../actions/templateActions';

import { TEMPLATE_MODULE } from '../../../constants/constants';

const styles = theme => ({
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
  selectForm: {
    margin: 10,
  },
  selectDiv: {
    textAlign: 'center',
  },
});

//TODO: Add functionality to create new ConfigMap with Trial Template
class AddDialog extends React.Component {
  onConfigMapNamespaceChange = event => {
    let templateData = this.props.trialTemplatesData;
    let newConfigMapNamespace = event.target.value;
    let newConfigMapName = this.props.updatedConfigMapName;

    if (newConfigMapNamespace !== this.props.updatedConfigMapNamespace) {
      let namespaceIndex = templateData.findIndex(function(trialTemplate, i) {
        return trialTemplate.ConfigMapNamespace === newConfigMapNamespace;
      });

      newConfigMapName = templateData[namespaceIndex].ConfigMaps[0].ConfigMapName;
    }

    this.props.changeTemplate(
      newConfigMapNamespace,
      newConfigMapName,
      this.props.updatedConfigMapPath,
      this.props.updatedTemplateYaml,
    );
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
  };

  render() {
    const { classes } = this.props;

    return this.props.configMapNamespaceIndex !== -1 ? (
      <div>
        <Dialog open={this.props.addOpen} onClose={this.props.closeDialog}>
          <DialogTitle id="alert-dialog-title" className={classes.header}>
            {'Template Creator'}
            <Typography className={classes.headerTypography}>
              {'Select ConfigMap Namespace and Name'}
            </Typography>
          </DialogTitle>
          <DialogContent>
            <div className={classes.selectDiv}>
              <FormControl variant="outlined" className={classes.selectForm}>
                <InputLabel>Namespace</InputLabel>
                <Select
                  value={this.props.updatedConfigMapNamespace}
                  onChange={this.onConfigMapNamespaceChange}
                  className={classes.selectBox}
                  input={<OutlinedInput labelWidth={90} />}
                >
                  {this.props.trialTemplatesData.map((trialTemplate, i) => {
                    return (
                      <MenuItem value={trialTemplate.ConfigMapNamespace} key={i}>
                        {trialTemplate.ConfigMapNamespace}
                      </MenuItem>
                    );
                  })}
                </Select>
              </FormControl>
              <FormControl variant="outlined" className={classes.selectForm}>
                <InputLabel>Name</InputLabel>
                <Select
                  value={this.props.updatedConfigMapName}
                  onChange={this.onConfigMapNameChange}
                  className={classes.selectBox}
                  input={<OutlinedInput labelWidth={50} />}
                >
                  {this.props.trialTemplatesData[this.props.configMapNamespaceIndex].ConfigMaps.map(
                    (configMap, i) => {
                      return (
                        <MenuItem value={configMap.ConfigMapName} key={i}>
                          {configMap.ConfigMapName}
                        </MenuItem>
                      );
                    },
                  )}
                </Select>
              </FormControl>
            </div>
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
                !this.props.updatedConfigMapPath ||
                !this.props.updatedTemplateYaml ||
                // Path can't contain spaces
                this.props.updatedConfigMapPath.indexOf(' ') !== -1
              }
              onClick={this.submitAddTemplate}
              color={'primary'}
            >
              Save
            </Button>
            <Button onClick={this.props.closeDialog} color={'primary'}>
              Discard
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    ) : (
      <div />
    );
  }
}

const mapStateToProps = state => {
  let templatesData = state[TEMPLATE_MODULE].trialTemplatesData;

  let nsIndex = templatesData.findIndex(function(trialTemplate, i) {
    return trialTemplate.ConfigMapNamespace === state[TEMPLATE_MODULE].updatedConfigMapNamespace;
  });

  return {
    addOpen: state[TEMPLATE_MODULE].addOpen,
    trialTemplatesData: state[TEMPLATE_MODULE].trialTemplatesData,
    configMapNamespaceIndex: nsIndex,
    updatedConfigMapNamespace: state[TEMPLATE_MODULE].updatedConfigMapNamespace,
    updatedConfigMapName: state[TEMPLATE_MODULE].updatedConfigMapName,
    updatedConfigMapPath: state[TEMPLATE_MODULE].updatedConfigMapPath,
    updatedTemplateYaml: state[TEMPLATE_MODULE].updatedTemplateYaml,
  };
};

export default connect(mapStateToProps, {
  closeDialog,
  addTemplate,
  changeTemplate,
})(withStyles(styles)(AddDialog));
