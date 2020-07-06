import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import FormGroup from '@material-ui/core/FormGroup';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import Switch from '@material-ui/core/Switch';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';

import { fetchNamespaces, filterExperiments, changeStatus } from '../../../actions/generalActions';

import { GENERAL_MODULE } from '../../../constants/constants';

const styles = theme => ({
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  },
  filter: {
    margin: '0 auto',
    textAlign: 'center',
  },
  selectBox: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    width: 200,
    height: 56,
    textAlign: 'left',
  },
  selectLabel: {
    marginLeft: '8px',
  },
});

class FilterPanel extends React.Component {
  componentDidMount() {
    if (this.props.globalNamespace !== '') {
      this.props.filterExperiments(this.props.experimentName, this.props.globalNamespace);
    } else {
      this.props.fetchNamespaces();
      this.props.filterExperiments(this.props.experimentName, this.props.experimentNamespace);
    }
  }

  handleStatus = name => event => {
    this.props.changeStatus(name, event.target.checked);
  };

  onNameChange = event => {
    this.props.filterExperiments(event.target.value, this.props.experimentNamespace);
  };

  onNamespaceChange = event => {
    this.props.filterExperiments(this.props.experimentName, event.target.value);
  };

  render() {
    const { classes } = this.props;

    return (
      <div className={classes.filter}>
        <FormGroup row>
          <FormControl variant="outlined">
            <InputLabel className={classes.selectLabel}>Namespace</InputLabel>
            {this.props.globalNamespace === '' ? (
              <Select
                value={this.props.experimentNamespace}
                onChange={this.onNamespaceChange}
                className={classes.selectBox}
                label="Namespace"
              >
                {this.props.namespaces.map((namespace, i) => {
                  return (
                    <MenuItem value={namespace} key={i}>
                      {namespace}
                    </MenuItem>
                  );
                })}
              </Select>
            ) : (
              <Select
                value={this.props.experimentNamespace}
                className={classes.selectBox}
                disabled
                label="Namespace"
              >
                <MenuItem value={this.props.experimentNamespace}>
                  {this.props.experimentNamespace}
                </MenuItem>
              </Select>
            )}
          </FormControl>
          <TextField
            id="outlined-name"
            label="Name"
            className={classes.textField}
            value={this.props.experimentName}
            onChange={this.onNameChange}
            variant="outlined"
          />
          {Object.keys(this.props.filterStatus).map((filter, i) => {
            return (
              <FormControlLabel
                key={i}
                control={
                  <Switch
                    checked={this.props.filterStatus[filter]}
                    onChange={this.handleStatus(filter)}
                    value={filter}
                    color={'primary'}
                  />
                }
                label={filter}
              />
            );
          })}
        </FormGroup>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    experimentName: state[GENERAL_MODULE].experimentName,
    experimentNamespace: state[GENERAL_MODULE].experimentNamespace,
    filterStatus: state[GENERAL_MODULE].filterStatus,
    namespaces: state[GENERAL_MODULE].namespaces,
    globalNamespace: state[GENERAL_MODULE].globalNamespace,
  };
};

export default connect(mapStateToProps, { filterExperiments, changeStatus, fetchNamespaces })(
  withStyles(styles)(FilterPanel),
);
