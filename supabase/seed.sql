INSERT INTO auth.users (id, instance_id, aud, role, email, encrypted_password, email_confirmed_at)
VALUES (
  '5516e359-6c9c-4ebb-a409-52373d536d50', 
  gen_random_uuid(), -- Instance ID
  'authenticated',   -- Audience
  'authenticated',   -- Role
  'noah@email.com', 
  crypt('password', gen_salt('bf')), -- Hash the password
  NOW() -- Email confirmed
),
(
  '72c5e147-c5d8-4840-8787-6f8637e537b5', 
  gen_random_uuid(), -- Instance ID
  'authenticated',   -- Audience
  'authenticated',   -- Role
  'bob@yahoo.com', 
  crypt('password', gen_salt('bf')), -- Hash the password
  NOW() -- Email confirmed
),
(
  'f73b3d99-44f4-4fbc-9e23-17a310202b07', 
  gen_random_uuid(), -- Instance ID
  'authenticated',   -- Audience
  'authenticated',   -- Role
  'eric@gmail.com', 
  crypt('password', gen_salt('bf')), -- Hash the password
  NOW() -- Email confirmed
);

INSERT INTO users (id, first_name, last_name) 
VALUES (
    '5516e359-6c9c-4ebb-a409-52373d536d50', 
    'Noah',
    'Libeskind'
),(
    '72c5e147-c5d8-4840-8787-6f8637e537b5', 
    'Bob',
    'Smith'
),(
    'f73b3d99-44f4-4fbc-9e23-17a310202b07', 
    'Eric',
    'Jamison'
);

INSERT INTO workspaces (id, name, owner_id) 
VALUES (
    '6cf86c8e-5e38-4af1-8755-295ad12ed91b',
    'Friends of Noah',
    '5516e359-6c9c-4ebb-a409-52373d536d50'
),
(
    '9e3018cd-89dd-4af8-8fee-9aaaba3549b7',
    'Bird Watching',
    '5516e359-6c9c-4ebb-a409-52373d536d50'
);

INSERT INTO workspace_members (member_id, workspace_id) 
VALUES (
    'f73b3d99-44f4-4fbc-9e23-17a310202b07', -- Eric
    '6cf86c8e-5e38-4af1-8755-295ad12ed91b' -- Friends of Noah
),
(
    '5516e359-6c9c-4ebb-a409-52373d536d50', -- Noah
    '6cf86c8e-5e38-4af1-8755-295ad12ed91b' -- Friends of Noah
),
(
    '72c5e147-c5d8-4840-8787-6f8637e537b5', -- Bob
    '9e3018cd-89dd-4af8-8fee-9aaaba3549b7' -- Bird Watching
),
(
    '5516e359-6c9c-4ebb-a409-52373d536d50', -- Noah
    '9e3018cd-89dd-4af8-8fee-9aaaba3549b7' -- Bird Watching
);
