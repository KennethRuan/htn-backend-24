-- +goose Up
-- +goose StatementBegin
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

CREATE INDEX idx_users_skills_user_id ON users_skills(user_id);
CREATE INDEX idx_users_skills_skill_id ON users_skills(skill_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX idx_users_skills_user_id;
DROP INDEX idx_users_skills_skill_id;

DROP TABLE users_skills;
DROP TABLE users;
DROP TABLE skills;
-- +goose StatementEnd
