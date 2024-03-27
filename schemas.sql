CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(63) UNIQUE NOT NULL,
  username VARCHAR(63) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE urls(
  id SERIAL PRIMARY KEY,
  url_id VARCHAR(63) UNIQUE NOT NULL,
  url VARCHAR(2000) NOT NULL,
  short_url VARCHAR(63) UNIQUE NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
);

CREATE TABLE user_urls(
  id SERIAL PRIMARY KEY,
  url_id VARCHAR(63) NOT NULL,
  user_id VARCHAR(63) NOT NULL,

  UNIQUE (url_id, user_id),
  FOREIGN KEY (url_id) REFERENCES urls(url_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE statistics(
  id BIGSERIAL PRIMARY KEY,
  url_id VARCHAR(63) NOT NULL,
  ip_address VARCHAR(45),
  user_agent VARCHAR(255),
  referer_url VARCHAR(2000),
  latitude DECIMAL(10, 8),
  Longitude DECIMAL(11, 8),
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

  FOREIGN KEY(url_id) REFERENCES urls(url_id) ON DELETE CASCADE
);

CREATE TABLE open_graphs(
  id SERIAL PRIMARY KEY,
  url_id VARCHAR(63) UNIQUE NOT NULL,
  title VARCHAR(255) NOT NULL,
  description VARCHAR(511) NOT NULL,
  image VARCHAR(255) NOT NULL,

  FOREIGN KEY(url_id) REFERENCES urls(url_id) ON DELETE CASCADE
);

CREATE TABLE admins(
  id SERIAL PRIMARY KEY,
  admin_id VARCHAR(63) UNIQUE NOT NULL,
  username VARCHAR(63) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL
);

CREATE TABLE error_logs(
  id SERIAL PRIMARY KEY,
  error_id VARCHAR(63) UNIQUE NOT NULL,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
