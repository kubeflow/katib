import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';
import Snackbar from '@material-ui/core/Snackbar';
import IconButton from '@material-ui/core/IconButton';

import CloseIcon from '@material-ui/icons/Close';

import { connect } from 'react-redux';
import { closeSnackbar } from '../../actions/generalActions';

const module = 'general';

const useStyles = makeStyles({
  close: {
    padding: 4,
  },
});

const Snack = props => {
  const classes = useStyles();

  const vertical = 'top';
  const horizontal = 'center';
  return (
    <Snackbar
      anchorOrigin={{
        vertical: vertical,
        horizontal: horizontal,
      }}
      open={props.snackOpen}
      autoHideDuration={6000}
      onClose={props.closeSnackbar}
      ContentProps={{
        'aria-describedby': 'message-id',
      }}
      message={<span id="message-id">{props.snackText}</span>}
      action={[
        <IconButton
          key="close"
          aria-label="Close"
          color="inherit"
          className={classes.close}
          onClick={props.closeSnackbar}
        >
          <CloseIcon />
        </IconButton>,
      ]}
    />
  );
};

const mapStateToProps = state => {
  return {
    snackText: state[module].snackText,
    snackOpen: state[module].snackOpen,
  };
};

export default connect(mapStateToProps, { closeSnackbar })(Snack);
