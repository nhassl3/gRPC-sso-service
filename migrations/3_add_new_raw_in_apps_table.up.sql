INSERT INTO apps (id, name, secret)
VALUES (1, "test", "test-secret-word_1")
ON CONFLICT DO NOTHING;
