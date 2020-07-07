import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/mode-yaml';

import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import CreateIcon from '@material-ui/icons/Create';

import { openDialog } from '../../../actions/templateActions';

const useStyles = makeStyles({
  root: {
    width: '98%',
    margin: '0 auto',
  },

  buttons: {
    marginTop: 30,
    marginLeft: 20,
  },
  icon: {
    margin: 4,
  },
});
const dialogTypeEdit = 'edit';
const dialogTypeDelete = 'delete';

const TemplatePanel = props => {
  const classes = useStyles();

  const openEditDialog = (configMapNamespace, configMapName, configMapPath, templateYaml) => () => {
    props.openDialog(
      dialogTypeEdit,
      configMapNamespace,
      configMapName,
      configMapPath,
      templateYaml,
    );
  };

  const openDeleteDialog = (configMapNamespace, configMapName, configMapPath) => () => {
    props.openDialog(dialogTypeDelete, configMapNamespace, configMapName, configMapPath);
  };

  return (
    <div className={classes.root}>
      <AceEditor
        mode="yaml"
        theme="sqlserver"
        value={props.templateYaml}
        tabSize={2}
        fontSize={16}
        width={'100%'}
        showPrintMargin={false}
        autoScrollEditorIntoView={true}
        readOnly={true}
        maxLines={160}
        minLines={10}
      />
      <Button
        className={classes.buttons}
        variant={'contained'}
        color={'primary'}
        onClick={openEditDialog(
          props.configMapNamespace,
          props.configMapName,
          props.configMapPath,
          props.templateYaml,
        )}
      >
        <CreateIcon color={'secondary'} className={classes.icon} />
        Edit
      </Button>

      <Button
        className={classes.buttons}
        variant={'contained'}
        color={'primary'}
        onClick={openDeleteDialog(
          props.configMapNamespace,
          props.configMapName,
          props.configMapPath,
        )}
      >
        <DeleteIcon color={'secondary'} className={classes.icon} />
        Delete
      </Button>
    </div>
  );
};

const mapStateToProps = () => {
  return {};
};

export default connect(mapStateToProps, { openDialog })(TemplatePanel);
