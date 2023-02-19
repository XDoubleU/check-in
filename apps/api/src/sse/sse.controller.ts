import { Controller, Param, Sse } from "@nestjs/common"
import { SseService } from "./sse.service"
import { Observable } from "rxjs"
import { Location } from "types"

@Controller("sse")
export class SseController {
  constructor(private readonly sseService: SseService) {}

  @Sse()
  sseAllLocations(): Observable<Location> {
    return this.sseService.sendAllLocationUpdates()
  }

  @Sse(":id")
  sseSingleLocation(@Param("id") id: string): Observable<Location> {
    return this.sseService.sendSingleLocationUpdates(id)
  }
}
