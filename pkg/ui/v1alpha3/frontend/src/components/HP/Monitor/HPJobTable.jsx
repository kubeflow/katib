import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';

import { connect } from 'react-redux';
import { fetchHPJobTrialInfo } from '../../../actions/hpMonitorActions';

const module = "hpMonitor";

const styles = theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing.unit * 3,
    overflowX: 'auto',
  },
  table: {
    minWidth: 700
  },
  hover: {
    '&:hover': {
      cursor: "pointer",
    }
  },
  created: {
    color: theme.colors.created,
  },
  running: {
      color: theme.colors.running,
  },
  succeeded: {
      color: theme.colors.succeeded,
  },
  killed: {
      color: theme.colors.killed
  },
  failed: {
      color: theme.colors.failed,
  }
});

class HPJobTable extends React.Component {

  fetchAndOpenDialogTrial = (trialName) => (event) => {
    this.props.fetchHPJobTrialInfo(trialName);
  }

  render () {
    const { classes } = this.props;

    let header = [];
    let data = [];
    if (this.props.jobData && this.props.jobData.length > 1) {
      header = this.props.jobData[0];
      data = this.props.jobData.slice(1).sort(function(a,b) {
        if(a[1] < b[1]) { return -1; }
        if(a[1] > b[1]) { return 1; }
        return 0;
      });
    }
    return (
      <Paper className={classes.root}>
        {this.props.jobData.length > 1 &&
          <Table className={classes.table}>
            <TableHead>
              <TableRow>
                {header.map(header => (
                  <TableCell>{header}</TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {data.map((row, id) => (
                <TableRow key={id}>
                  {row.map((element, index) => {
                    if (index === 0 && row[1] == "Succeeded") {
                      return (
                        <TableCell className={classes.hover} component="th" scope="row" onClick={this.fetchAndOpenDialogTrial(element)} key={index}>
                          {element}
                        </TableCell>
                      )
                    } else if (index === 1) {
                      if (element === "Created") {
                        return (
                          <TableCell className={classes.created}>{element}</TableCell>
                        )
                      } else if (element === "Running") {
                        return (
                          <TableCell className={classes.running}>{element}</TableCell>
                        )
                      } else if (element === "Succeeded") {
                        return (
                          <TableCell className={classes.succeeded}>{element}</TableCell>
                        )
                      } else if (element === "Killed") {
                        return (
                          <TableCell className={classes.killed}>{element}</TableCell>
                        )
                      } else if (element === "Failed") {
                        return (
                          <TableCell className={classes.failed}>{element}</TableCell>
                        )
                      }
                    } else {
                      return (
                        <TableCell>{element}</TableCell>
                      )
                    }
                  })}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        }
      </Paper>
    );
  }
}


const mapStateToProps = (state) => ({
  jobData: state[module].jobData,
})

export default connect(mapStateToProps, { fetchHPJobTrialInfo })(withStyles(styles)(HPJobTable));
