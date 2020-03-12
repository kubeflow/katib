import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/core/styles';

import TextField from '@material-ui/core/TextField';
import Select from '@material-ui/core/Select';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import FormControl from '@material-ui/core/FormControl';
import FormGroup from '@material-ui/core/FormGroup';

import { fetchNamespaces } from '../../../actions/generalActions';

const module = 'template';
const generalModule = 'general';

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
});

//TODO: Enable Filter Pannel
class FilterPanel extends React.Component {
  componentDidMount() {
    this.props.fetchNamespaces();
    // this.props.filterTemplates(this.props.templatesNamespace, this.props.templatesConfigMapName);
  }

  onNamespaceChange = event => {
    // this.props.filterTemplates(event.target.value, this.props.templatesNamespace);
  };

  onConfigMapNameChange = event => {
    // this.props.filterTemplates(this.props.templatesConfigMapName, event.target.value);
  };

  render() {
    const { classes } = this.props;

    return (
      <div>
        <FormControl variant="outlined">
          <InputLabel>Namespace</InputLabel>
          <Select
            value={this.props.templatesNamespace}
            onChange={this.onNamespaceChange}
            className={classes.selectBox}
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
          value={this.props.templatesConfigMapName}
          onChange={this.onNameChange}
          margin="normal"
          variant="outlined"
        />
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    namespaces: state[generalModule].namespaces,
    templatesNamespace: state[module].templatesNamespace,
    templatesConfigMapName: state[module].templatesConfigMapName,
    templatesConfigMapsList: state[module].templatesConfigMapsList,
  };
};

export default connect(mapStateToProps, { fetchNamespaces })(withStyles(styles)(FilterPanel));
