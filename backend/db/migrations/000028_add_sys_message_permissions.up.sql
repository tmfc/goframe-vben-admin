-- 1. Add Permissions
INSERT INTO sys_permission (name, description, status) VALUES ('System:Message', 'Message Center root', 1);
INSERT INTO sys_permission (name, description, parent_id, status) 
SELECT 'System:Message:List', 'Message list', id, 1 FROM sys_permission WHERE name = 'System:Message';
INSERT INTO sys_permission (name, description, parent_id, status) 
SELECT 'System:Message:SetRead', 'Mark message as read', id, 1 FROM sys_permission WHERE name = 'System:Message:List';
INSERT INTO sys_permission (name, description, parent_id, status) 
SELECT 'System:Message:SetAllRead', 'Mark all messages as read', id, 1 FROM sys_permission WHERE name = 'System:Message:List';

-- 2. Add Casbin Rules (Allow access to /message/* for roles)
-- Format: p, role, tenantId, path, method
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'super', '1', '/message/*', '*');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'admin', '1', '/message/*', '*');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'user', '1', '/message/*', '*');
INSERT INTO casbin_rule (ptype, v0, v1, v2, v3) VALUES ('p', 'guest', '1', '/message/*', 'GET');
