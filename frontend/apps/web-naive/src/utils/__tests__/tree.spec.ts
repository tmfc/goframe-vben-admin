import { describe, it, expect } from 'vitest';
import { listToTree, flattenTree, collectExpandedKeys } from '../tree';

describe('Tree Utils', () => {
  const list = [
    { id: '1', parentId: '0', name: 'Root' },
    { id: '2', parentId: '1', name: 'Child 1' },
    { id: '3', parentId: '1', name: 'Child 2' },
    { id: '4', parentId: '2', name: 'Grandchild' },
    { id: '5', parentId: '0', name: 'Root 2' },
  ];

  it('listToTree should convert list to tree', () => {
    const tree = listToTree(list, { id: 'id', pid: 'parentId', children: 'children' });
    expect(tree).toHaveLength(2);
    expect(tree[0].id).toBe('1');
    expect(tree[0].children).toHaveLength(2);
    expect(tree[0].children[0].id).toBe('2');
    expect(tree[0].children[0].children).toHaveLength(1);
    expect(tree[0].children[0].children[0].id).toBe('4');
  });

  it('listToTree should handle custom field names', () => {
    const customList = [
      { key: 1, parent: 0, title: 'Root' },
      { key: 2, parent: 1, title: 'Child' },
    ];
    const tree = listToTree(customList, { id: 'key', pid: 'parent', children: 'subs' });
    expect(tree).toHaveLength(1);
    expect(tree[0].subs).toHaveLength(1);
  });

  it('flattenTree should convert tree to list', () => {
    const tree = [
      {
        id: '1',
        children: [
          { id: '2', children: [{ id: '4' }] },
          { id: '3' },
        ],
      },
      { id: '5' },
    ];
    const flattened = flattenTree(tree, { children: 'children' });
    expect(flattened).toHaveLength(5);
    expect(flattened.find(i => i.id === '4')).toBeDefined();
  });

  it('collectExpandedKeys should collect keys up to depth', () => {
    const tree = [
      {
        id: '1',
        children: [
          { id: '2', children: [{ id: '4' }] }, // Depth 2, 3
          { id: '3' },
        ],
      },
      { id: '5' },
    ];
    // Collect keys for depth <= 2 (so id 1, 2, 3, 5 should be expanded if they have children? usually logic is to expand node if it has children and is within depth)
    // The previous implementation in index.vue:
    // if (depth <= 2) expandedKeys.push(node.id);
    
    const keys = collectExpandedKeys(tree, 2, { id: 'id', children: 'children' });
    expect(keys).toContain('1');
    expect(keys).toContain('2');
    // 3 and 5 are depth 1 but have no children? Or do we expand leaf nodes? 
    // The logic in index.vue was: if (depth <= 2) push(id).
    // So 1 (d1), 2 (d2), 3 (d2), 5 (d1).
    expect(keys).toContain('3');
    expect(keys).toContain('5');
    expect(keys).not.toContain('4'); // Depth 3
  });
});
