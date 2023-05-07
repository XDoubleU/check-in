import { type ErrorDto } from "types-custom"
import { type APIResponse } from "./types"
import Router from "next/router"
import queryString from "query-string"

export async function fetchHandler<T = undefined>(
  input: string,
  init?: RequestInit,
  query?: queryString.StringifiableRecord
): Promise<APIResponse<T>> {
  return fetchHandlerBase(true, input, init, query)
}

export async function fetchHandlerNoRefresh<T = undefined>(
  input: string,
  init?: RequestInit,
  query?: queryString.StringifiableRecord
): Promise<APIResponse<T>> {
  return fetchHandlerBase(false, input, init, query)
}

// eslint-disable-next-line max-lines-per-function
async function fetchHandlerBase<T = undefined>(
  refresh: boolean,
  input: string,
  init?: RequestInit,
  query: queryString.StringifiableRecord = {}
): Promise<APIResponse<T>> {
  const url = queryString.stringifyUrl(
    {
      url: `${process.env.NEXT_PUBLIC_API_URL ?? ""}/${input}`,
      query: query
    },
    { skipEmptyString: true }
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
    // if refreshing the tokens still gives a 401 the user isn't logged in
    if (refreshResponse.status === 401) {
      await Router.push("signin")
      response = refreshResponse
    }
    // else retry the original call
    else {
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

async function refreshTokens(): Promise<Response> {
  const url = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/auth/refresh`

  return await fetch(url, {
    credentials: "include"
  })
}
