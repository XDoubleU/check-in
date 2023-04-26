/* eslint-disable sonarjs/no-duplicate-string */
/* eslint-disable max-lines-per-function */

import { normalizeName } from "../src/helpers/normalization"

describe("Location Name Normalization (unit)", () => {
  it("Valid simple normalization", () => {
    const name = "Test name $14"

    const output = normalizeName(name)

    expect(output).toBe("test-name-14")
  })

  it("Valid normalization with space at start", () => {
    const name = " Test name $14"

    const output = normalizeName(name)

    expect(output).toBe("test-name-14")
  })

  it("Valid normalization with space at end", () => {
    const name = "Test name $14 "

    const output = normalizeName(name)

    expect(output).toBe("test-name-14")
  })
})
