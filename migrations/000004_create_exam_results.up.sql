-- ============================================================
-- 000004: Exam Results
-- ============================================================
-- Tasarım kararı: TYT'de 10, AYT'de 7 ders var.
-- Her ders için 3 kolon (correct/wrong/net) = 30-51 kolon.
-- Bu yaklaşım yerine JSONB kullanıyoruz:
--
--   AVANTAJLAR:
--   • Yeni ders eklendiğinde migration gerekmez
--   • Schema daha temiz
--   • Partial update kolaylaştı
--   • Farklı sınav tipleri (LGS, KPSS) eklenebilir
--
--   DEZAVANTAJLAR:
--   • Sorgu biraz daha karmaşık (JSONB operatörleri)
--   • GIN index gerekiyor
--
--   Alternatif: Eğer sadece aggregate sorgular yapılacaksa
--   (toplam net, ortalama) flat kolon daha hızlı.
--   Burada mixed kullanım var, JSONB daha esnektir.
-- ============================================================

CREATE TABLE exam_results (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    exam_type  exam_type   NOT NULL,
    exam_date  DATE        NOT NULL,
    -- Her dersin correct/wrong/net değerleri:
    -- {"turkish": {"correct": 28, "wrong": 4, "net": 27.00},
    --  "math":    {"correct": 15, "wrong": 8, "net": 12.33}, ...}
    scores     JSONB       NOT NULL DEFAULT '{}',
    total_net  NUMERIC(6,2) NOT NULL DEFAULT 0,   -- uygulama katmanında hesaplanır
    note       TEXT,                               -- öğrencinin notu (opsiyonel)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Aynı güne aynı tipte birden fazla deneme girilmesin diye kısıtlama YOK
    -- (öğrenci günde birden fazla deneme çözebilir)
    CONSTRAINT chk_total_net CHECK (total_net >= -100 AND total_net <= 200)
);

CREATE INDEX idx_exam_results_user_date
    ON exam_results (user_id, exam_date DESC);

CREATE INDEX idx_exam_results_user_type
    ON exam_results (user_id, exam_type, exam_date DESC);

-- JSONB içinde belirli ders netlerini aramak için (opsiyonel)
CREATE INDEX idx_exam_results_scores_gin
    ON exam_results USING GIN (scores);

CREATE TRIGGER trg_exam_results_updated_at
    BEFORE UPDATE ON exam_results
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ─── SCORES VALIDATION ────────────────────────────────────
-- JSONB'nin geçerli formatta girildiğini doğrulayan check constraint
-- Sadece non-empty olduğunu kontrol ediyoruz, detaylı validasyon app katmanında
ALTER TABLE exam_results
    ADD CONSTRAINT chk_scores_not_empty
    CHECK (scores != '{}' AND jsonb_typeof(scores) = 'object');
