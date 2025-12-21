INSERT INTO users (id, name, email) VALUES
('11111111-1111-1111-1111-111111111111', 'User One', 'u1@test.com'),
('22222222-2222-2222-2222-222222222222', 'User Two', 'u2@test.com');

INSERT INTO groups (id, name) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Test Group');

INSERT INTO group_members (group_id, user_id) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '22222222-2222-2222-2222-222222222222');
