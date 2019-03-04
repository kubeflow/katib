import { combineReducers } from 'redux'
import generalReducer from './general';
// import nasCreateReducer from './nasCreate';
import nasCreateReducer from './nasCreate';
import nasMonitorReducer from './nasMonitor';
import hpCreateReducer from './hpCreate';
// import hpMonitorReducer from './hpMonitor';
import templateReducer from './template';
import hpMonitorReducer from './hpMonitor';

const rootReducer = combineReducers({
    ["general"]: generalReducer,
    ["template"]: templateReducer,
    ["hpCreate"]: hpCreateReducer,
    ["hpMonitor"]: hpMonitorReducer,
    ["nasCreate"]: nasCreateReducer,
    ["nasMonitor"]: nasMonitorReducer,
})

export default rootReducer;