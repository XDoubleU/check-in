import { Injectable } from "@nestjs/common"
import { filter, map, type Observable, Subject } from "rxjs"
import { type LocationUpdateEventDto } from "types-custom"
import { convertToLocationUpdateEventDto } from "../helpers/conversion"
import { type LocationEntity } from "../entities"

@Injectable()
export class WsService {
  private readonly locationUpdates = new Subject<LocationUpdateEventDto>()

  public addLocationUpdate(location: LocationEntity): void {
    this.locationUpdates.next(convertToLocationUpdateEventDto(location))
  }

  public sendAllLocationUpdates(): Observable<LocationUpdateEventDto> {
    return this.locationUpdates.asObservable()
  }

  public sendSingleLocationUpdates(
    normalizedName: string
  ): Observable<LocationUpdateEventDto> {
    return this.locationUpdates.asObservable().pipe(
      filter((location) => location.normalizedName === normalizedName),
      map((location) => location)
    )
  }
}
