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

import { closeDialog, editTemplate, changeTemplate } from '../../../actions/templateActions';

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
});

class EditDialog extends React.Component {
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

  submitEditTemplate = () => {
    this.props.editTemplate(
      this.props.edittedTemplateNamespace,
      this.props.edittedTemplateConfigMapName,
      this.props.currentTemplateName,
      this.props.edittedTemplateName,
      this.props.edittedTemplateYaml,
    );
  };

  render() {
    const { classes } = this.props;
    return (
      <Dialog open={this.props.editOpen} onClose={this.props.closeDialog} maxWidth={'xl'}>
        <DialogTitle id="alert-dialog-title" className={classes.header}>
          {'Template Editor'}
          <Typography className={classes.headerTypography}>
            {'Namespace: ' + this.props.edittedTemplateNamespace}
          </Typography>

          <Typography className={classes.headerTypography}>
            {'ConfigMap: ' + this.props.edittedTemplateConfigMapName}
          </Typography>
        </DialogTitle>
        <DialogContent>
          <TextField
            value={this.props.edittedTemplateName}
            className={classes.textField}
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
            onClick={this.submitEditTemplate}
            color={'primary'}
          >
            Save
          </Button>
          <Button onClick={this.props.closeDialog} color={'primary'}>
            Discard
          </Button>
        </DialogActions>
      </Dialog>
    );
  }
}

const mapStateToProps = state => {
  return {
    editOpen: state[module].editOpen,
    edittedTemplateNamespace: state[module].edittedTemplateNamespace,
    edittedTemplateConfigMapName: state[module].edittedTemplateConfigMapName,
    currentTemplateName: state[module].currentTemplateName,
    edittedTemplateName: state[module].edittedTemplateName,
    edittedTemplateYaml: state[module].edittedTemplateYaml,
    edittedTemplateConfigMapSelectList: state[module].edittedTemplateConfigMapSelectList,
  };
};

export default connect(mapStateToProps, { closeDialog, editTemplate, changeTemplate })(
  withStyles(styles)(EditDialog),
);
