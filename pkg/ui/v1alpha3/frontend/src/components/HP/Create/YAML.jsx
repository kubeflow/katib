import React from 'react';
import { connect } from 'react-redux';
import makeStyles from '@material-ui/styles/makeStyles';
import 'brace/mode/javascript';
import 'brace/theme/tomorrow';
import AceEditor from 'react-ace';
import Button from '@material-ui/core/Button';

import { changeYaml } from '../../../actions/hpCreateActions';
import { submitYaml } from '../../../actions/generalActions';

const module = 'hpCreate';
const generalModule = 'general';

const useStyles = makeStyles({
  root: {
    flexGrow: 1,
  },
  editor: {
    margin: '0 auto',
  },
  submit: {
    textAlign: 'center',
    marginTop: 10,
  },
  progress: {
    height: 10,
    margin: 10,
  },
  close: {
    padding: 4,
  },
});

const YAML = props => {
  const onYamlChange = value => {
    props.changeYaml(value);
  };

  const submitWholeYaml = () => {
    props.submitYaml(props.yaml, props.globalNamespace);
  };

  const classes = useStyles();
  return (
    <div className={classes.root}>
      <h1>Generate</h1>
      <hr />
      {/* {props.loading && <LinearProgress className={classes.progress}/>} */}
      <div className={classes.editor}>
        <AceEditor
          mode="text"
          theme="tomorrow"
          value={props.yaml}
          onChange={onYamlChange}
          name="yaml-editor"
          editorProps={{ $blockScrolling: true }}
          tabSize={2}
          enableLiveAutocompletion={true}
          fontSize={14}
          width={'100%'}
          height={700}
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
    yaml: state[module].currentYaml,
    globalNamespace: state[generalModule].globalNamespace,
  };
};

export default connect(mapStateToProps, { changeYaml, submitYaml })(YAML);
