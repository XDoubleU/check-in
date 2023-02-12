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
  ignorePatterns: [".eslintrc.js", "*.config.*", "**/build/**", "**test**"],
  rules: {
    "@typescript-eslint/no-unused-vars": "error",
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