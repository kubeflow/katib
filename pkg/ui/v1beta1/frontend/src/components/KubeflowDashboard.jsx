import React from 'react';
import { CHANGE_GLOBAL_NAMESPACE } from '../actions/generalActions';
import generalReducer from '../reducers/general';

// TODO: Think about better way of update namespace
// Right now, after updating namespace in select list
// user must go to Kubeflow dashboard home page to update Global Namespace in Katib UI
function onGlobalNamespaceChange(namespace) {
  generalReducer(undefined, { type: CHANGE_GLOBAL_NAMESPACE, globalNamespace: namespace });
}

class KubeflowDashboard extends React.Component {
  componentDidMount() {
    window.addEventListener('DOMContentLoaded', function(event) {
      if (window.centraldashboard && window.centraldashboard.CentralDashboardEventHandler) {
        // Init method will invoke the callback with the event handler instance
        // and a boolean indicating whether the page is iframed or not
        window.centraldashboard.CentralDashboardEventHandler.init(function(cdeh) {
          // Binds a callback that gets invoked anytime the Dashboard's
          // namespace is changed
          cdeh.onNamespaceSelected = namespace => {
            onGlobalNamespaceChange(namespace);
          };
        });
      }
    });
  }

  render() {
    return <div />;
  }
}

export default KubeflowDashboard;
