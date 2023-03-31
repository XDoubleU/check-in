import { Controller, Header, Param, Sse } from "@nestjs/common"
import { type LocationUpdateEvent, SseService } from "./sse.service"
import { Observable } from "rxjs"
import { Public } from "../auth/decorators/public.decorator"

@Controller("sse")
export class SseController {
  private readonly sseService: SseService

  public constructor(sseService: SseService) {
    this.sseService = sseService
  }

  @Public()
  @Sse()
  @Header("Access-Control-Allow-Origin", "*")
  @Header("Content-Type", "text/event-stream")
  @Header("Cache-Control", "no-cache, no-store")
  public sseAllLocations(): Observable<LocationUpdateEvent> {
    return this.sseService.sendAllLocationUpdates()
  }

  @Sse(":normalizedName")
  public sseSingleLocation(
    @Param("normalizedName") normalizedName: string
  ): Observable<LocationUpdateEvent> {
    return this.sseService.sendSingleLocationUpdates(normalizedName)
  }
}
