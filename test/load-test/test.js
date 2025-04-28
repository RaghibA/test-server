import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 4000, // Number of virtual users to simulate
  duration: '60s', // Duration of the test
};

export default function () {
  const i = __ITER + 1; // Incrementing counter for data
  const payload = JSON.stringify({
    data: `incrementing data ${i}`,
  });
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post('https://test-service.raghiba.com/test', payload, params);
  sleep(1); 
}