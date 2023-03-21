import { type Location } from "types-custom"

const SSE_URL = `${process.env.NEXT_PUBLIC_API_URL ?? ""}/sse`

export function checkinsEventSource(location: Location): EventSource {
  return new EventSource(`${SSE_URL}/${location.normalizedName}`, {
    withCredentials: true
  })
}
