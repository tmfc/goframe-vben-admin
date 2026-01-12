import { createJiti } from "../../../node_modules/.pnpm/jiti@2.6.1/node_modules/jiti/lib/jiti.mjs";

const jiti = createJiti(import.meta.url, {
  "interopDefault": true,
  "alias": {
    "@vben/turbo-run": "/Users/william/.conduit/worktrees/goframe-vben-admin/20260112-113357-9ee810e1/frontend/scripts/turbo-run"
  },
  "transformOptions": {
    "babel": {
      "plugins": []
    }
  }
})

/** @type {import("/Users/william/.conduit/worktrees/goframe-vben-admin/20260112-113357-9ee810e1/frontend/scripts/turbo-run/src/index.js")} */
const _module = await jiti.import("/Users/william/.conduit/worktrees/goframe-vben-admin/20260112-113357-9ee810e1/frontend/scripts/turbo-run/src/index.ts");

export default _module?.default ?? _module;