import { URLSearchParams } from "url"

// TODO: test this
export default class Query {  
  urlSearchParams : URLSearchParams

  constructor(data: { [key: string]: number | string | undefined }) {
    for (const key of Object.keys(data)){
      if (data[key] === undefined || data[key] === null) {
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
    return this.urlSearchParams.toString()
  }
}