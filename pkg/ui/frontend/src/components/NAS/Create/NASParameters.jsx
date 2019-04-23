import React from 'react';

import PropTypes from 'prop-types';
import withStyles from '@material-ui/styles/withStyles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';

import CommonParametersMeta from './Params/CommonMeta';
import CommonParametersSpec from './Params/CommonSpec';
import SuggestionSpec from './Params/SuggestionSpec';
import WorkerSpecParam from './Params/Worker';

import { submitNASJob } from '../../../actions/nasCreateActions';


import { connect } from 'react-redux';
import NASConfig from './Params/NASConfig';

const module = "nasCreate";

const styles = theme => ({
    root: {
        width: '90%',
        margin: '0 auto',
    },
    submit: {
        textAlign: 'center',
        marginTop: 10,
    },
    textField: {
        marginLeft: 4,
        marginRight: 4,
        width: '100%'
    },
    help: {
        padding: 4 / 2,
        verticalAlign: "middle",
    },
    section: {
        padding: 4,
    },
    parameter: {
        padding: 2,
    },
    formControl: {
        margin: 4,
        width: '100%',
    },
    selectEmpty: {
        marginTop: 10,
    },
    addButton: {
        margin: 10,
    }
})

const SectionInTypography = (name, classes) => {
    return (
        <div className={classes.section}>
            <Grid container>
                <Grid item xs={12} sm={12}>
                    <Typography variant="h6">
                        {name}
                    </Typography>
                <hr />
                </Grid>
            </Grid>
        </div>
    )
}

// probably get render into a function
const deCapitalizeFirstLetterAndAppend = (source, destination) => {
    source.map((parameter, i) => {
        let value = Number(parameter.value)
        destination[parameter.name.toLowerCase()] = (isNaN(value) ? parameter.value : value)
    })
}

const addSuggestionParameters = (spec, destination) => {
    spec.map((parameter, i) => {
        destination.push(parameter)
    })
}

const addOperations = (source, destination) => {
    source.map((operation, index) => {
        let parameters = []
        operation.parameterconfigs.map((param, i) => {
            let  tempParam = {}
            tempParam.name = param.name
            tempParam.parametertype = param.parameterType
            tempParam.feasible = {}
            if (param.feasible === "list") {
                tempParam.feasible.list = param.list.map((param, i) => param.value)
            } else {
                tempParam.feasible.min = param.min
                tempParam.feasible.max = param.max
            }
            parameters.push(tempParam)
        })
        destination.push({
            operationType: operation.operationType,
            parameterconfigs: parameters,
        })
    })
}

const NASParameters = (props) => {
    const submitNASJob = () => {
        let data = {}
        data.metadata = {}
        deCapitalizeFirstLetterAndAppend(props.commonParametersMetadata, data.metadata)
        data.spec = {}
        deCapitalizeFirstLetterAndAppend(props.commonParametersSpec, data.spec)
        data.spec.nasConfig = {}
        data.spec.nasConfig.graphConfig = {
            numLayers: Number(props.numLayers),
            inputSize: props.inputSize.map(size => Number(size)),
            outputSize: props.outputSize.map(size => Number(size)),
        }
        data.spec.nasConfig.operations = []
        addOperations(props.operations, data.spec.nasConfig.operations)
        data.spec.workerSpec = {
            goTemplate: {
                templatePath: props.worker,
            }
        }
        data.spec.metricsnames = props.metricsName.map((metrics, i) => metrics.value)
        data.spec.suggestionSpec = {}
        data.spec.suggestionSpec.suggestionAlgorithm = props.suggestionAlgorithm;
        data.spec.suggestionSpec.requestNumber = (!isNaN(Number(props.requestNumber)) ? Number(props.requestNumber) : 1) 
        data.spec.suggestionSpec.suggestionParameters = []
        addSuggestionParameters(props.suggestionParameters, data.spec.suggestionSpec.suggestionParameters)
        props.submitNASJob(data)
    }

    const { classes } = props;
    
    return (
            <div className={classes.root}>
                {/* Common Metadata */}
                {SectionInTypography("Metadata", classes)}
                <CommonParametersMeta />
                {SectionInTypography("Spec", classes)}
                <CommonParametersSpec />
                {SectionInTypography("NAS Config", classes)}
                <NASConfig />
                {SectionInTypography("Worker Spec", classes)}
                <WorkerSpecParam />
                {SectionInTypography("Suggestion Parameters", classes)} 
                <SuggestionSpec />
                <div className={classes.submit}>
                    <Button variant="contained" color={"primary"} className={classes.button} onClick={submitNASJob}>
                        Deploy
                    </Button>
                </div>        
            </div>
    )
}

const mapStateToProps = (state) => ({
    commonParametersMetadata: state[module].commonParametersMetadata,
    commonParametersSpec: state[module].commonParametersSpec,
    numLayers: state[module].numLayers,
    inputSize: state[module].inputSize,
    outputSize: state[module].outputSize,
    operations: state[module].operations,
    worker: state[module].worker,
    metricsName: state[module].metricsName,
    suggestionAlgorithm: state[module].suggestionAlgorithm,
    requestNumber: state[module].requestNumber,
    suggestionParameters: state[module].suggestionParameters,
})

HPParameters.propTypes = {
    numLayers: PropTypes.number,
    worker: PropTypes.string,
    requestNumber: PropTypes.number,
    suggestionAlgorithm: PropTypes.string,
    metricsName: PropTypes.arrayOf(PropTypes.string),
}

export default connect(mapStateToProps, { submitNASJob })(withStyles(styles)(NASParameters));