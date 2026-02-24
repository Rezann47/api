CREATE TABLE messages (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sender_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    receiver_id  UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content      TEXT NOT NULL CHECK (char_length(content) BETWEEN 1 AND 2000),
    is_read      BOOLEAN NOT NULL DEFAULT false,
    read_at      TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT different_users CHECK (sender_id <> receiver_id)
);

-- Konuşma listesi & mesaj geçmişi için
CREATE INDEX idx_messages_sender_receiver ON messages (sender_id, receiver_id, created_at DESC);
CREATE INDEX idx_messages_receiver_sender ON messages (receiver_id, sender_id, created_at DESC);

-- Okunmamış sayısı için
CREATE INDEX idx_messages_unread ON messages (receiver_id, is_read) WHERE is_read = false;
