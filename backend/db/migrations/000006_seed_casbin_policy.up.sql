-- Seed Casbin policies to allow all authenticated users in default tenant.
-- Roles correspond to consts.RoleSuper, consts.RoleAdmin
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3)
VALUES
  ('p', 'super', '00000000-0000-0000-0000-000000000000', '*', '*'),
  ('p', 'admin', '00000000-0000-0000-0000-000000000000', '*', '*'),
  ('p', 'staff', '00000000-0000-0000-0000-000000000000', '*', '*');
