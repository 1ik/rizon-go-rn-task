-- +goose Up
CREATE TABLE feedbacks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedbacks_user_id ON feedbacks(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_feedbacks_user_id;
DROP TABLE IF EXISTS feedbacks;
