import React from 'react';
import { connect } from 'react-redux';

import { Link } from 'react-router-dom';

import { makeStyles, withStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';

import Menu from './Menu';

import { toggleMenu } from '../../actions/generalActions';

const useStyles = makeStyles({
  menuButton: {
    marginLeft: -12,
    marginRight: 20,
  },
});

const KatibLink = withStyles({
  root: {
    textDecoration: 'none',
    '&:hover': {
      color: '#40a9ff',
    },
  },
})(Typography);

const Header = props => {
  const classes = useStyles();

  const toggleMenu = event => {
    props.toggleMenu(true);
  };

  return (
    <div>
      <AppBar position={'static'} color={'primary'}>
        <Toolbar>
          <IconButton
            className={classes.menuButton}
            color={'inherit'}
            aria-label={'Menu'}
            onClick={toggleMenu}
          >
            <MenuIcon />
          </IconButton>
          <KatibLink variant={'h5'} color={'secondary'} component={Link} to="/">
            Katib
          </KatibLink>
        </Toolbar>
        <Menu />
      </AppBar>
    </div>
  );
};

export default connect(null, { toggleMenu })(Header);
