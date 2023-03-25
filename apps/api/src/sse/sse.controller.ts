import { Controller, Param, Res, Sse } from "@nestjs/common"
import { type LocationUpdateEvent, SseService } from "./sse.service"
import { Observable } from "rxjs"
import { Public } from "../auth/decorators/public.decorator"
import { Response } from "express"

@Controller("sse")
export class SseController {
  private readonly sseService: SseService

  public constructor(sseService: SseService) {
    this.sseService = sseService
  }

  @Public()
  @Sse()
  public sseAllLocations(
    @Res() res: Response
  ): Observable<LocationUpdateEvent> {
    res.set("Access-Control-Allow-Origin", "*")
    return this.sseService.sendAllLocationUpdates()
  }

  @Sse(":normalizedName")
  public sseSingleLocation(
    @Param("normalizedName") normalizedName: string
  ): Observable<LocationUpdateEvent> {
    return this.sseService.sendSingleLocationUpdates(normalizedName)
  }
}
