-- Users
INSERT INTO users (name, email, password) 
VALUES
('user 1', 'user1@gmail.com', '$2a$05$wQ8lYAdEw7ZzF3OSzWeCKee8wc0KWxbBqfJpNu.lb.f1rvuSyy/I2');

-- POS
INSERT INTO pos (user_id, name, type, color)
VALUES
(1, 'pos 1', 0, '#FF00FF');