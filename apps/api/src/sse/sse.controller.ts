import { Controller, Param, Sse } from "@nestjs/common"
import { LocationUpdateEvent, SseService } from "./sse.service"
import { Observable } from "rxjs"
import { Public } from "../auth/decorators/public.decorator"

@Controller("sse")
export class SseController {
  constructor(private readonly sseService: SseService) {}

  @Public()
  @Sse()
  sseAllLocations(): Observable<LocationUpdateEvent> {
    return this.sseService.sendAllLocationUpdates()
  }

  @Sse(":normalizedName")
  sseSingleLocation(@Param("normalizedName") normalizedName: string): Observable<LocationUpdateEvent> {
    return this.sseService.sendSingleLocationUpdates(normalizedName)
  }
}
