import React from 'react';
import { connect } from 'react-redux';

import AceEditor from 'react-ace';
import 'ace-builds/src-noconflict/theme-xcode';
import 'ace-builds/src-noconflict/mode-json';

import { withStyles } from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';

import { closeDialogSuggestion } from '../../actions/generalActions';

import { GENERAL_MODULE } from '../../constants/constants';

const styles = theme => ({
  header: {
    textAlign: 'center',
    width: 900,
  },
  aceMarker: {
    position: 'absolute',
    backgroundColor: '#FFFF00',
  },
});

const SuggestionStatusFailed = '"type": "Failed",';

const SuggestionInfoDialog = props => {
  const { classes } = props;

  let markers = [];
  let failedSuggestion = '';
  if (props.dialogSuggestionOpen) {
    let jsonString = JSON.stringify(props.suggestion, null, 2);
    let jsonSplit = jsonString.split(/\n/);
    for (var i = jsonSplit.length - 1; i >= 0; i--) {
      if (jsonSplit[i].trim() === SuggestionStatusFailed) {
        markers.push({
          startRow: i - 1,
          endRow: i + 7,
          className: classes.aceMarker,
          fullLine: true,
        });
        failedSuggestion = ' failed, check conditions';
        break;
      }
    }
  }

  return (
    <Dialog open={props.dialogSuggestionOpen} onClose={props.closeDialogSuggestion} maxWidth={'xl'}>
      <DialogTitle id="alert-dialog-title" className={classes.header}>
        {props.suggestion.metadata && props.suggestion.metadata.name
          ? 'Suggestion ' +
            JSON.stringify(props.suggestion.metadata.name, null, 2) +
            failedSuggestion
          : ''}
      </DialogTitle>
      <DialogContent>
        <AceEditor
          mode="json"
          theme="xcode"
          value={JSON.stringify(props.suggestion, null, 2)}
          tabSize={2}
          fontSize={13}
          width={'100%'}
          showPrintMargin={false}
          autoScrollEditorIntoView={true}
          maxLines={40}
          minLines={10}
          readOnly={true}
          setOptions={{ useWorker: false }}
          markers={markers}
        />
      </DialogContent>
    </Dialog>
  );
};

const mapStateToProps = state => {
  return {
    dialogSuggestionOpen: state[GENERAL_MODULE].dialogSuggestionOpen,
    suggestion: state[GENERAL_MODULE].suggestion,
  };
};

export default connect(mapStateToProps, { closeDialogSuggestion })(
  withStyles(styles)(SuggestionInfoDialog),
);
