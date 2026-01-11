CREATE TABLE
   IF NOT EXISTS users (
      id UUID PRIMARY KEY,
      username VARCHAR(255) NOT NULL,
      email VARCHAR(255) NOT NULL UNIQUE,
      password VARCHAR(255) NOT NULL,
      profile_picture TEXT,
      created_at TIMESTAMP
      WITH
         TIME ZONE DEFAULT NOW (),
         updated_at TIMESTAMP
      WITH
         TIME ZONE DEFAULT NOW ()
   );

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);