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

  submitEditTemplate = () => {
    this.props.editTemplate(
      this.props.updatedConfigMapNamespace,
      this.props.updatedConfigMapName,
      this.props.configMapPath,
      this.props.updatedConfigMapPath,
      this.props.updatedTemplateYaml,
    );
  };

  render() {
    const { classes } = this.props;
    return (
      <Dialog open={this.props.editOpen} onClose={this.props.closeDialog} maxWidth={'xl'}>
        <DialogTitle id="alert-dialog-title" className={classes.header}>
          {'Template Editor'}
          <Typography className={classes.headerTypography}>
            {'ConfigMap Namespace: ' + this.props.updatedConfigMapNamespace}
          </Typography>

          <Typography className={classes.headerTypography}>
            {'ConfigMap Name: ' + this.props.updatedConfigMapName}
          </Typography>
        </DialogTitle>
        <DialogContent>
          <TextField
            value={this.props.updatedConfigMapPath}
            className={classes.textField}
            onChange={this.onConfigMapPathChange}
            label="Template Config Map Path"
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
            disabled={!this.props.updatedConfigMapPath || !this.props.updatedTemplateYaml}
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
    updatedConfigMapNamespace: state[module].updatedConfigMapNamespace,
    updatedConfigMapName: state[module].updatedConfigMapName,
    configMapPath: state[module].configMapPath,
    updatedConfigMapPath: state[module].updatedConfigMapPath,
    updatedTemplateYaml: state[module].updatedTemplateYaml,
  };
};

export default connect(mapStateToProps, { closeDialog, editTemplate, changeTemplate })(
  withStyles(styles)(EditDialog),
);
