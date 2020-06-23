import { combineReducers } from 'redux';
import generalReducer from './general';
import nasCreateReducer from './nasCreate';
import nasMonitorReducer from './nasMonitor';
import hpCreateReducer from './hpCreate';
import templateReducer from './template';
import hpMonitorReducer from './hpMonitor';

import {
  GENERAL_MODULE,
  HP_CREATE_MODULE,
  HP_MONITOR_MODULE,
  NAS_CREATE_MODULE,
  NAS_MONITOR_MODULE,
  TEMPLATE_MODULE,
} from '../constants/constants';

const rootReducer = combineReducers({
  [GENERAL_MODULE]: generalReducer,
  [HP_CREATE_MODULE]: hpCreateReducer,
  [HP_MONITOR_MODULE]: hpMonitorReducer,
  [NAS_CREATE_MODULE]: nasCreateReducer,
  [NAS_MONITOR_MODULE]: nasMonitorReducer,
  [TEMPLATE_MODULE]: templateReducer,
});

export default rootReducer;
