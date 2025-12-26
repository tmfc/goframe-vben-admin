const { createJiti } = require("../../../node_modules/.pnpm/jiti@2.6.1/node_modules/jiti/lib/jiti.cjs")

const jiti = createJiti(__filename, {
  "interopDefault": true,
  "alias": {
    "@vben/tailwind-config": "/Users/william/work/my_opensource/goframe-vben-admin/frontend/internal/tailwind-config"
  },
  "transformOptions": {
    "babel": {
      "plugins": []
    }
  }
})

/** @type {import("/Users/william/work/my_opensource/goframe-vben-admin/frontend/internal/tailwind-config/src/index.js")} */
module.exports = jiti("/Users/william/work/my_opensource/goframe-vben-admin/frontend/internal/tailwind-config/src/index.ts")