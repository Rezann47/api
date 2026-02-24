CREATE TABLE study_plans (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_by   UUID NOT NULL REFERENCES users(id),
    title        VARCHAR(200) NOT NULL DEFAULT 'Çalışma Planı',
    plan_date    DATE NOT NULL,
    note         TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE study_plan_items (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id            UUID NOT NULL REFERENCES study_plans(id) ON DELETE CASCADE,
    subject_id         UUID NOT NULL REFERENCES subjects(id),
    topic_id           UUID REFERENCES topics(id),
    duration_minutes   INT NOT NULL DEFAULT 30,
    display_order      SMALLINT NOT NULL DEFAULT 0,
    is_completed       BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at       TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_study_plans_user_date ON study_plans (user_id, plan_date);
CREATE INDEX idx_study_plan_items_plan ON study_plan_items (plan_id);

CREATE OR REPLACE FUNCTION update_study_plan_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_study_plans_updated_at
    BEFORE UPDATE ON study_plans
    FOR EACH ROW EXECUTE FUNCTION update_study_plan_updated_at();

CREATE TRIGGER trg_study_plan_items_updated_at
    BEFORE UPDATE ON study_plan_items
    FOR EACH ROW EXECUTE FUNCTION update_study_plan_updated_at();
