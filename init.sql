CREATE DATABASE IF NOT EXISTS pet_daisy;
USE pet_daisy;

CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(36) UNIQUE NOT NULL,
    display_name varchar(32) UNIQUE NOT NULL,
    sync_code varchar(6) DEFAULT NULL,
    pet_count int DEFAULT 0
);

