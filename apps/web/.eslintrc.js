module.exports = {
  extends: ["custom/nextjs.js"],
  parserOptions: {
    project: "tsconfig.json",
    tsconfigRootDir: __dirname,
    sourceType: "module"
  }
}
