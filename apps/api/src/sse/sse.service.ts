import { Injectable } from "@nestjs/common"
import { LocationEntity } from "mikro-orm-config"
import { filter, map, Observable, Subject } from "rxjs"

export interface LocationUpdateEventData {
  normalizedName: string,
  available: number,
  capacity: number
}

export interface LocationUpdateEvent {
  data: LocationUpdateEventData
}

@Injectable()
export class SseService {
  private locationUpdates = new Subject<LocationUpdateEvent>()

  addLocationUpdate(location: LocationEntity): void {
    const newLocationUpdate: LocationUpdateEvent = {
      data: {
        normalizedName: location.normalizedName,
        available: location.available,
        capacity: location.capacity
      }
    }

    this.locationUpdates.next(newLocationUpdate)
  }

  sendAllLocationUpdates(): Observable<LocationUpdateEvent> {
    return this.locationUpdates.asObservable()
  }

  sendSingleLocationUpdates(normalizedName: string): Observable<LocationUpdateEvent> {
    return this.locationUpdates.asObservable().pipe(
      filter(location => location.data.normalizedName === normalizedName),
      map(location => location)
    )
  }
}
