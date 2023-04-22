module.exports = {
  extends: ["custom"],
  parserOptions: {
    project: ["./tsconfig.json", "./tsconfig.test.json"],
    tsconfigRootDir: __dirname,
    sourceType: "module"
  },
  ignorePatterns: ["migration-files"]
}
