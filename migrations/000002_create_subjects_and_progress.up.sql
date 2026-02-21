-- ============================================================
-- 000002: Subjects, Topics & Student Progress
-- ============================================================

CREATE TYPE exam_type AS ENUM ('TYT', 'AYT');

-- ─── SUBJECTS ─────────────────────────────────────────────
-- Seed data ile doldurulacak (TYT/AYT dersleri)
CREATE TABLE subjects (
    id            UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(100) NOT NULL,
    exam_type     exam_type    NOT NULL,
    display_order SMALLINT     NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_subjects_name_exam UNIQUE (name, exam_type)
);

CREATE INDEX idx_subjects_exam_type ON subjects (exam_type, display_order);

CREATE TRIGGER trg_subjects_updated_at
    BEFORE UPDATE ON subjects
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─── TOPICS ───────────────────────────────────────────────
CREATE TABLE topics (
    id            UUID       PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_id    UUID         NOT NULL REFERENCES subjects (id) ON DELETE CASCADE,
    name          VARCHAR(200) NOT NULL,
    display_order SMALLINT     NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_topics_name_subject UNIQUE (name, subject_id)
);

CREATE INDEX idx_topics_subject_id ON topics (subject_id, display_order);

CREATE TRIGGER trg_topics_updated_at
    BEFORE UPDATE ON topics
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─── STUDENT TOPIC PROGRESS ───────────────────────────────
-- Sadece tamamlanan konular kaydedilir (sparse model).
-- is_completed=false satırı tutmak yerine satırı siliyoruz.
-- Bu yaklaşım büyük tablolarda çok daha verimli.
CREATE TABLE student_topic_progress (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID        NOT NULL REFERENCES users    (id) ON DELETE CASCADE,
    topic_id        UUID        NOT NULL REFERENCES topics   (id) ON DELETE CASCADE,
    completion_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_progress_user_topic UNIQUE (user_id, topic_id)
);

-- Öğrencinin tüm ilerlemesini çekmek için (en çok kullanılan sorgu)
CREATE INDEX idx_progress_user_id      ON student_topic_progress (user_id);

-- Belirli bir dersin tamamlanma yüzdesi hesaplamak için
CREATE INDEX idx_progress_user_topic   ON student_topic_progress (user_id, topic_id);

-- Eğitmenin öğrencisini görüntülemesi için de aynı index yeterli

CREATE TRIGGER trg_progress_updated_at
    BEFORE UPDATE ON student_topic_progress
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
