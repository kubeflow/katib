import React from 'react';
import { connect } from 'react-redux';

import { withStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';

import { filterTemplates } from '../../../actions/templateActions';

import { GENERAL_MODULE, TEMPLATE_MODULE } from '../../../constants/constants';

const styles = theme => ({
  selectBox: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
    width: 200,
    height: 56,
  },
  textField: {
    marginLeft: theme.spacing.unit,
    marginRight: theme.spacing.unit,
  },
  selectLabel: {
    marginLeft: '8px',
  },
});

class FilterPanel extends React.Component {
  componentDidMount() {
    this.props.filterTemplates(
      this.props.filteredConfigMapNamespace,
      this.props.filteredConfigMapName,
    );
  }

  onConfigMapNamespaceChange = event => {
    this.props.filterTemplates(event.target.value, this.props.filteredConfigMapName);
  };

  onConfigMapNameChange = event => {
    this.props.filterTemplates(this.props.filteredConfigMapNamespace, event.target.value);
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        <FormControl variant="outlined">
          <InputLabel className={classes.selectLabel}>ConfigMap Namespace</InputLabel>
          <Select
            value={this.props.filteredConfigMapNamespace}
            onChange={this.onConfigMapNamespaceChange}
            className={classes.selectBox}
            label="ConfigMap Namespace"
          >
            {this.props.namespaces.map((namespace, i) => {
              return (
                <MenuItem value={namespace} key={i}>
                  {namespace}
                </MenuItem>
              );
            })}
          </Select>
        </FormControl>
        <TextField
          label="ConfigMap Name"
          className={classes.textField}
          value={this.props.filteredConfigMapName}
          onChange={this.onConfigMapNameChange}
          variant="outlined"
        />
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    namespaces: state[GENERAL_MODULE].namespaces,
    filteredConfigMapNamespace: state[TEMPLATE_MODULE].filteredConfigMapNamespace,
    filteredConfigMapName: state[TEMPLATE_MODULE].filteredConfigMapName,
  };
};

export default connect(mapStateToProps, { filterTemplates })(withStyles(styles)(FilterPanel));
