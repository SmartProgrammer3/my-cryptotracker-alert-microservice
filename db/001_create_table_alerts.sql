CREATE TABLE alerts (
    id            BIGINT        UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    symbol        VARCHAR(20)   NOT NULL,
    target_price  DECIMAL(20,8) NOT NULL,
    direction     ENUM('ABOVE', 'BELOW') NOT NULL,
    status        ENUM('PENDING', 'TRIGGERED', 'CANCELLED') NOT NULL,
    created_at    DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    triggered_at  DATETIME      NULL,

    INDEX idx_symbol_status (symbol, status)
);