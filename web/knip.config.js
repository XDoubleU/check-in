module.exports = {
    "entry": ["src/pages/index.tsx"],
    "project": ["**/*.{ts,tsx}"],
    "ignore": [
        "src/api-wrapper/types/*",
    ],
    "ignoreDependencies": [
        "jest-environment-jsdom"
    ]
}