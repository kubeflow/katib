import React from 'react';

import Divider from '@material-ui/core/Divider';

import TrialConfigMap from './TrialConfigMap';
import TrialParameters from './TrialParameters';

class TrialTemplate extends React.Component {
  render() {
    return (
      <div>
        <TrialConfigMap></TrialConfigMap>
        <Divider />
        <TrialParameters></TrialParameters>
      </div>
    );
  }
}

export default TrialTemplate;
