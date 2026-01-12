DELETE FROM casbin_rule WHERE v2 LIKE '/message/%';
DELETE FROM sys_permission WHERE name LIKE 'System:Message%';
