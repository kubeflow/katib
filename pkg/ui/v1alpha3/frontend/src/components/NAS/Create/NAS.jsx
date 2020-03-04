import React from 'react';
import makeStyles from '@material-ui/styles/makeStyles';
import { Tabs } from 'antd';
import 'antd/dist/antd.css';

import YAML from './YAML';
import NASParameters from './NASParameters';

const TabPane = Tabs.TabPane;

const useStyles = makeStyles({
  root: {
    flexGrow: 1,
    marginTop: 40,
  },
});

const NAS = props => {
  const classes = useStyles();

  return (
    <div className={classes.root}>
      <Tabs defaultActiveKey="1">
        <TabPane tab="YAML File" key="1">
          <YAML />
        </TabPane>
        <TabPane tab="Parameters" key="2">
          <NASParameters />
        </TabPane>
      </Tabs>
    </div>
  );
};

export default NAS;
