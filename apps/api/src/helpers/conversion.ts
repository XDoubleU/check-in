import { type LocationEntity } from "mikro-orm-config"
import { type LocationUpdateEventDto } from "types-custom"

export function convertToLocationUpdateEventDto(
  location: LocationEntity
): LocationUpdateEventDto {
  return {
    normalizedName: location.normalizedName,
    available: location.available,
    capacity: location.capacity,
    yesterdayFullAt: location.yesterdayFullAt
  }
}
