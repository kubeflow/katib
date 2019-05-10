import * as actions from '../actions/nasMonitorActions';

const initialState = {
    filter: '',
    filterType: {
        "Running": true,
        "Failed": true,
        "Completed": true,
    },
    jobsList: [
    ],
    filteredJobsList: [
    ],
    loading: true,
    steps: [
        {
          "name": "Generation 0",
          "architecture": "digraph G {\n\t0->1;\n\t1->2;\n\t0->2;\n\t2->3;\n\t1->3;\n\t3->4;\n\t0->4;\n\t1->4;\n\t4->5;\n\t2->5;\n\t3->5;\n\t5->6;\n\t6->7;\n\t0->7;\n\t1->7;\n\t4->7;\n\t7->8;\n\t1->8;\n\t2->8;\n\t3->8;\n\t5->8;\n\t8->9;\n\t9->10;\n\t10->11;\n\t0 [ label=\"Input\" ];\n\t1 [ label=\"7x7 sep_conv\\n7 channels\" ];\n\t10 [ label=\"FullConnect\\nSoftmax\" ];\n\t11 [ label=\"Output\" ];\n\t2 [ label=\"7x7 depth_conv\\n\" ];\n\t3 [ label=\"7x7 sep_conv\\n7 channels\" ];\n\t4 [ label=\"5x5 sep_conv\\n5 channels\" ];\n\t5 [ label=\"3x3 conv\\n3 channels\" ];\n\t6 [ label=\"3x3 sep_conv\\n3 channels\" ];\n\t7 [ label=\"3x3 sep_conv\\n3 channels\" ];\n\t8 [ label=\"7x7 sep_conv\\n7 channels\" ];\n\t9 [ label=\"GlobalAvgPool\" ];\n\n}\n",
          "metricsname": [],
          "metricsvalue": [],
          "trialname": "Trial 1"
        },
        {
          "name": "Generation 1",
          "architecture": "digraph G {\n\t0->1;\n\t1->2;\n\t2->3;\n\t1->3;\n\t3->4;\n\t1->4;\n\t2->4;\n\t4->5;\n\t3->5;\n\t5->6;\n\t0->6;\n\t2->6;\n\t3->6;\n\t6->7;\n\t1->7;\n\t4->7;\n\t5->7;\n\t7->8;\n\t0->8;\n\t1->8;\n\t3->8;\n\t4->8;\n\t5->8;\n\t6->8;\n\t8->9;\n\t9->10;\n\t10->11;\n\t0 [ label=\"Input\" ];\n\t1 [ label=\"7x7 sep_conv\\n7 channels\" ];\n\t10 [ label=\"FullConnect\\nSoftmax\" ];\n\t11 [ label=\"Output\" ];\n\t2 [ label=\"7x7 depth_conv\\n\" ];\n\t3 [ label=\"3x3 conv\\n3 channels\" ];\n\t4 [ label=\"7x7 sep_conv\\n7 channels\" ];\n\t5 [ label=\"5x5 sep_conv\\n5 channels\" ];\n\t6 [ label=\"3x3 sep_conv\\n3 channels\" ];\n\t7 [ label=\"5x5 conv\\n5 channels\" ];\n\t8 [ label=\"7x7 sep_conv\\n7 channels\" ];\n\t9 [ label=\"GlobalAvgPool\" ];\n\n}\n",
          "metricsname": [
            "Validation-Accuracy"
          ],
          "metricsvalue": [
            "0.6473"
          ],
          "trialname": "Trial 2"
        },
        {
          "name": "Generation 2",
          "trialname": "Trial 3",
          "architecture": "digraph G {\n\t0->1;\n\t1->2;\n\t2->3;\n\t0->3;\n\t3->4;\n\t1->4;\n\t2->4;\n\t4->5;\n\t0->5;\n\t1->5;\n\t3->5;\n\t5->6;\n\t1->6;\n\t3->6;\n\t4->6;\n\t6->7;\n\t0->7;\n\t3->7;\n\t4->7;\n\t5->7;\n\t7->8;\n\t1->8;\n\t4->8;\n\t5->8;\n\t6->8;\n\t8->9;\n\t9->10;\n\t10->11;\n\t0 [ label=\"Input\" ];\n\t1 [ label=\"5x5 sep_conv\\n5 channels\" ];\n\t10 [ label=\"FullConnect\\nSoftmax\" ];\n\t11 [ label=\"Output\" ];\n\t2 [ label=\"3x3 depth_conv\\n\" ];\n\t3 [ label=\"3x3 sep_conv\\n3 channels\" ];\n\t4 [ label=\"3x3 sep_conv\\n3 channels\" ];\n\t5 [ label=\"7x7 conv\\n7 channels\" ];\n\t6 [ label=\"3x3 sep_conv\\n3 channels\" ];\n\t7 [ label=\"5x5 sep_conv\\n5 channels\" ];\n\t8 [ label=\"7x7 depth_conv\\n\" ];\n\t9 [ label=\"GlobalAvgPool\" ];\n\n}\n",
          "metricsname": [],
          "metricsvalue": []
        }
      ]
};

const nasMonitorReducer = (state = initialState, action) => {
    switch (action.type) {
        case actions.FILTER_JOBS:
            let jobs = state.jobsList.slice();
            let newList = jobs.filter(job => job.name.includes(action.filter));

            let avTypes = Object.assign({}, state.filterType);
            var typeKeys = Object.keys(avTypes);

            var avFilters = typeKeys.filter((key) => {
                return avTypes[key]
            });

            let filteredJobs = newList.filter(job => avFilters.includes(job.status));

            return {
                ...state,
                filteredJobsList: filteredJobs,
                filter: action.filter,
            }
        case actions.CHANGE_TYPE:
            const types = Object.assign({}, state.filterType)
            types[action.filter] = action.checked;
            var keys = Object.keys(types);

            var filters = keys.filter((key) => {
                return types[key]
            });
            const jobsList = state.jobsList.slice();
            const filtered = jobsList.filter(job => filters.includes(job.status));
            
            return {
                ...state,
                filterType: types,
                filteredJobsList: filtered,
            }
        case actions.FETCH_NAS_JOBS_SUCCESS:
            jobs = action.jobs
            avTypes = Object.assign({}, state.filterType);
            typeKeys = Object.keys(avTypes);

            avFilters = typeKeys.filter((key) => {
                return avTypes[key]
            });

            filteredJobs = jobs.filter(job => avFilters.includes(job.status));
            return {
                ...state,
                jobsList: action.jobs,
                filteredJobsList: filteredJobs,
            }
        case actions.FETCH_NAS_JOB_INFO_REQUEST:
            return {
                ...state,
                loading: true,
            }
        case actions.FETCH_NAS_JOB_INFO_FAILURE:
            return {
                ...state,
                loading: false,
            }
        case actions.FETCH_NAS_JOB_INFO_SUCCESS:
            return {
                ...state,
                loading: false,
                steps: action.steps,
            }
        default:
            return state;
    }
};


export default nasMonitorReducer;