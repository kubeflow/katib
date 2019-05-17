import React from 'react';

import PropTypes from 'prop-types';
import withStyles from '@material-ui/styles/withStyles';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';

import CommonParametersMeta from './Params/CommonMeta';
import CommonParametersSpec from './Params/CommonSpec';
import Objective from "./Params/Objective"
import Algorithm from "./Params/Algorithm"
import TrialSpecParam from './Params/Trial';

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
        let name = parameter.name.charAt(0).toLowerCase() + parameter.name.slice(1)
        destination[name] = (isNaN(value) ? parameter.value : value)
    })
}

const addAlgorithmSettings = (spec, destination) => {
    spec.map((parameter, i) => {
        destination.push(parameter)
    })
}

const addOperations = (source, destination) => {
    source.map((operation, index) => {
        let parameters = []
        operation.parameters.map((param, i) => {
            let  tempParam = {}
            tempParam.name = param.name
            tempParam.parametertype = param.parameterType
            tempParam.feasibleSpace = {}
            if (param.feasibleSpace === "list") {
                tempParam.feasibleSpace.list = param.list.map((param, i) => param.value)
            } else {
                tempParam.feasibleSpace.min = param.min
                tempParam.feasibleSpace.max = param.max
                tempParam.feasibleSpace.step = param.step
            }
            parameters.push(tempParam)
        })
        destination.push({
            operationType: operation.operationType,
            parameters: parameters,
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
        
        data.spec.objective = {}
        deCapitalizeFirstLetterAndAppend(props.objective, data.spec.objective)
        data.spec.objective.additionalMetricNames = props.additionalMetricNames.map((metrics, i) => metrics.value)
        
        data.spec.algorithm = {}
        data.spec.algorithm.algorithmName = props.algorithmName
        data.spec.algorithm.algorithmSettings = [] 
        addAlgorithmSettings(props.algorithmSettings, data.spec.algorithm.algorithmSettings)
       
        data.spec.nasConfig = {}
        data.spec.nasConfig.graphConfig = {
            numLayers: Number(props.numLayers),
            inputSize: props.inputSize.map(size => Number(size)),
            outputSize: props.outputSize.map(size => Number(size)),
        }
        data.spec.nasConfig.operations = []
        addOperations(props.operations, data.spec.nasConfig.operations)

        data.spec.trialTemplate = {
            goTemplate: {
                templatePath: props.trial,
            }
        }

        console.log("DATA BEFORE BACKEND")
        console.log(data)
         props.submitNASJob(data)
    }

    const { classes } = props;
    
    return (
            <div className={classes.root}>
                {/* Common Metadata */}
                {SectionInTypography("Metadata", classes)}
                <CommonParametersMeta />
                {SectionInTypography("Common Parameters", classes)}
                <CommonParametersSpec />
                {SectionInTypography("Objective", classes)}
                <Objective />
                {SectionInTypography("Algorithm", classes)}
                <Algorithm/>
                {SectionInTypography("NAS Config", classes)}
                <NASConfig />
                {SectionInTypography("Trial Spec", classes)}
                <TrialSpecParam />
                <div className={classes.submit}>
                    <Button variant="contained" color={"primary"} className={classes.button} onClick={submitNASJob}>
                        Deploy
                    </Button>
                </div>        
            </div>
    )
}

// TODO: think of a better way of passing those
const mapStateToProps = (state) => ({
    commonParametersMetadata: state[module].commonParametersMetadata,
    commonParametersSpec: state[module].commonParametersSpec,
    objective: state[module].objective,
    additionalMetricNames: state[module].additionalMetricNames,
    algorithmName: state[module].algorithmName,
    algorithmSettings: state[module].algorithmSettings,
    numLayers: state[module].numLayers,
    inputSize: state[module].inputSize,
    outputSize: state[module].outputSize,
    operations: state[module].operations,
    trial: state[module].trial
})

//TODO: Added validation and remove it
// NASParameters.propTypes = {
//     numLayers: PropTypes.number,
//     trial: PropTypes.string,
//     requestNumber: PropTypes.number,
//     suggestionAlgorithm: PropTypes.string,
//     metricsName: PropTypes.arrayOf(PropTypes.string),
// }

export default connect(mapStateToProps, { submitNASJob })(withStyles(styles)(NASParameters));
