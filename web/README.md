# Check-In Web

## Lint Commands

Format code:        `yarn format`
Run linter:         `yarn lint`
Run autofix linter: `yarn lint:fix`

## Dev Commands

None; the web-client needs the API and thus the main docker-compose file should be used.

## Build Commands

Build:                                                `yarn build`
Export project as HTML (for running as static site):  `yarn export`

## Test Commands

For now: None

Name suggestions for future implementation:
Run tests:                `yarn test`
Run tests with coverage:  `yarn test:cov`

## Other Commands

Generate types from OpenAPI Spec: `yarn swag`

## CI Commands

Generate lint report:  `yarn lint:report`
Generate test report:  `yarn test:report`
