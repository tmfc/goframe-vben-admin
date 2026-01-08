import { describe, it, expect } from 'vitest';
import { useFormSchema } from '../data';

describe('User Data', () => {
  it('should have a schemas array with at least one schema', () => {
    const schemas = useFormSchema();
    expect(Array.isArray(schemas)).toBe(true);
    expect(schemas.length).toBeGreaterThan(0);
  });
});
