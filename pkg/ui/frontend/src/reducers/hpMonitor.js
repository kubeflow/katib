
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
};

const hpMonitorReducer = (state = initialState, action) => {
    switch (action.type) {
        default:
            return state;
    }
};

export default hpMonitorReducer;