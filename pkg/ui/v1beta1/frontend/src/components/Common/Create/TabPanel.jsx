import React from 'react';

import { withStyles, makeStyles } from '@material-ui/core/styles';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

import HPYAML from '../../HP/Create/YAML';
import HPParameters from '../../HP/Create/HPParameters';
import NASYAML from '../../NAS/Create/YAML';
import NASParameters from '../../NAS/Create/NASParameters';

import * as constants from '../../../constants/constants';

const useStyles = makeStyles({
  root: {
    marginTop: 40,
  },
});

const MyTabs = withStyles({
  root: {
    borderBottom: '1px solid #e8e8e8',
    marginBottom: 15,
  },
  indicator: {
    backgroundColor: '#1890ff',
  },
})(Tabs);

const MyTab = withStyles(theme => ({
  root: {
    textTransform: 'none',
    marginRight: 40,
    minWidth: 40,
    fontWeight: theme.typography.fontWeightRegular,
    fontSize: 14,
    opacity: 1,
    '&:hover': {
      color: '#40a9ff',
    },
    '&$selected': {
      color: '#1890ff',
      fontWeight: theme.typography.fontWeightMedium,
    },
    '&:focus': {
      color: '#1890ff',
    },
  },
  selected: {},
}))(props => <Tab disableRipple {...props} />);

const TabsPanel = props => {
  const [tabIndex, setTabIndex] = React.useState(0);

  const onTabChange = (event, newIndex) => {
    setTabIndex(newIndex);
  };
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <MyTabs value={tabIndex} onChange={onTabChange}>
        <MyTab label="YAML File" />
        <MyTab label="Parameters" />
      </MyTabs>
      {props.match.path === constants.LINK_HP_CREATE ? (
        tabIndex === 0 ? (
          <HPYAML />
        ) : (
          <HPParameters />
        )
      ) : tabIndex === 0 ? (
        <NASYAML />
      ) : (
        <NASParameters />
      )}
    </div>
  );
};

export default TabsPanel;
