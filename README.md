# Hack the North Backend Challenge 2024

Submission for the 2024 HTN BE Challenge.

The application consists of a Go backend, PostgreSQL database, and React frontend.

## About The Project

This project was my first time trying out Go for backend development. It was a lot of learning and I am excited to continue experimenting with it!
The frontend, built in React and Tailwind, provides a simple UI for testing the endpoints.

![](https://github.com/KennethRuan/htn-backend-24/blob/main/demo.gif)

I was interested in tying together the service as a full-stack web application for two main reasons.

1. Make It Easier To Test Endpoints
2. Get a More Comprehensive Experience of Go Backend Development and its Integrations

Additionally, in the scenario where more endpoints were implemented, the dashboard would offer a centralized location to showcase what functions were available.

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

The client and server will both be available at `localhost:8080`.

## API Endpoints

- All Users Endpoint **(/api/users GET)**
- User Information Endpoint **(/api/users/:id GET)**
- Updating User Data Endpoint **(/api/users/:id PUT)**
- Skills Frequency Endpoint **(/api/skills/?min_frequency={}&max_frequency={} GET)**

## Database Schema

For this project, to try something new, I decided to forego an ORM.
Instead I managed the database migration myself and wrote raw SQL queries.

Attached below are the table schemas, which can also be found under the migrations folder.

```SQL
CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  company VARCHAR(255),
  phone VARCHAR(255)
);

CREATE TABLE skills (
  id UUID PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE users_skills (
  user_id UUID REFERENCES users(id),
  skill_id UUID REFERENCES skills(id),
  rating INTEGER NOT NULL,
  PRIMARY KEY (user_id, skill_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
);
```

### Users_Skills Junction

My choice to represent the user skills using a junction table, amongst other options, was driven by a couple of thought processes.
For example, it would have been possible to represent skills as a `jsonb` column under the `users` table, however, I believe the junction is a more flexible representation.

- Junction tables are quite standard for representing many-to-many relationships, allowing us to make efficient and flexible queries.

- Additionally, they give us flexibility and scalability. We can adjust skills and their relationship to users as we see fit.

- It also enforces consistency and transparency in our data. If we were to update the information associated with skills, in a `jsonb` solution, it would be hard to tell if all of the data had been migrated and if any existing code is still writing data in the old format. Having a separate table allows changes to be transparent to other developers and forces us to adhere to proper migration practices.

### Indexes

On my table, I have also made two indexes.

```SQL
CREATE INDEX idx_users_skills_user_id ON users_skills(user_id);
CREATE INDEX idx_users_skills_skill_id ON users_skills(skill_id);
```

These allow for querying and joins on those columns to be performed more efficiently at the cost of a slightly slower write times.
Since we will be frequently joining the `users_skills` junction to the `users` and `skills`, I felt this was a beneficial tradeoff.

### Migrations

Migrations are done with a Go tool called [goose](https://github.com/pressly/goose).
Although, we do not have many migrations for this project, I experimented with it to demonstrate how we can make this backend implementation more extensible and developer-friendly.

## Next Steps

Time-willing, I would have been interested in implementing the event QR scanner and the Hardware signout project.
I believe that both of these features would be interesting challenges to implement on the backend, but also pose a cool full-stack design challenge.

Besides designing a scalable backend, since I already had the frontend set up for this project, I would have loved to experiment with ideas for the user flow for volunteers and hackers alike.

A separate idea that I would love to implement is a Snack Shop system!

At TreeHacks 2024, they ran a Snack corner where you could go up and request any assortment of snacks from their collection.
I think it would be really fun to build out a system where hackers regularly get snack vouchers to redeem, but can also participate in workshops/events to earn special vouchers.
