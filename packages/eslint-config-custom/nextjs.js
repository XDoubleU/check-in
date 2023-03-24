const { commonNamingConvention } = require("./shared");

module.exports = {
  extends: [
    "next/core-web-vitals",
    "./index.js"
  ],
  rules: {
    "@typescript-eslint/explicit-function-return-type": "off",
    "@typescript-eslint/naming-convention": [
      "error",
      ...commonNamingConvention,
      {
        "selector": "function",
        "format": ["PascalCase"]
      }
    ],
    "@typescript-eslint/no-misused-promises": [
      "error",
      {
        "checksVoidReturn": {
          "attributes": false
        }
      }
    ]
  }
}