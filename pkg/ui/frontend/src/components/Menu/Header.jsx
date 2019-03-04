import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';

import Menu from './Menu';

import { connect } from 'react-redux';
import { toggleMenu } from '../../actions/generalActions';


const useStyles = makeStyles({
    root: {
        flexGrow: 1,
    },
    grow: {
        flexGrow: 1,
    },
    menuButton: {
        marginLeft: -12,
        marginRight: 20,
    },
});

const Header = (props) => {
    const classes = useStyles();

    const toggleMenu = (event) => {
        props.toggleMenu(true);
    }

    return (
        <div className={classes.root}>
            <AppBar position={"static"} color={"primary"}>
                <Toolbar>
                    <IconButton className={classes.menuButton} color={"inherit"} aria-label={"Menu"} onClick={toggleMenu}>
                        <MenuIcon/>
                    </IconButton>
                    <Typography variant={"headline"} color={"secondary"}>
                        Katib
                    </Typography>
                </Toolbar>
                <Menu />
            </AppBar>
        </div>
    )
}

export default connect(null, { toggleMenu }, )(Header);