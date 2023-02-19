import { Injectable } from "@nestjs/common"
import { CheckIn, Location, School } from "types"
import { PrismaService } from "../prisma.service"
import { SseService } from "../sse/sse.service"
import { LocationsService } from "../locations/locations.service"

@Injectable()
export class CheckInsService extends PrismaService {
  constructor(
    private readonly sseService: SseService,
    private readonly locationsService: LocationsService
  ) 
  {
    super()
  }

  async create(location: Location, school: School): Promise<CheckIn | null> {
    const checkIn = await this.checkIn.create({
      data: {
        locationId: location.id,
        capacity: location.capacity,
        schoolId: school.id
      }
    })

    const updatedLocation = await this.locationsService.getById(location.id)
    if (!updatedLocation) {
      return null
    }

    this.sseService.addLocationUpdate(updatedLocation)

    return checkIn
  }
}
