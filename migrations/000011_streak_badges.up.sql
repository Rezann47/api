-- users tablosuna streak kolonları ekle
ALTER TABLE users ADD COLUMN IF NOT EXISTS current_streak  INT          NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS longest_streak  INT          NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_study_date DATE;

-- Rozetler tablosu
CREATE TABLE IF NOT EXISTS badges (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_key  VARCHAR(50) NOT NULL,
    badge_name VARCHAR(100) NOT NULL,
    badge_icon VARCHAR(10)  NOT NULL,
    earned_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, badge_key)
);

CREATE INDEX IF NOT EXISTS idx_badges_user_id ON badges(user_id);

-- Haftalık liderlik view'ı
CREATE OR REPLACE VIEW leaderboard_weekly AS
SELECT
    u.id,
    u.name AS full_name,
    u.avatar_id,
    u.current_streak,
    COALESCE(SUM(p.duration_minutes), 0)::INT AS total_minutes,
    COALESCE(COUNT(p.id), 0)::INT             AS session_count
FROM users u
LEFT JOIN pomodoros p
    ON p.user_id = u.id
    AND p.started_at >= NOW() - INTERVAL '7 days'
WHERE u.role = 'student'
  AND u.deleted_at IS NULL
GROUP BY u.id, u.name AS full_name, u.avatar_id, u.current_streak
ORDER BY total_minutes DESC;
