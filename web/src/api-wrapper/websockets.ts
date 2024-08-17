import { type SubscribeMessageDto, type Location } from "./types/apiTypes"

const WS_URL = process.env.NEXT_PUBLIC_API_URL?.replace("http", "ws") ?? ""

export function checkinsWebsocket(location: Location): WebSocket {
  const webSocket = new WebSocket(WS_URL)

  webSocket.onopen = (): void => {
    let message: SubscribeMessageDto = {
      subject: "single-location",
      normalizedName: location.normalizedName
    }
    webSocket.send(JSON.stringify(message))

    message = {
      subject: "state"
    }
    webSocket.send(JSON.stringify(message))
  }

  return webSocket
}
