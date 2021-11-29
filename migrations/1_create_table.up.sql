CREATE TABLE error_logs
(
    id            INTEGER PRIMARY KEY,
    time          TIME,
    stack_trace   VARCHAR(3000),
    user_agent    VARCHAR(255),
    http_code     INTEGER,
    app_name      VARCHAR(255),
    function_name VARCHAR(255)
);