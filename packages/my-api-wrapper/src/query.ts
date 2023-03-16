export default class Query {  
  private readonly urlSearchParams : URLSearchParams

  public constructor(data: Record<string, number | string | undefined>) {
    this.urlSearchParams = new URLSearchParams()

    for (const key of Object.keys(data)){
      if (data[key] === undefined) {
        continue
      }

      if (typeof data[key] === "string" && data[key] !== ""){
        this.urlSearchParams.append(key, data[key] as string)
      }

      if (typeof data[key] === "number"){
        this.urlSearchParams.append(key, (data[key] as number).toString())
      }
    }
  }

  public toString(): string {
    const result = this.urlSearchParams.toString()

    if (result.length > 0){
      return `?${result}`
    }

    return result
  }
}