/** @type {import('ts-jest').JestConfigWithTsJest} */

const nextJest = require('next/jest')

const createJestConfig = nextJest({
  dir: "./",
});

const config = {
  preset: "ts-jest",
  moduleDirectories: ["node_modules", "src", "utils", __dirname],
  testEnvironment: "jest-environment-jsdom",
  testMatch: ["<rootDir>/tests/**/*.test.tsx"],
  collectCoverageFrom: ["<rootDir>/src/**/*.{ts,tsx}"],
  setupFilesAfterEnv: ["<rootDir>/jest.setup.js"],
  coveragePathIgnorePatterns: [
    "_app.tsx",
    "_document.tsx",
    "layouts/Error.tsx"
  ]
};

module.exports =  async () => ({
  ...(await createJestConfig(config)()),
  transformIgnorePatterns: [
    'node_modules/(?!(query-string|decode-uri-component|split-on-first|filter-obj)/)'
  ]
})

