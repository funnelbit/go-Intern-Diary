CREATE TABLE user (
  `id` BIGINT UNSIGNED NOT NULL,
  `name` VARBINARY(32) NOT NULL,
  `password_hash` VARBINARY(254) NOT NULL,

  `created_at` DATETIME(6) NOT NULL,
  `updated_at` DATETIME(6) NOT NULL,

  PRIMARY KEY (id),
  UNIQUE KEY (name),

  KEY (created_at),
  KEY (updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE user_session (
  `id` BIGINT UNSIGNED NOT NULL,
  `user_id` BIGINT UNSIGNED NOT NULL,
  `expred_at` BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
