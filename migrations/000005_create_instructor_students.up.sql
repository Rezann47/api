-- ============================================================
-- 000005: Instructor-Student Relationship
-- ============================================================

CREATE TABLE instructor_students (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    instructor_id UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    student_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_instructor_student UNIQUE (instructor_id, student_id),
    -- Bir kullanıcı hem eğitmen hem öğrenci olmasın
    CONSTRAINT chk_not_self CHECK (instructor_id != student_id)
);

-- Eğitmenin öğrencilerini listesi
CREATE INDEX idx_instructor_students_instructor
    ON instructor_students (instructor_id);

-- Öğrencinin hangi eğitmenlere bağlı olduğu (ters arama)
CREATE INDEX idx_instructor_students_student
    ON instructor_students (student_id);
