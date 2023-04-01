module.exports = {
  extends: [
    "custom"
  ],
  parserOptions: {
    project: "./tsconfig.json",
    tsconfigRootDir: __dirname,
    sourceType: "module"
  },
  ignorePatterns: ["src/migrations"]
}