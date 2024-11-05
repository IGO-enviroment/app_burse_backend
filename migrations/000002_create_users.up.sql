CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,

  first_name VARCHAR(255),
  middle_name VARCHAR(255),
  last_name VARCHAR(255),

  email VARCHAR(255) UNIQUE NOT NULL,

  digest_password VARCHAR(255) NOT NULL,

  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
