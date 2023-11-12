CREATE TABLE IF NOT EXISTS weather (
    uid INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    city VARCHAR(15) NOT NULL,
    user_name VARCHAR(15) NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    weather_forecast JSON NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_valid_until ON weather(valid_until);
