import http from 'k6/http';
import { check, fail, sleep } from 'k6';

// Init stage
export const options = {
    // A number specifying the number of VUs to run concurrently.
    vus: 10,
    // A string specifying the total duration of the test run.
    duration: '30s',

    // Customize the statistics to include in the summary
    summaryTrendStats: ['min', 'max', 'mean', 'p(90)', 'p(95)'],

    // Gradually ramp up or ramp down the number of virtual users (VUs) over time.
    stages: [
        { duration: '30s', target: 20 }, // Ramp-up to 20 VUs over 30 seconds
        { duration: '1m', target: 20 }, // Stay at 20 VUs for 1 minute
        { duration: '10s', target: 50 }, // Ramp-up to 50 VUs over 10 seconds
        { duration: '1m', target: 50 }, // Stay at 50 VUs for 1 minute
        { duration: '30s', target: 0 }, // Ramp-down to 0 VUs over 30 seconds
    ],
};

// (Optional) Setup stage, preparing data before execute test
// The setup function runs once, before any VUs start their execution.
// The return value from the setup function is passed to the default function of each VU.
export function setup() {
    let totalUser = __ENV.TOTAL_USER || 10; // Access the environment variable of TOTAL_USER

    let mockData = { id: 1, title: 'foo' };
    return mockData;
}

// Execution stage
export default function (mockData) {
    // Define the URL to request
    let url = 'https://jsonplaceholder.typicode.com/posts';

    // Define the parameters, including headers if needed
    let params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    // Define the payload for a POST request
    let payload = JSON.stringify({
        title: 'foo',
        body: 'bar',
        userId: 1,
    });

    // Make a POST request
    // https://k6.io/docs/javascript-api/k6-http/
    let res = http.post(url, payload, params);

    // Check if the response status is 200 (OK)
    let success = check(res, {
        'status is 200': (res) => res.status === 200,
        'response time is less than 1000ms': (r) => r.timings.duration < 1000,
    });

    // If the check fails, trigger a failure
    if (!success) {
        fail(`Request failed with status ${res.status}`); // This will return the function
    }

    // Log the response body for debugging (optional)
    console.log(res.body);

    sleep(1);
}

// (Optional) Teardown Stage
export function teardown(data) {
    // Perform cleanup or analysis
    console.log('Test completed, cleaning up...');
}
