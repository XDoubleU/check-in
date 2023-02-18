module.exports = {
  parser: "@typescript-eslint/parser",
  plugins: [
    "@typescript-eslint"
  ],
  extends: [
    "turbo",
    "plugin:@typescript-eslint/recommended"
  ],
  root: true,
  env: {
    "node": true,
    "jest": true
  },
  ignorePatterns: [".eslintrc.js", "*.config.*", "**/dist/**", "**test**"],
  rules: {
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        "ignoreRestSiblings": true
      }
    ],
    "@typescript-eslint/no-explicit-any": "error",
    "@typescript-eslint/explicit-function-return-type": "error",
    "sort-imports": 
    [
      "warn", 
      { 
        "ignoreCase": true, 
        "ignoreDeclarationSort": true 
      }
    ],
    "semi": [2, "never"],
    "quotes": [2, "double"]
  }
}