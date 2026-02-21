-- ============================================================
-- 000003: Pomodoro Sessions
-- ============================================================

CREATE TABLE pomodoros (
    id               UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES users    (id) ON DELETE CASCADE,
    subject_id       UUID                 REFERENCES subjects (id) ON DELETE SET NULL,
    duration_minutes SMALLINT    NOT NULL CHECK (duration_minutes > 0 AND duration_minutes <= 480),
    -- "date" değil TIMESTAMPTZ: zaman dilimi farkı olan kullanıcılarda
    -- sadece DATE tutmak hatalı günlere düşürebilir.
    -- Uygulama tarafında DATE hesaplaması user timezone'una göre yapılır.
    started_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tarih filtresi + toplam süre hesabı için kritik composite index
CREATE INDEX idx_pomodoros_user_started
    ON pomodoros (user_id, started_at DESC);

-- Ders bazlı filtreleme (opsiyonel feature için hazır)
CREATE INDEX idx_pomodoros_subject
    ON pomodoros (subject_id) WHERE subject_id IS NOT NULL;

CREATE TRIGGER trg_pomodoros_updated_at
    BEFORE UPDATE ON pomodoros
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
