import { describe, it, expect } from 'vitest';
import { generatePermissionCode } from '../utils';

describe('Menu Utils', () => {
  const mockMenus = [
    { id: '1', name: 'System', permissionCode: 'Menu:System' },
    { id: '2', name: 'User', parentId: '1', permissionCode: 'Menu:System:User' },
  ];

  it('should generate code for catalog', () => {
    const code = generatePermissionCode('Tools', 'catalog');
    expect(code).toBe('Catalog:Tools');
  });

  it('should generate code for button with parent', () => {
    const code = generatePermissionCode('Add', 'button', '2', mockMenus as any);
    expect(code).toBe('Menu:System:User:Add');
  });

  it('should generate deep path for menu', () => {
    const code = generatePermissionCode('Logs', 'menu', '2', mockMenus as any);
    expect(code).toBe('Menu:System:User:Logs');
  });

  it('should handle special characters in name', () => {
    const code = generatePermissionCode('User management', 'menu');
    expect(code).toBe('Menu:UserManagement');
  });
});
