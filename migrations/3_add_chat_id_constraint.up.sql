ALTER TABLE chat_ids
    ADD CONSTRAINT app_name CHECK (LENGTH(app_name) > 0);