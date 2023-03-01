import Query from "../src/query"

describe("Query class", () => {
  test("All query parameters are undefined", () => {
    const page = undefined
    const search = undefined
  
    const query = new Query({
      page,
      search
    })
  
    expect(query.toString()).toBe("")
  })

  test("All query string parameters are empty strings", () => {
    const page = undefined
    const search = ""
  
    const query = new Query({
      page,
      search
    })
  
    expect(query.toString()).toBe("")
  })

  test("All query parameters have valid values", () => {
    const page = 5
    const search = "test"
  
    const query = new Query({
      page,
      search
    })
  
    expect(query.toString()).toBe("?page=5&search=test")
  })
})

