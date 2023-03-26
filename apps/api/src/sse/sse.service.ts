import { Injectable } from "@nestjs/common"
import { type LocationEntity } from "mikro-orm-config"
import { filter, map, type Observable, Subject } from "rxjs"
import { type LocationUpdateEventDto } from "types-custom"
import { convertToLocationUpdateEventDto } from "../helpers/conversion"

export interface LocationUpdateEvent {
  data: LocationUpdateEventDto
}

@Injectable()
export class SseService {
  private readonly locationUpdates = new Subject<LocationUpdateEvent>()

  public addLocationUpdate(location: LocationEntity): void {
    this.locationUpdates.next({
      data: convertToLocationUpdateEventDto(location)
    })
  }

  public sendAllLocationUpdates(): Observable<LocationUpdateEvent> {
    return this.locationUpdates.asObservable()
  }

  public sendSingleLocationUpdates(
    normalizedName: string
  ): Observable<LocationUpdateEvent> {
    return this.locationUpdates.asObservable().pipe(
      filter((location) => location.data.normalizedName === normalizedName),
      map((location) => location)
    )
  }
}
