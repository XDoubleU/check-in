import { type Location } from "types-custom"

// eslint-disable-next-line @typescript-eslint/no-non-null-assertion
const WS_URL = process.env.NEXT_PUBLIC_API_URL?.replace("http", "ws") ?? ""

export function checkinsWebsocket(location: Location): WebSocket {
  const webSocket = new WebSocket(WS_URL)

  webSocket.onopen = (): void => {
    webSocket.send(
      JSON.stringify({
        event: "single-location",
        data: {
          normalizedName: location.normalizedName
        }
      })
    )
  }

  return webSocket
}
