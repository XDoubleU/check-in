import { Injectable } from "@nestjs/common"
import { filter, map, Observable, Subject } from "rxjs"
import { Location } from "types"

@Injectable()
export class SseService {
  private locationUpdates = new Subject<Location>()

  addLocationUpdate(location: Location): void {
    this.locationUpdates.next(location)
  }

  sendAllLocationUpdates(): Observable<Location> {
    return this.locationUpdates.asObservable()
  }

  sendSingleLocationUpdates(locationId: string): Observable<Location> {
    return this.locationUpdates.asObservable().pipe(
      filter(location => location.id === locationId),
      map(location => location)
    )
  }
}
