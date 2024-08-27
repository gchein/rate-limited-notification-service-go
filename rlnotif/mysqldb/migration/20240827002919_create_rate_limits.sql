DROP TABLE IF EXISTS rate_limits;
CREATE TABLE rate_limits (
  id                 BIGINT PRIMARY KEY AUTO_INCREMENT,
  notification_type  VARCHAR(255) NOT NULL,
  time_window        VARCHAR(255) NOT NULL,
  max_limit          INT NOT NULL,
  created_at         DATETIME NOT NULL,
  updated_at         DATETIME NOT NULL,
  UNIQUE(notification_type, time_window)
);
