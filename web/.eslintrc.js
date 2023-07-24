module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  plugins: [
    "@typescript-eslint",
    "import",
    "sonarjs",
    "redundant-undefined",
    "deprecation"
  ],
  extends: [
    "next/core-web-vitals",
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "plugin:@typescript-eslint/recommended-requiring-type-checking",
    "plugin:@typescript-eslint/strict",
    "plugin:import/typescript",
    "plugin:sonarjs/recommended"
  ],
  ignorePatterns: [
    ".eslintrc.js",
    "**/dist/**",
    "*.config.*",
    "**/coverage/**",
    "**/node_modules/**",
    "schema.d.ts"
  ],
  parserOptions: {
    project: "./tsconfig.json",
    tsconfigRootDir: __dirname,
    sourceType: "module"
  },
  rules: {
    "deprecation/deprecation": "error",
    "@typescript-eslint/explicit-member-accessibility": "error",
    "@typescript-eslint/member-ordering": "error",
    "@typescript-eslint/no-require-imports": "error",
    "@typescript-eslint/parameter-properties": "error",
    "@typescript-eslint/prefer-readonly": "error",
    "max-lines-per-function": "error",
    "no-duplicate-imports": "error",
    "no-warning-comments": "error",
    "redundant-undefined/redundant-undefined": "error",
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        "ignoreRestSiblings": true
      }
    ],
    "@typescript-eslint/consistent-type-imports": [
      "error",
      {
        "fixStyle": "inline-type-imports"
      }
    ],
    "import/consistent-type-specifier-style": [
      "error", 
      "prefer-inline"
    ],
    "@typescript-eslint/naming-convention": [
      "error",
      {
        "selector": "default",
        "format": ["camelCase"],
        "leadingUnderscore": "allow",
        "trailingUnderscore": "allow",
      },
      {
        "selector": "variable",
        "format": ["camelCase", "UPPER_CASE"],
        "leadingUnderscore": "allow",
        "trailingUnderscore": "allow",
      },
      {
        "selector": "typeLike",
        "format": ["PascalCase"],
      },
      { 
        "selector": "objectLiteralProperty",
        "format": null
      },
      {
        "selector": "variable",
        "modifiers": ["const", "exported"],
        "format": ["PascalCase", "UPPER_CASE"]
      },
      {
        "selector": "enumMember",
        "format": ["PascalCase"]
      },
      {
        "selector": "function",
        "format": ["PascalCase", "camelCase"]
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
  },
  settings: {
    "import/resolver": {
      "typescript": {
        "alwaysTryTypes": true,
        "project": [
          "./tsconfig.json",
        ]
      }
    },
  }
}
