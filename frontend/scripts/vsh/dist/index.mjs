import { createJiti } from "../../../node_modules/.pnpm/jiti@2.6.1/node_modules/jiti/lib/jiti.mjs";

const jiti = createJiti(import.meta.url, {
  "interopDefault": true,
  "alias": {
    "@vben/vsh": "/Users/william/work/my_opensource/goframe-vben-admin/frontend/scripts/vsh"
  },
  "transformOptions": {
    "babel": {
      "plugins": []
    }
  }
})

/** @type {import("/Users/william/work/my_opensource/goframe-vben-admin/frontend/scripts/vsh/src/index.js")} */
const _module = await jiti.import("/Users/william/work/my_opensource/goframe-vben-admin/frontend/scripts/vsh/src/index.ts");

export default _module?.default ?? _module;