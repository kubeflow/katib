import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-sqlserver';
import 'ace-builds/src-noconflict/mode-yaml';

import { makeStyles } from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';

import { changeYaml } from '../../../actions/nasCreateActions';
import { submitYaml } from '../../../actions/generalActions';

import { GENERAL_MODULE, NAS_CREATE_MODULE } from '../../../constants/constants';

const useStyles = makeStyles({
  editor: {
    margin: '0 auto',
  },
  submit: {
    textAlign: 'center',
    marginTop: 10,
  },
  button: {
    margin: 15,
  },
});

const YAML = props => {
  const onYamlChange = value => {
    props.changeYaml(value);
  };

  const submitWholeYaml = () => {
    props.submitYaml(props.currentYaml, props.globalNamespace);
  };

  const classes = useStyles();
  return (
    <div>
      <Typography variant={'h5'}>{'Generate'}</Typography>
      <hr />
      <div className={classes.editor}>
        <AceEditor
          mode="yaml"
          theme="sqlserver"
          value={props.currentYaml}
          tabSize={2}
          fontSize={14}
          width={'auto'}
          showPrintMargin={false}
          autoScrollEditorIntoView={true}
          maxLines={32}
          minLines={32}
          onChange={onYamlChange}
        />
      </div>
      <div className={classes.submit}>
        <Button
          variant="contained"
          color={'primary'}
          className={classes.button}
          onClick={submitWholeYaml}
        >
          Deploy
        </Button>
      </div>
    </div>
  );
};

const mapStateToProps = state => {
  return {
    currentYaml: state[NAS_CREATE_MODULE].currentYaml,
    globalNamespace: state[GENERAL_MODULE].globalNamespace,
  };
};

export default connect(mapStateToProps, { changeYaml, submitYaml })(YAML);
