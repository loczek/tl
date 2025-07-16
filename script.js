import http from "k6/http";

export let options = {
  scenarios: {
    variable_rps: {
      executor: "ramping-arrival-rate",
      startRate: 5,
      preAllocatedVUs: 50,
      maxVUs: 100,
      stages: [
        { target: 5, duration: "2m" },
        { target: 15, duration: "4m" },
        { target: 5, duration: "5m" },
        { target: 25, duration: "3m" },
        { target: 10, duration: "1m" },
        { target: 2, duration: "1m" },
        { target: 20, duration: "5m" },
        { target: 5, duration: "2m" },
      ],
    },
  },
};

export default function () {
  http.get("http://localhost:3000/pr9Rlg8");
}
