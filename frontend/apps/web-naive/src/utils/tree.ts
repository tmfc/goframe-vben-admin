
interface TreeConfig {
  id?: string;
  pid?: string;
  children?: string;
}

const DEFAULT_CONFIG: Required<TreeConfig> = {
  id: 'id',
  pid: 'parentId',
  children: 'children',
};

function getConfig(config?: TreeConfig): Required<TreeConfig> {
  return Object.assign({}, DEFAULT_CONFIG, config);
}

export function listToTree<T = any>(list: T[], config?: TreeConfig): T[] {
  const conf = getConfig(config);
  const nodeMap = new Map();
  const roots: T[] = [];
  const { id, children, pid } = conf;

  for (const item of list) {
    // @ts-ignore
    nodeMap.set(item[id], { ...item, [children]: [] });
  }

  for (const item of list) {
    // @ts-ignore
    const parentId = item[pid];
    // @ts-ignore
    const node = nodeMap.get(item[id]);
    if (!node) continue;
    
    // Treat 0, '0', null, undefined as root if parent doesn't exist
    if (
      parentId === null || 
      parentId === undefined || 
      parentId === '0' || 
      parentId === 0 || 
      !nodeMap.has(parentId)
    ) {
      roots.push(node);
    } else {
      // @ts-ignore
      nodeMap.get(parentId)[children].push(node);
    }
  }
  return roots;
}

export function flattenTree<T = any>(tree: T[], config?: TreeConfig): T[] {
  const { children } = getConfig(config);
  const result: T[] = [];
  function walk(nodes: T[]) {
    for (const node of nodes) {
      result.push(node);
      // @ts-ignore
      if (node[children] && node[children].length > 0) {
        // @ts-ignore
        walk(node[children]);
      }
    }
  }
  walk(tree);
  return result;
}

export function collectExpandedKeys<T = any>(tree: T[], maxDepth: number, config?: TreeConfig): string[] {
  const { id, children } = getConfig(config);
  const keys: string[] = [];
  function walk(nodes: T[], depth: number) {
    for (const node of nodes) {
      if (depth <= maxDepth) {
        // @ts-ignore
        keys.push(String(node[id]));
      }
      // @ts-ignore
      if (node[children] && node[children].length > 0) {
        // @ts-ignore
        walk(node[children], depth + 1);
      }
    }
  }
  walk(tree, 1);
  return keys;
}

export function sortTree<T = any>(tree: T[], compareFn: (a: T, b: T) => number, config?: TreeConfig): T[] {
  const { children } = getConfig(config);
  tree.sort(compareFn);
  for (const node of tree) {
    // @ts-ignore
    if (node[children] && node[children].length > 0) {
      // @ts-ignore
      sortTree(node[children], compareFn, config);
    }
  }
  return tree;
}
