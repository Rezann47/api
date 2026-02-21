-- ============================================================
-- 000006: Seed Data — TYT & AYT Subjects + Topics
-- ============================================================
-- Bu migration sadece INSERT yapar, her çalışmada tekrarlanabilir
-- olması için ON CONFLICT DO NOTHING kullanıyoruz.

-- ─── TYT SUBJECTS ─────────────────────────────────────────
INSERT INTO subjects (id, name, exam_type, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000001', 'Türkçe',    'TYT', 1),
    ('a0000001-0000-0000-0000-000000000002', 'Matematik', 'TYT', 2),
    ('a0000001-0000-0000-0000-000000000003', 'Geometri',  'TYT', 3),
    ('a0000001-0000-0000-0000-000000000004', 'Fizik',     'TYT', 4),
    ('a0000001-0000-0000-0000-000000000005', 'Kimya',     'TYT', 5),
    ('a0000001-0000-0000-0000-000000000006', 'Biyoloji',  'TYT', 6),
    ('a0000001-0000-0000-0000-000000000007', 'Tarih',     'TYT', 7),
    ('a0000001-0000-0000-0000-000000000008', 'Coğrafya',  'TYT', 8),
    ('a0000001-0000-0000-0000-000000000009', 'Felsefe',   'TYT', 9),
    ('a0000001-0000-0000-0000-000000000010', 'Din',       'TYT', 10)
ON CONFLICT (name, exam_type) DO NOTHING;

-- ─── AYT SUBJECTS ─────────────────────────────────────────
INSERT INTO subjects (id, name, exam_type, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000001', 'Matematik',  'AYT', 1),
    ('b0000002-0000-0000-0000-000000000002', 'Fizik',      'AYT', 2),
    ('b0000002-0000-0000-0000-000000000003', 'Kimya',      'AYT', 3),
    ('b0000002-0000-0000-0000-000000000004', 'Biyoloji',   'AYT', 4),
    ('b0000002-0000-0000-0000-000000000005', 'Edebiyat',   'AYT', 5),
    ('b0000002-0000-0000-0000-000000000006', 'Tarih',      'AYT', 6),
    ('b0000002-0000-0000-0000-000000000007', 'Coğrafya',   'AYT', 7)
ON CONFLICT (name, exam_type) DO NOTHING;

-- ─── TYT TOPICS ───────────────────────────────────────────

-- Türkçe (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000001', 'Sözcükte Anlam',            1),
    ('a0000001-0000-0000-0000-000000000001', 'Cümlede Anlam',             2),
    ('a0000001-0000-0000-0000-000000000001', 'Paragraf',                  3),
    ('a0000001-0000-0000-0000-000000000001', 'Ses Bilgisi',               4),
    ('a0000001-0000-0000-0000-000000000001', 'Yazım Kuralları',           5),
    ('a0000001-0000-0000-0000-000000000001', 'Noktalama İşaretleri',      6),
    ('a0000001-0000-0000-0000-000000000001', 'Sözcük Türleri',            7),
    ('a0000001-0000-0000-0000-000000000001', 'Cümlenin Ögeleri',          8),
    ('a0000001-0000-0000-0000-000000000001', 'Cümle Türleri',             9),
    ('a0000001-0000-0000-0000-000000000001', 'Fiil Çekimi',               10),
    ('a0000001-0000-0000-0000-000000000001', 'Anlatım Bozukluğu',         11)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Matematik (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000002', 'Sayılar ve Sayı Sistemleri', 1),
    ('a0000001-0000-0000-0000-000000000002', 'Bölme ve Bölünebilme',       2),
    ('a0000001-0000-0000-0000-000000000002', 'EBOB - EKOK',                3),
    ('a0000001-0000-0000-0000-000000000002', 'Rasyonel Sayılar',           4),
    ('a0000001-0000-0000-0000-000000000002', 'Ondalık Sayılar',            5),
    ('a0000001-0000-0000-0000-000000000002', 'Basamak Kavramı',            6),
    ('a0000001-0000-0000-0000-000000000002', 'Üslü Sayılar',               7),
    ('a0000001-0000-0000-0000-000000000002', 'Köklü Sayılar',              8),
    ('a0000001-0000-0000-0000-000000000002', 'Çarpanlara Ayırma',          9),
    ('a0000001-0000-0000-0000-000000000002', 'Denklemler',                 10),
    ('a0000001-0000-0000-0000-000000000002', 'Eşitsizlikler',              11),
    ('a0000001-0000-0000-0000-000000000002', 'Mutlak Değer',               12),
    ('a0000001-0000-0000-0000-000000000002', 'Oran Orantı',                13),
    ('a0000001-0000-0000-0000-000000000002', 'Problemler',                 14),
    ('a0000001-0000-0000-0000-000000000002', 'Kümeler',                    15),
    ('a0000001-0000-0000-0000-000000000002', 'Mantık',                     16),
    ('a0000001-0000-0000-0000-000000000002', 'Fonksiyonlar',               17),
    ('a0000001-0000-0000-0000-000000000002', 'Polinomlar',                 18),
    ('a0000001-0000-0000-0000-000000000002', 'Olasılık',                   19),
    ('a0000001-0000-0000-0000-000000000002', 'İstatistik',                 20)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Geometri (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000003', 'Temel Kavramlar',        1),
    ('a0000001-0000-0000-0000-000000000003', 'Üçgenler',               2),
    ('a0000001-0000-0000-0000-000000000003', 'Özel Üçgenler',          3),
    ('a0000001-0000-0000-0000-000000000003', 'Çokgenler',              4),
    ('a0000001-0000-0000-0000-000000000003', 'Dörtgenler',             5),
    ('a0000001-0000-0000-0000-000000000003', 'Çember ve Daire',        6),
    ('a0000001-0000-0000-0000-000000000003', 'Katı Cisimler',          7),
    ('a0000001-0000-0000-0000-000000000003', 'Analitik Geometri',      8)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Fizik (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000004', 'Fizik Bilimine Giriş',   1),
    ('a0000001-0000-0000-0000-000000000004', 'Madde ve Özellikleri',   2),
    ('a0000001-0000-0000-0000-000000000004', 'Hareket',                3),
    ('a0000001-0000-0000-0000-000000000004', 'Kuvvet',                 4),
    ('a0000001-0000-0000-0000-000000000004', 'Enerji',                 5),
    ('a0000001-0000-0000-0000-000000000004', 'Isı ve Sıcaklık',        6),
    ('a0000001-0000-0000-0000-000000000004', 'Elektrostatik',          7),
    ('a0000001-0000-0000-0000-000000000004', 'Dalgalar',               8)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Kimya (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000005', 'Kimya Bilimine Giriş',   1),
    ('a0000001-0000-0000-0000-000000000005', 'Atom ve Yapısı',         2),
    ('a0000001-0000-0000-0000-000000000005', 'Periyodik Tablo',        3),
    ('a0000001-0000-0000-0000-000000000005', 'Kimyasal Bağlar',        4),
    ('a0000001-0000-0000-0000-000000000005', 'Maddenin Halleri',       5),
    ('a0000001-0000-0000-0000-000000000005', 'Karışımlar',             6),
    ('a0000001-0000-0000-0000-000000000005', 'Asit Baz',               7),
    ('a0000001-0000-0000-0000-000000000005', 'Kimyasal Tepkimeler',    8),
    ('a0000001-0000-0000-0000-000000000005', 'Mol Kavramı',            9)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Biyoloji (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000006', 'Hücre',                  1),
    ('a0000001-0000-0000-0000-000000000006', 'Canlıların Sınıflandırılması', 2),
    ('a0000001-0000-0000-0000-000000000006', 'Kalıtım',                3),
    ('a0000001-0000-0000-0000-000000000006', 'Ekosistem',              4),
    ('a0000001-0000-0000-0000-000000000006', 'Sinir Sistemi',          5),
    ('a0000001-0000-0000-0000-000000000006', 'Endokrin Sistem',        6),
    ('a0000001-0000-0000-0000-000000000006', 'Üreme ve Gelişme',       7)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Tarih (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000007', 'Tarih Bilimi',                       1),
    ('a0000001-0000-0000-0000-000000000007', 'İlk Türk Devletleri',               2),
    ('a0000001-0000-0000-0000-000000000007', 'İslam Tarihi',                       3),
    ('a0000001-0000-0000-0000-000000000007', 'Osmanlı Kuruluş Dönemi',            4),
    ('a0000001-0000-0000-0000-000000000007', 'Osmanlı Yükselme Dönemi',           5),
    ('a0000001-0000-0000-0000-000000000007', 'Osmanlı Gerileme ve Dağılma',       6),
    ('a0000001-0000-0000-0000-000000000007', 'Kurtuluş Savaşı',                   7),
    ('a0000001-0000-0000-0000-000000000007', 'Atatürk İlke ve İnkılapları',       8),
    ('a0000001-0000-0000-0000-000000000007', 'Türkiye Cumhuriyeti Tarihi',        9)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Coğrafya (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000008', 'Harita Bilgisi',           1),
    ('a0000001-0000-0000-0000-000000000008', 'Türkiyenin Coğrafi Konumu', 2),
    ('a0000001-0000-0000-0000-000000000008', 'İklim',                    3),
    ('a0000001-0000-0000-0000-000000000008', 'Nüfus ve Yerleşme',        4),
    ('a0000001-0000-0000-0000-000000000008', 'Ekonomik Faaliyetler',     5),
    ('a0000001-0000-0000-0000-000000000008', 'Doğal Afetler',            6),
    ('a0000001-0000-0000-0000-000000000008', 'Çevre Sorunları',          7)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Felsefe (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000009', 'Felsefenin Konusu',        1),
    ('a0000001-0000-0000-0000-000000000009', 'Bilgi Felsefesi',          2),
    ('a0000001-0000-0000-0000-000000000009', 'Varlık Felsefesi',         3),
    ('a0000001-0000-0000-0000-000000000009', 'Ahlak Felsefesi',          4),
    ('a0000001-0000-0000-0000-000000000009', 'Siyaset Felsefesi',        5),
    ('a0000001-0000-0000-0000-000000000009', 'Sanat Felsefesi',          6),
    ('a0000001-0000-0000-0000-000000000009', 'Din Felsefesi',            7),
    ('a0000001-0000-0000-0000-000000000009', 'Mantık',                   8),
    ('a0000001-0000-0000-0000-000000000009', 'Psikoloji',                9),
    ('a0000001-0000-0000-0000-000000000009', 'Sosyoloji',                10)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Din (TYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('a0000001-0000-0000-0000-000000000010', 'İslam Dini ve İnanç',     1),
    ('a0000001-0000-0000-0000-000000000010', 'Kuran ve Yorumu',         2),
    ('a0000001-0000-0000-0000-000000000010', 'Hz. Muhammed',            3),
    ('a0000001-0000-0000-0000-000000000010', 'İbadetler',               4),
    ('a0000001-0000-0000-0000-000000000010', 'İslam Ahlakı',            5),
    ('a0000001-0000-0000-0000-000000000010', 'Din ve Laiklik',          6)
ON CONFLICT (name, subject_id) DO NOTHING;

-- ─── AYT TOPICS ───────────────────────────────────────────

-- Matematik (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000001', 'Türev',                    1),
    ('b0000002-0000-0000-0000-000000000001', 'İntegral',                 2),
    ('b0000002-0000-0000-0000-000000000001', 'Logaritma',                3),
    ('b0000002-0000-0000-0000-000000000001', 'Trigonometri',             4),
    ('b0000002-0000-0000-0000-000000000001', 'Karmaşık Sayılar',         5),
    ('b0000002-0000-0000-0000-000000000001', 'Diziler',                  6),
    ('b0000002-0000-0000-0000-000000000001', 'Seriler',                  7),
    ('b0000002-0000-0000-0000-000000000001', 'Kombinatorik',             8),
    ('b0000002-0000-0000-0000-000000000001', 'Olasılık (İleri)',         9),
    ('b0000002-0000-0000-0000-000000000001', 'Analitik Geometri (İleri)',10),
    ('b0000002-0000-0000-0000-000000000001', 'Vektörler',                11),
    ('b0000002-0000-0000-0000-000000000001', 'Matrisler',                12)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Fizik (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000002', 'Kuvvet ve Hareket',        1),
    ('b0000002-0000-0000-0000-000000000002', 'Newton Kanunları',         2),
    ('b0000002-0000-0000-0000-000000000002', 'Enerji ve Güç',            3),
    ('b0000002-0000-0000-0000-000000000002', 'İtme ve Momentum',         4),
    ('b0000002-0000-0000-0000-000000000002', 'Dairesel Hareket',         5),
    ('b0000002-0000-0000-0000-000000000002', 'Gravitasyon',              6),
    ('b0000002-0000-0000-0000-000000000002', 'Basit Harmonik Hareket',   7),
    ('b0000002-0000-0000-0000-000000000002', 'Elektrik ve Manyetizma',   8),
    ('b0000002-0000-0000-0000-000000000002', 'Modern Fizik',             9),
    ('b0000002-0000-0000-0000-000000000002', 'Optik',                    10)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Kimya (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000003', 'Gazlar',                   1),
    ('b0000002-0000-0000-0000-000000000003', 'Sıvı Çözeltiler',          2),
    ('b0000002-0000-0000-0000-000000000003', 'Kimyasal Denge',           3),
    ('b0000002-0000-0000-0000-000000000003', 'Elektrokimya',             4),
    ('b0000002-0000-0000-0000-000000000003', 'Organik Kimya',            5),
    ('b0000002-0000-0000-0000-000000000003', 'Kimya ve Enerji',          6),
    ('b0000002-0000-0000-0000-000000000003', 'Reaksiyon Hızı',           7)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Biyoloji (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000004', 'Hücre Bölünmeleri',        1),
    ('b0000002-0000-0000-0000-000000000004', 'Kalıtım ve Genetik',       2),
    ('b0000002-0000-0000-0000-000000000004', 'DNA ve Genetik Kod',       3),
    ('b0000002-0000-0000-0000-000000000004', 'Evrim',                    4),
    ('b0000002-0000-0000-0000-000000000004', 'Sistemler Biyolojisi',     5),
    ('b0000002-0000-0000-0000-000000000004', 'Ekoloji',                  6),
    ('b0000002-0000-0000-0000-000000000004', 'Biyoteknoloji',            7)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Edebiyat (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000005', 'Türk Edebiyatına Giriş',   1),
    ('b0000002-0000-0000-0000-000000000005', 'Divan Edebiyatı',          2),
    ('b0000002-0000-0000-0000-000000000005', 'Tanzimat Edebiyatı',       3),
    ('b0000002-0000-0000-0000-000000000005', 'Servet-i Fünun',           4),
    ('b0000002-0000-0000-0000-000000000005', 'Milli Edebiyat',           5),
    ('b0000002-0000-0000-0000-000000000005', 'Cumhuriyet Edebiyatı',     6),
    ('b0000002-0000-0000-0000-000000000005', 'Roman ve Hikaye',          7),
    ('b0000002-0000-0000-0000-000000000005', 'Şiir Türleri',             8)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Tarih (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000006', 'Osmanlı Tarihi (Detay)',          1),
    ('b0000002-0000-0000-0000-000000000006', 'I. Dünya Savaşı',                 2),
    ('b0000002-0000-0000-0000-000000000006', 'Kurtuluş Savaşı (Detay)',         3),
    ('b0000002-0000-0000-0000-000000000006', 'Cumhuriyet Dönemi',               4),
    ('b0000002-0000-0000-0000-000000000006', 'Atatürk İlkeleri (Detay)',        5),
    ('b0000002-0000-0000-0000-000000000006', 'II. Dünya Savaşı ve Türkiye',     6),
    ('b0000002-0000-0000-0000-000000000006', 'Soğuk Savaş Dönemi',              7),
    ('b0000002-0000-0000-0000-000000000006', 'Çağdaş Türk ve Dünya Tarihi',    8)
ON CONFLICT (name, subject_id) DO NOTHING;

-- Coğrafya (AYT)
INSERT INTO topics (subject_id, name, display_order) VALUES
    ('b0000002-0000-0000-0000-000000000007', 'Doğal Sistemler',           1),
    ('b0000002-0000-0000-0000-000000000007', 'Beşeri Sistemler',          2),
    ('b0000002-0000-0000-0000-000000000007', 'Küresel Ortam',             3),
    ('b0000002-0000-0000-0000-000000000007', 'Türkiye Ekonomisi',         4),
    ('b0000002-0000-0000-0000-000000000007', 'Bölgeler',                  5),
    ('b0000002-0000-0000-0000-000000000007', 'Küresel Çevre Sorunları',   6)
ON CONFLICT (name, subject_id) DO NOTHING;
