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

const module = 'template';

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
  componentDidMount() {
    if (this.props.trialTemplatesList.length != 0) {
      let configMaps = this.props.trialTemplatesList[0].ConfigMapsList;
      if (configMaps.length != 0) {
        let configMapsList = [];
        configMaps.forEach(configMap => configMapsList.push(configMap.ConfigMapName));
        this.props.changeTemplate(
          this.props.edittedTemplateNamespace,
          this.props.edittedTemplateConfigMapName,
          this.props.edittedTemplateName,
          this.props.edittedTemplateYaml,
          configMapsList,
        );
      }
    }
  }

  onNamespaceChange = event => {
    let newNamespace = event.target.value;

    let namespaceIndex = this.props.trialTemplatesList.findIndex(function(trialTemplate, i) {
      return trialTemplate.Namespace === newNamespace;
    });

    if (this.props.trialTemplatesList.length != 0) {
      let configMaps = this.props.trialTemplatesList[namespaceIndex].ConfigMapsList;
      //TODO: add logic when configMapsList is empty
      if (configMaps.length != 0) {
        let configMapsList = [];
        configMaps.forEach(configMap => configMapsList.push(configMap.ConfigMapName));

        this.props.changeTemplate(
          newNamespace,
          configMapsList[0],
          this.props.edittedTemplateName,
          this.props.edittedTemplateYaml,
          configMapsList,
        );
      }
    }
  };

  onConfigMapNameChange = event => {
    this.props.changeTemplate(
      this.props.edittedTemplateNamespace,
      event.target.value,
      this.props.edittedTemplateName,
      this.props.edittedTemplateYaml,
      this.props.edittedTemplateConfigMapSelectList,
    );
  };

  onNameChange = event => {
    this.props.changeTemplate(
      this.props.edittedTemplateNamespace,
      this.props.edittedTemplateConfigMapName,
      event.target.value,
      this.props.edittedTemplateYaml,
      this.props.edittedTemplateConfigMapSelectList,
    );
  };

  onYamlChange = newTemplateYaml => {
    this.props.changeTemplate(
      this.props.edittedTemplateNamespace,
      this.props.edittedTemplateConfigMapName,
      this.props.edittedTemplateName,
      newTemplateYaml,
      this.props.edittedTemplateConfigMapSelectList,
    );
  };

  submitAddTemplate = () => {
    this.props.addTemplate(
      this.props.edittedTemplateNamespace,
      this.props.edittedTemplateConfigMapName,
      this.props.edittedTemplateName,
      this.props.edittedTemplateYaml,
    );
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        <Dialog open={this.props.addOpen} onClose={this.props.closeDialog}>
          <DialogTitle id="alert-dialog-title" className={classes.header}>
            {'Template Creator'}
            <Typography className={classes.headerTypography}>
              {'Select Namespace and ConfigMap'}
            </Typography>
          </DialogTitle>
          <DialogContent>
            <div className={classes.selectDiv}>
              <FormControl variant="outlined" className={classes.selectForm}>
                <InputLabel>Namespace</InputLabel>
                <Select
                  value={this.props.edittedTemplateNamespace}
                  onChange={this.onNamespaceChange}
                  className={classes.selectBox}
                  input={<OutlinedInput labelWidth={160} />}
                >
                  {this.props.trialTemplatesList.map((trialTemplate, i) => {
                    return (
                      <MenuItem value={trialTemplate.Namespace} key={i}>
                        {trialTemplate.Namespace}
                      </MenuItem>
                    );
                  })}
                </Select>
              </FormControl>
              <FormControl variant="outlined" className={classes.selectForm}>
                <InputLabel>ConfigMap</InputLabel>
                <Select
                  value={this.props.edittedTemplateConfigMapName}
                  onChange={this.onConfigMapNameChange}
                  className={classes.selectBox}
                  input={<OutlinedInput labelWidth={160} />}
                >
                  {this.props.edittedTemplateConfigMapSelectList.map((configMap, i) => {
                    return (
                      <MenuItem value={configMap} key={i}>
                        {configMap}
                      </MenuItem>
                    );
                  })}
                </Select>
              </FormControl>
            </div>
            <TextField
              className={classes.textField}
              value={this.props.edittedTemplateName}
              onChange={this.onNameChange}
              label="Template name"
              placeholder="Template name"
            />
            <br />
            <AceEditor
              mode="yaml"
              theme="sqlserver"
              value={this.props.edittedTemplateYaml}
              tabSize={2}
              fontSize={13}
              width={'100%'}
              showPrintMargin={false}
              autoScrollEditorIntoView={true}
              maxLines={30}
              minLines={10}
              onChange={this.onYamlChange}
            />
          </DialogContent>
          <DialogActions>
            <Button
              disabled={!this.props.edittedTemplateName || !this.props.edittedTemplateYaml}
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
    );
  }
}

const mapStateToProps = state => {
  return {
    addOpen: state[module].addOpen,
    trialTemplatesList: state[module].trialTemplatesList,
    edittedTemplateNamespace: state[module].edittedTemplateNamespace,
    edittedTemplateConfigMapName: state[module].edittedTemplateConfigMapName,
    edittedTemplateName: state[module].edittedTemplateName,
    edittedTemplateYaml: state[module].edittedTemplateYaml,
    edittedTemplateConfigMapSelectList: state[module].edittedTemplateConfigMapSelectList,
  };
};

export default connect(mapStateToProps, {
  closeDialog,
  addTemplate,
  changeTemplate,
})(withStyles(styles)(AddDialog));
