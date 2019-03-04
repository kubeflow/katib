
const initialState = {
    filter: '',
    filterType: {
        "Running": true,
        "Failed": true,
        "Succeeded": true,
    },
    jobsList: [
        {
            name: "Job 1",
            status: "Running",
            id: "1", 
        },
        {
            name: "Job 2",
            status: "Failed",
            id: "2", 
        },
        {
            name: "Job 3",
            status: "Succeded",
            id: "3", 
        }
    ],
    filteredJobsList: [
        {
            name: "Job 1",
            status: "Running",
            id: "1", 
        },
        {
            name: "Job 2",
            status: "Failed",
            id: "2", 
        },
        {
            name: "Job 3",
            status: "Succeded",
            id: "3", 
        }
    ],
    steps: [
        {
            id: "1",
            name: "Generation 1",
            architecture: "ASD",
            metricsName: "Accuracy",
            metricsValue: "0.99",
            link: "link",
        },
        {
            id: "2",
            name: "Generation 2",
            architecture: "DAS",
            metricsName: "Accuracy",
            metricsValue: "0.99999",
            link: "link",
        },
        {
            id: "3",
            name: "Generation 3",
            architecture: "DAS",
            metricsName: "Accuracy",
            metricsValue: "0.999999999",
            link: "link",
        },
    ]
};

const nasMonitorReducer = (state = initialState, action) => {
    switch (action.type) {
        default:
            return state;
    }
};

export default nasMonitorReducer;