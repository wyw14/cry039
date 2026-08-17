INSERT INTO office_areas(id,floor,environment,capacity) VALUES ('north-2','2F','open',80),('north-3','3F','open',65) ON CONFLICT DO NOTHING;
