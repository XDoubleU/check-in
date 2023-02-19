module.exports = {
  extends: [
    "next/core-web-vitals",
    "custom"
  ],
  parserOptions: {
    project: "tsconfig.json",
    tsconfigRootDir: __dirname,
    sourceType: "module"
  },
  overrides: [
    {
      files: [ "*" ],
      rules: {
        "@typescript-eslint/explicit-function-return-type": "off",
      }
    }
  ]
}