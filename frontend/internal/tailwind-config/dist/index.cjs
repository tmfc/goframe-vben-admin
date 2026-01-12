const { createJiti } = require("../../../node_modules/.pnpm/jiti@2.6.1/node_modules/jiti/lib/jiti.cjs")

const jiti = createJiti(__filename, {
  "interopDefault": true,
  "alias": {
    "@vben/tailwind-config": "/Users/william/.conduit/worktrees/goframe-vben-admin/20260112-113357-9ee810e1/frontend/internal/tailwind-config"
  },
  "transformOptions": {
    "babel": {
      "plugins": []
    }
  }
})

/** @type {import("/Users/william/.conduit/worktrees/goframe-vben-admin/20260112-113357-9ee810e1/frontend/internal/tailwind-config/src/index.js")} */
module.exports = jiti("/Users/william/.conduit/worktrees/goframe-vben-admin/20260112-113357-9ee810e1/frontend/internal/tailwind-config/src/index.ts")