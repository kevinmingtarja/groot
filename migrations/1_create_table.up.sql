CREATE TABLE error_logs
(
    id            SERIAL PRIMARY KEY,
    time          TIME,
    request_url   VARCHAR(255),
    stack_trace   VARCHAR(3000),
    user_agent    VARCHAR(255),
    http_code     INTEGER,
    app_name      VARCHAR(255),
    function_name VARCHAR(255)
);
