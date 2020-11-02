CREATE DATABASE IF NOT EXISTS chat;

USE chat;

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    identifier VARCHAR(50) UNIQUE,
    full_name VARCHAR(200)
);

CREATE TABLE rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    identifier VARCHAR(50) UNIQUE,
    valid BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    by_user INT,
    to_room INT,
    message VARCHAR(300),
    deleted BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (by_user) REFERENCES users(id),
    FOREIGN KEY (to_room) REFERENCES rooms(id)
);