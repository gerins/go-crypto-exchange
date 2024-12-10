import http from 'k6/http';
import { check, fail, sleep } from 'k6';

// Init stage
export const options = {
    // A number specifying the number of VUs to run concurrently.
    vus: 2000,

    // A string specifying the total duration of the test run.
    // Cannot run simultaneously with stages
    duration: '1m',

    // Customize the statistics to include in the summary
    summaryTrendStats: ['min', 'max', 'p(90)', 'p(95)'],

    // Gradually ramp up or ramp down the number of virtual users (VUs) over time.
    // stages: [
    //     { duration: '30s', target: 20 }, // Ramp-up to 20 VUs over 30 seconds
    //     { duration: '1m', target: 20 }, // Stay at 20 VUs for 1 minute
    //     { duration: '10s', target: 50 }, // Ramp-up to 50 VUs over 10 seconds
    //     { duration: '1m', target: 50 }, // Stay at 50 VUs for 1 minute
    //     { duration: '30s', target: 0 }, // Ramp-down to 0 VUs over 30 seconds
    // ],
};

// (Optional) Setup stage, preparing data before execute test
// The setup function runs once, before any VUs start their execution.
// The return value from the setup function is passed to the default function of each VU.

// Load the CSV file
const emailList = open('../emails.csv'); // Ensure this 'open' function exists and works in your environment

// Parse CSV content into an array of emails
let userEmails = emailList
    .split('\n') // Split the content by newlines
    .map((line) => line.trim()) // Remove any extra spaces
    .filter((line) => line); // Exclude empty lines

export function setup() {
    let listToken = [];

    userEmails.forEach((email) => {
        let params = {
            headers: {
                'Content-Type': 'application/json',
            },
        };
        let payload = JSON.stringify({
            email: email,
            password: 'admin',
        });

        let res = http.post('http://localhost:8070/api/v1/user/login', payload, params);

        // Parse the response body as JSON
        let responseBody = res.json();

        // Ensure the response body contains the expected 'data.token'
        if (responseBody && responseBody.data && responseBody.data.token) {
            listToken.push(responseBody.data.token); // Add the token to the listToken array
        } else {
            console.log(`Login failed for user: ${email}, response: ${JSON.stringify(responseBody)}`);
        }
    });

    return listToken;
}

// Execution stage
export default function (listToken) {
    // Define the URL to request
    let url = 'http://localhost:8070/api/v1/order';

    // Define the parameters, including headers if needed
    let params = {
        headers: {
            'Content-Type': 'application/json',
            Authorization: 'Bearer ' + listToken[getRandomInt(0, 5)],
        },
    };

    // Define the payload for a POST request
    let payload = JSON.stringify({
        pair_code: 'DOGEIDRT',
        quantity: getRandomInt(1, 100),
        price: getRandomInt(100, 1000),
        side: getRandomSide(),
        type: 'LIMIT',
    });

    // Make a POST request
    // https://k6.io/docs/javascript-api/k6-http/
    let res = http.post(url, payload, params);

    // Check if the response status is 200 (OK)
    let success = check(res, {
        'status is 200': (res) => res.status === 200,
        'response time is less than 3000ms': (r) => r.timings.duration < 3000,
    });

    // If the check fails, trigger a failure
    if (!success) {
        fail(`Request failed with status ${res.status}`); // This will return the function
    }

    sleep(1);
}

// (Optional) Teardown Stage
export function teardown(data) {
    // Perform cleanup or analysis
    console.log('Test completed, cleaning up...');
}

function getRandomInt(min, max) {
    min = Math.ceil(min);
    max = Math.floor(max);
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function getRandomSide() {
    return Math.random() < 0.5 ? 'BUY' : 'SELL';
}
