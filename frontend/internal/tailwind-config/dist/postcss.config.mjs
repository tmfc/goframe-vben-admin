import { createJiti } from "../../../node_modules/.pnpm/jiti@2.6.1/node_modules/jiti/lib/jiti.mjs";

const jiti = createJiti(import.meta.url, {
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

/** @type {import("/Users/william/work/my_opensource/goframe-vben-admin/frontend/internal/tailwind-config/src/postcss.config.js")} */
const _module = await jiti.import("/Users/william/work/my_opensource/goframe-vben-admin/frontend/internal/tailwind-config/src/postcss.config.ts");

export default _module?.default ?? _module;