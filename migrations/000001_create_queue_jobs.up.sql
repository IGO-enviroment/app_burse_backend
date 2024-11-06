CREATE TABLE IF NOT EXISTS queue_jobs (
  id SERIAL PRIMARY KEY,

  reserv_at TIMESTAMP,

  processed BOOLEAN DEFAULT false,
  queue_name VARCHAR(255) NOT NULL,
  run_at TIMESTAMP,

  method VARCHAR(255) NOT NULL,
  item TEXT NOT NULL,

  created_at TIMESTAMP DEFAULT NOW()	
);
