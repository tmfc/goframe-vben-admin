-- Remove seeded Casbin policies for default tenant.
DELETE FROM casbin_rule
WHERE ptype = 'p'
  AND v1 = '00000000-0000-0000-0000-000000000000'
  AND v2 = '*'
  AND v3 = '*'
  AND v0 IN ('super', 'admin', 'staff');
