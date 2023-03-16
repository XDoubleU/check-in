const { commonNamingConvention } = require("./shared");

module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  plugins: [
    "@typescript-eslint",
    "import"
  ],
  extends: [
    "turbo",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/recommended-requiring-type-checking",
    "plugin:@typescript-eslint/strict"
  ],
  ignorePatterns: [
    ".eslintrc.js",
    "**/dist/**",
    "*.config.*"
  ],
  rules: {
    "@typescript-eslint/explicit-function-return-type": "error",
    "@typescript-eslint/explicit-member-accessibility": "error",
    "@typescript-eslint/member-ordering": "error",
    "@typescript-eslint/no-require-imports": "error",
    "@typescript-eslint/parameter-properties": "error",
    "@typescript-eslint/prefer-readonly": "error",
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        "ignoreRestSiblings": true
      }
    ],
    "@typescript-eslint/naming-convention": [
      "error",
      ...commonNamingConvention
    ],
    "@typescript-eslint/consistent-type-imports": [
      "error",
      {
        "fixStyle": "inline-type-imports"
      }
    ],
    "import/no-duplicates": [
      "error", 
      {
        "prefer-inline": true
      }
    ],
    "import/consistent-type-specifier-style": [
      "error", 
      "prefer-inline"
    ]
  }
}