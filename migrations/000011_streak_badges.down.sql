DROP VIEW  IF EXISTS leaderboard_weekly;
DROP TABLE IF EXISTS badges;
ALTER TABLE users DROP COLUMN IF EXISTS current_streak;
ALTER TABLE users DROP COLUMN IF EXISTS longest_streak;
ALTER TABLE users DROP COLUMN IF EXISTS last_study_date;
