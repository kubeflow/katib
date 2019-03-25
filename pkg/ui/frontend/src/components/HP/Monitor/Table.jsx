import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Paper from '@material-ui/core/Paper';

import { connect } from 'react-redux';
import { fetchWorkerInfo } from '../../../actions/hpMonitorActions';

const module = "hpMonitor";

const styles = theme => ({
  root: {
    width: '100%',
    marginTop: theme.spacing.unit * 3,
    overflowX: 'auto',
  },
  table: {
    minWidth: 700,
  },
  hover: {
    '&:hover': {
      cursor: "pointer",
    }
  }
});

class HPTable extends React.Component {

  fetchAndOpenDialog = (id) => (event) => {
    this.props.fetchWorkerInfo(this.props.id, id);
  }

  render () {
    const { classes } = this.props;

    let header = [];
    let data = [];
    if (this.props.jobData && this.props.jobData.length !== 0) {
      header = this.props.jobData[0];
      data = this.props.jobData.slice(1)
    }
    return (
      <Paper className={classes.root}>
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
                  if (index === 0) {
                    return (
                      <TableCell className={classes.hover} component="th" scope="row" onClick={this.fetchAndOpenDialog(element)} key={index}>
                        {element}
                      </TableCell>
                    )
                  } else {
                    return (
                      <TableCell align="right">{element}</TableCell>
                    )
                  }
                })}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </Paper>
    );
  }
}


const mapStateToProps = (state) => ({
  jobData: state[module].jobData,
})

export default connect(mapStateToProps, { fetchWorkerInfo })(withStyles(styles)(HPTable));