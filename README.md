# Hack the North Backend Challenge 2024

Submission for the 2024 HTN BE Challenge.

The application consists of a Go backend, PostgreSQL database, and React frontend.

You can find it hosted here ()

## About The Project

This project was my first time trying out Go for backend development. It was a lot of learning and I am excited to continue experimenting with it!
The frontend, built in React and Tailwind, provides a simple UI for testing the endpoints. I was interested in tying together the service as a full-stack web application for two main reasons.

1. Make It Easier To Test Endpoints
2. Get a More Comprehensive Experience of Go Backend Development and its Integrations

Additionally, in the scenario where more endpoints were implemented, there would be a centralized dashboard to showcase what functions were available.
That being said, you can definitely test all the endpoints with Postman, if you would prefer!

## Installation

If you would like to install the project locally.
Clone the repository

```
git clone https://github.com/KennethRuan/htn-backend-24
```

and then run docker compose up

```
docker compose up
```

## API Endpoints

- All Users Endpoint (/api/users GET)
- User Information Endpoint (/api/users/:id GET)
- Updating User Data Endpoint (/api/users/:id PUT)
- Skills Frequency Endpoint (/api/skills/?min_frequency={}&max_frequency={})
