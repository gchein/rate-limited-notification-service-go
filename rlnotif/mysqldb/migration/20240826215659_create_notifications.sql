CREATE TABLE notifications (
  id                 BIGINT PRIMARY KEY AUTO_INCREMENT,
  notification_type  VARCHAR(255) NOT NULL,
  message            VARCHAR(255) NOT NULL,
  user_id            BIGINT NOT NULL,
  created_at         DATETIME NOT NULL,
  updated_at         DATETIME NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);
