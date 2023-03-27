import Query from "../src/query"

describe("Query class", () => {
  it("All query parameters are undefined", () => {
    const page = undefined
    const search = undefined

    const query = new Query({
      page,
      search
    })

    expect(query.toString()).toBe("")
  })

  it("All query string parameters are empty strings", () => {
    const page = undefined
    const search = ""

    const query = new Query({
      page,
      search
    })

    expect(query.toString()).toBe("")
  })

  it("All query parameters have valid values", () => {
    const page = 5
    const search = "test"

    const query = new Query({
      page,
      search
    })

    expect(query.toString()).toBe("?page=5&search=test")
  })
})