export async function fetchHandler(input: URL | RequestInfo, init?: RequestInit): Promise<Response | null> {
  const fetchCall = async (): Promise<Response> => {
    return await fetch(input, {
      credentials: "include",
      ...init
    })
  }
  
  let response = await fetchCall()
  if (response.status === 401) {
    const refreshResponse = await refreshTokens()
    if (refreshResponse.status === 401) {
      return null
    }
    response = await fetchCall()
  } else if (response.status === 404) {
    return null
  }

  return response
}

async function refreshTokens(): Promise<Response> {
  const url = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/auth/refresh`
  const refreshResponse = await fetch(url, {
    credentials: "include"
  })

  return refreshResponse
}