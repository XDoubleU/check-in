import {
  MessageBody,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer
} from "@nestjs/websockets"
import { Observable } from "rxjs"
import { type LocationUpdateEventDto } from "types-custom"
import { WebSocketServer as Server } from "ws"
import { WsService } from "./ws.service"

@WebSocketGateway({
  cors: {
    origin: "*"
  }
})
export class WsGateway {
  @WebSocketServer()
  public server!: Server

  private readonly wsService: WsService

  public constructor(wsService: WsService) {
    this.wsService = wsService
  }

  @SubscribeMessage("all-locations")
  public wsAllLocations(): Observable<LocationUpdateEventDto> {
    return this.wsService.sendAllLocationUpdates()
  }

  @SubscribeMessage("single-location")
  public wsSingleLocation(
    @MessageBody("normalizedName") normalizedName: string
  ): Observable<LocationUpdateEventDto> {
    return this.wsService.sendSingleLocationUpdates(normalizedName)
  }
}
