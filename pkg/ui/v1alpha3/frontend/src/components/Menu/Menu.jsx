import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import makeStyles from '@material-ui/styles/makeStyles';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import Typography from '@material-ui/core/Typography';
import Divider from '@material-ui/core/Divider';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Collapse from '@material-ui/core/Collapse';
import TuneIcon from '@material-ui/icons/Tune';
import NoteAddIcon from '@material-ui/icons/NoteAdd';
import WatchLaterIcon from '@material-ui/icons/WatchLater';
import SearchIcon from '@material-ui/icons/Search';
import SettingsIcon from '@material-ui/icons/Settings';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import InfoIcon from '@material-ui/icons/Info';

import { toggleMenu } from '../../actions/generalActions';

const module = 'general';

const useStyles = makeStyles({
  list: {
    width: 250,
  },
  nested: {
    paddingLeft: 10 * 4,
  },
});

const Menu = props => {
  const [hp, setHP] = React.useState(false);
  const [nas, setNAS] = React.useState(false);

  const toggleHP = () => {
    setHP(!hp);
  };

  const toggleNAS = () => {
    setNAS(!nas);
  };

  const classes = useStyles();

  const onCloseMenu = () => {
    props.toggleMenu(false);
  };

  // Add links
  const color = 'primary';
  const iconColor = 'primary';
  const variant = 'title';
  return (
    <div>
      <Drawer open={props.menuOpen} onClose={onCloseMenu}>
        <List>
          {/* HP */}
          <ListItem button onClick={toggleHP}>
            <ListItemIcon>
              <TuneIcon color={iconColor} />
            </ListItemIcon>
            <ListItemText>
              <Typography variant={variant} color={color}>
                HP
              </Typography>
            </ListItemText>
            {hp ? <ExpandLess /> : <ExpandMore />}
          </ListItem>
          <Collapse in={hp} timeout="auto" unmountOnExit>
            <List component="div" disablePadding>
              <ListItem
                button
                className={classes.nested}
                component={Link}
                to="/katib/hp"
                onClick={onCloseMenu}
              >
                <ListItemIcon>
                  <NoteAddIcon color={iconColor} />
                </ListItemIcon>
                <ListItemText>
                  <Typography variant={variant} color={color}>
                    Submit
                  </Typography>
                </ListItemText>
              </ListItem>
              <ListItem
                button
                className={classes.nested}
                component={Link}
                to="/katib/hp_monitor"
                onClick={onCloseMenu}
              >
                <ListItemIcon>
                  <WatchLaterIcon color={iconColor} />
                </ListItemIcon>
                <ListItemText>
                  <Typography variant={variant} color={color}>
                    Monitor
                  </Typography>
                </ListItemText>
              </ListItem>
            </List>
          </Collapse>
          <Divider />
          {/* NAS */}
          <ListItem button onClick={toggleNAS}>
            <ListItemIcon>
              <SearchIcon color={iconColor} />
            </ListItemIcon>
            <ListItemText>
              <Typography variant={variant} color={color}>
                NAS
              </Typography>
            </ListItemText>
            {hp ? <ExpandLess /> : <ExpandMore />}
          </ListItem>
          <Collapse in={nas} timeout="auto" unmountOnExit>
            <List component="div" disablePadding>
              <ListItem
                button
                className={classes.nested}
                component={Link}
                to="/katib/nas"
                onClick={onCloseMenu}
              >
                <ListItemIcon>
                  <NoteAddIcon color={iconColor} />
                </ListItemIcon>
                <ListItemText>
                  <Typography variant={variant} color={color}>
                    Submit
                  </Typography>
                </ListItemText>
              </ListItem>
              <ListItem
                button
                className={classes.nested}
                component={Link}
                to="/katib/nas_monitor"
                onClick={onCloseMenu}
              >
                <ListItemIcon>
                  <WatchLaterIcon color={iconColor} />
                </ListItemIcon>
                <ListItemText>
                  <Typography variant={variant} color={color}>
                    Monitor
                  </Typography>
                </ListItemText>
              </ListItem>
            </List>
          </Collapse>
          <Divider />
          {/* TRIAL MANIFESTS */}
          <ListItem button component={Link} to="/katib/trial" onClick={onCloseMenu}>
            <ListItemIcon>
              <SettingsIcon color={iconColor} />
            </ListItemIcon>
            <ListItemText>
              <Typography variant={variant} color={color}>
                Trial Manifests
              </Typography>
            </ListItemText>
          </ListItem>
          <Divider />
        </List>
      </Drawer>
    </div>
  );
};

const mapStateToProps = state => {
  return {
    menuOpen: state[module].menuOpen,
  };
};

export default connect(mapStateToProps, { toggleMenu })(Menu);
