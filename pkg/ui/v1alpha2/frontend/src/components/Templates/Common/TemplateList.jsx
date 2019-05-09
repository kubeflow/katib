import React from 'react';
import { withStyles } from '@material-ui/core/styles';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import Typography from '@material-ui/core/Typography';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';

import TemplatePanel from './TemplatePanel';
import EditDialog from './EditDialog';
import DeleteDialog from './DeleteDialog';

import { connect } from 'react-redux';

const module = "template";


const styles = theme => ({
    root: {
        marginTop: 40,
        width: '100%',
    },
    heading: {
        fontSize: theme.typography.pxToRem(24),
        fontWeight: theme.typography.fontWeightRegular,
    },
});

class TemplateList extends React.Component {

    componentDidMount() {
    };

    render () {
        const { classes } = this.props;
        const templates = (this.props.type === "trial" ? this.props.trialTemplates : this.props.collectorTemplates);
        return (
            <div className={classes.root}>
                {templates.map((template, index) => {
                    return (
                        <ExpansionPanel key={index}>
                            <ExpansionPanelSummary expandIcon={<ExpandMoreIcon />}>
                                <Typography className={classes.heading}>
                                    {template.name}
                                </Typography>
                            </ExpansionPanelSummary>
                            <ExpansionPanelDetails>
                                <TemplatePanel type={this.props.type} text={template.yaml} index={index} />
                            </ExpansionPanelDetails>
                            <EditDialog type={this.props.type} />
                            <DeleteDialog type={this.props.type} />
                        </ExpansionPanel>
                    )
                })}
            </div>
        )
    }
}

const mapStateToProps = (state) => {
    return {
        collectorTemplates: state[module].collectorTemplates,
        trialTemplates: state[module].trialTemplates,
    };
};

export default connect(mapStateToProps, null)(withStyles(styles)(TemplateList));
