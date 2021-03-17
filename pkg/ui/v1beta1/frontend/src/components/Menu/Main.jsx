import React from 'react';

import { Link } from 'react-router-dom';

import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';

import { LINK_HP_CREATE, LINK_NAS_CREATE } from '../../constants/constants';

const useStyles = makeStyles({
  root: {
    margin: '0 auto',
    marginTop: 50,
    flexGrow: 1,
    width: '50%',
    height: 400,
    textAlign: 'center',
  },
  item: {
    padding: '40px !important',
    textDecoration: 'none !important',
  },
  block: {
    backgroundColor: '#4e4e4e',
    height: '100%',
    width: '100%',
    padding: 40,
    '&:hover': {
      backgroundColor: 'black',
    },
  },
  link: {
    textDecoration: 'none',
    color: '#1890ff',
  },
});

const Main = props => {
  const classes = useStyles();

  return (
    <Paper elevation={4} className={classes.root}>
      <Typography variant={'h4'}>Welcome to Katib</Typography>
      <Typography variant={'h6'}>Choose type of experiment</Typography>
      <br />
      <Grid container spacing={5} alignContent={'center'}>
        <Grid item xs={6} className={classes.item} component={Link} to={LINK_HP_CREATE}>
          <Paper className={classes.block}>
            <Typography variant={'h6'} color={'secondary'}>
              Hyperparameter Tuning
            </Typography>
          </Paper>
        </Grid>
        <Grid item xs={6} className={classes.item} component={Link} to={LINK_NAS_CREATE}>
          <Paper className={classes.block}>
            <Typography variant={'h6'} color={'secondary'}>
              Neural Architecture Search
            </Typography>
          </Paper>
        </Grid>
      </Grid>
      <br />
      <Typography variant={'h6'}>
        For usage instructions, see the{' '}
        <a
          href="https://www.kubeflow.org/docs/components/katib/"
          target="_blank"
          rel="noopener noreferrer"
          className={classes.link}
        >
          Kubeflow docs
        </a>
      </Typography>
      <Typography variant={'h6'}>
        To contribute to Katib, visit{' '}
        <a
          href="https://github.com/kubeflow/katib/"
          rel="noopener noreferrer"
          target="_blank"
          className={classes.link}
        >
          GitHub
        </a>
      </Typography>
    </Paper>
  );
};

export default Main;
