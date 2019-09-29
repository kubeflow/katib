import React  from 'react';

import Header from './Menu/Header';
import Snack from './Menu/Snack';

import { makeStyles } from '@material-ui/styles';

import  { Route } from 'react-router-dom';
import Main from './Menu/Main';
import HP from './HP/Create/HP';
import HPJobMonitor from './HP/Monitor/HPJobMonitor';
import HPJobInfo from './HP/Monitor/HPJobInfo';
import NAS from './NAS/Create/NAS';
import NASJobMonitor from './NAS/Monitor/NASJobMonitor';
import NASJobInfo from './NAS/Monitor/NASJobInfo';
import Trial from './Templates/Trial';


const useStyles = makeStyles({
    root: {
        width: '90%',
        margin: '0 auto',
        paddingTop: 20,
    }
});

const App = (props) => { 
    const classes = useStyles();
    return (
        <div className={classes.root}>
            <Header />
            <Route exact path="/" component={Main} />
            <Route path="/katib/hp" component={HP} />
            <Route exact path="/katib/hp_monitor" component={HPJobMonitor} />
            <Route path="/katib/hp_monitor/:namespace/:name" component={HPJobInfo} />
            <Route path="/katib/nas" component={NAS} />
            <Route exact path="/katib/nas_monitor" component={NASJobMonitor} />
            <Route path="/katib/nas_monitor/:name" component={NASJobInfo} />
            <Route path="/katib/trial" component={Trial} />
            <Snack />
        </div>
    )
};

export default App;
