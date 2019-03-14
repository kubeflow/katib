import React  from 'react';

import Header from './Menu/Header';
import Snack from './Menu/Snack';

import { makeStyles } from '@material-ui/styles';

import  { Route } from 'react-router-dom';
import HP from './HP/Create/HP';
import HPMonitor from './HP/Monitor/Monitor';
import HPJobInfo from './HP/Monitor/JobInfo';
import NAS from './NAS/Create/NAS';
import NASMonitor from './NAS/Monitor/Monitor';
import NASJobInfo from './NAS/Monitor/JobInfo';
import Worker from './Templates/Worker';
import Collector from './Templates/Collector';


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
            
            <Route path="/hp" component={HP} />
            <Route exact path="/hp_monitor" component={HPMonitor} />
            <Route path="/hp_monitor/:id" component={HPJobInfo} />
            <Route path="/nas" component={NAS} />
            <Route exact path="/nas_monitor" component={NASMonitor} />
            <Route path="/nas_monitor/:id" component={NASJobInfo} />
            <Route path="/worker" component={Worker} />
            <Route path="/collector" component={Collector} />
            {/* <Route exact path="/" component={GenerateFromYaml} />
            <Route path="/defaults/" component={GenerateFromParameters} />
            <Route path="/monitor/" component={Watch} /> */}
            <Snack />
        </div>
    )
};

export default App;