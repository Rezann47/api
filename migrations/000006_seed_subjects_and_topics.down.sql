-- Seed data rollback: önce topics (FK), sonra subjects
DELETE FROM topics WHERE subject_id IN (
    SELECT id FROM subjects WHERE exam_type IN ('TYT', 'AYT')
    AND id::text LIKE 'a0000001%' OR id::text LIKE 'b0000002%'
);
DELETE FROM subjects WHERE id::text LIKE 'a0000001%' OR id::text LIKE 'b0000002%';
