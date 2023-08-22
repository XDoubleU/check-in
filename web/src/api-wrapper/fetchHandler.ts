import { type APIResponse } from "./types"
import queryString from "query-string"
import { type ErrorDto } from "./types/apiTypes"
import { refreshTokens } from "./auth"

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function stringifyReplacer(_key: any, value: any): any {
  if (typeof value === "boolean") value = Boolean(value)
  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
  else if (!isNaN(value)) value = Number(value)
  return value
}

export async function fetchHandler<T = undefined>(
  input: string,
  method?: string,
  body?: unknown,
  query?: queryString.StringifiableRecord
): Promise<APIResponse<T>> {
  return fetchHandlerBase(true, input, method, body, query)
}

export async function fetchHandlerNoRefresh<T = undefined>(
  input: string,
  method?: string,
  body?: unknown,
  query?: queryString.StringifiableRecord
): Promise<APIResponse<T>> {
  return fetchHandlerBase(false, input, method, body, query)
}

// eslint-disable-next-line max-lines-per-function
async function fetchHandlerBase<T = undefined>(
  refresh: boolean,
  input: string,
  method?: string,
  body?: unknown,
  query: queryString.StringifiableRecord = {}
): Promise<APIResponse<T>> {
  const init: RequestInit = {}

  if (method) {
    init.method = method
  }

  if (body) {
    init.body = JSON.stringify(body, stringifyReplacer)
  }

  const url = queryString.stringifyUrl(
    {
      url: `${process.env.NEXT_PUBLIC_API_URL ?? ""}/${input}`,
      query: query
    },
    { skipEmptyString: true, arrayFormat: "comma" }
  )

  const fetchCall = async (): Promise<Response> => {
    return await fetch(url, {
      credentials: "include",
      headers: {
        "Content-Type": "application/json"
      },
      ...init
    })
  }

  // try the original call
  let response = await fetchCall()

  // when receiving 401 the accessToken is probably expired
  if (refresh && response.status === 401) {
    // refresh tokens
    const refreshResponse = await refreshTokens()
    // if refreshing the tokens returns 200, retry the original call
    if (refreshResponse.status === 200) {
      response = await fetchCall()
    }
  }

  if (!response.ok) {
    return {
      ok: response.ok,
      message: ((await response.json()) as ErrorDto).message
    }
  }

  const data = await response.text()
  if (data) {
    try {
      return {
        ok: response.ok,
        data: JSON.parse(data) as T
      }
    } catch (e) {
      return {
        ok: response.ok,
        data: data as T
      }
    }
  }

  return {
    ok: response.ok
  }
}
