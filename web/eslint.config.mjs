import { fixupConfigRules } from "@eslint/compat";
import _import from "eslint-plugin-import";
import sonarjs from "eslint-plugin-sonarjs";
import redundantUndefined from "eslint-plugin-redundant-undefined";
import jestDom from "eslint-plugin-jest-dom";
import testingLibrary from "eslint-plugin-testing-library";
import tsParser from "@typescript-eslint/parser";
import path from "node:path";
import { fileURLToPath } from "node:url";
import js from "@eslint/js";
import { FlatCompat } from "@eslint/eslintrc";
import tseslint from 'typescript-eslint';
import importPlugin from 'eslint-plugin-import';
import eslintConfigPrettier from "eslint-config-prettier"

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

export default [{
    ignores: [
        ".next",
        ".yarn",
        "eslint.config.mjs",
        "out",
        "public",
        "*.config.*",
        "coverage",
        "node_modules",
        "**/schema.d.ts",
        "jest.setup.js",
    ],
},
js.configs.recommended,
...tseslint.configs.strictTypeChecked,
...tseslint.configs.stylisticTypeChecked,
importPlugin.flatConfigs.typescript,
sonarjs.configs.recommended,
jestDom.configs["flat/recommended"],
testingLibrary.configs["flat/react"],
eslintConfigPrettier,
...fixupConfigRules(compat.extends(
    "next/core-web-vitals",
)), {
    plugins: {
        "redundant-undefined": redundantUndefined,
    },

    languageOptions: {
        parser: tsParser,
        ecmaVersion: 5,
        sourceType: "module",

        parserOptions: {
            project: ["./tsconfig.build.json", "./tsconfig.test.json"],
            tsconfigRootDir: __dirname,
            sourceType: "module"
        },
    },

    settings: {
        "import/resolver": {
            typescript: {
                alwaysTryTypes: true,
                project: ["./tsconfig.build.json", "./tsconfig.test.json"],
            },
        },
    },

    rules: {
        "@typescript-eslint/no-deprecated": "error",
        "@typescript-eslint/explicit-member-accessibility": "error",
        "@typescript-eslint/member-ordering": "error",
        "@typescript-eslint/no-require-imports": "error",
        "@typescript-eslint/parameter-properties": "error",
        "@typescript-eslint/prefer-readonly": "error",
        "max-lines-per-function": "error",
        "no-duplicate-imports": "error",
        "no-warning-comments": "error",
        "redundant-undefined/redundant-undefined": "error",

        "@typescript-eslint/no-unused-vars": ["error", {
            ignoreRestSiblings: true,
        }],

        "@typescript-eslint/consistent-type-imports": ["error", {
            fixStyle: "inline-type-imports",
        }],

        "import/consistent-type-specifier-style": ["error", "prefer-inline"],

        "@typescript-eslint/naming-convention": ["error", {
            selector: "default",
            format: ["camelCase"],
            leadingUnderscore: "allow",
            trailingUnderscore: "allow",
        }, {
            selector: "variable",
            format: ["camelCase", "UPPER_CASE"],
            leadingUnderscore: "allow",
            trailingUnderscore: "allow",
        }, {
            selector: "typeLike",
            format: ["PascalCase"],
        }, {
            selector: "objectLiteralProperty",
            format: null,
        }, {
            selector: "import",
            format: null,
        }, {
            selector: "variable",
            modifiers: ["const", "exported"],
            format: ["PascalCase", "UPPER_CASE"],
        }, {
            selector: "enumMember",
            format: ["PascalCase"],
        }, {
            selector: "function",
            format: ["PascalCase", "camelCase"],
        }],

        "@typescript-eslint/no-misused-promises": ["error", {
            checksVoidReturn: {
                attributes: false,
            },
        }],

        "testing-library/no-await-sync-events": ["error", {
            eventModules: ["fire-event"],
        }],
    },
}, {
    files: ["tests/**"],

    rules: {
        "max-lines-per-function": "off",
        "sonarjs/no-duplicate-string": "off",
    },
}];