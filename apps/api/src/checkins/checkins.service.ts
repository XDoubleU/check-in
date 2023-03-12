import { Injectable } from "@nestjs/common"
import { SseService } from "../sse/sse.service"
import { LocationsService } from "../locations/locations.service"
import { EntityRepository } from "@mikro-orm/core"
import { InjectRepository } from "@mikro-orm/nestjs"
import { CheckInEntity, LocationEntity, SchoolEntity } from "mikro-orm-config"

@Injectable()
export class CheckInsService {
  constructor(
    @InjectRepository(CheckInEntity)
    private readonly checkInsRepository: EntityRepository<CheckInEntity>,
    private readonly sseService: SseService,
    private readonly locationsService: LocationsService
  ) {}

  async create(location: LocationEntity, school: SchoolEntity): Promise<CheckInEntity | null> {
    const checkIn = new CheckInEntity(location, school)
    await this.checkInsRepository.persistAndFlush(checkIn)

    const updatedLocation = await this.locationsService.getById(location.id)
    if (!updatedLocation){
      return null
    }

    this.sseService.addLocationUpdate(updatedLocation)

    return checkIn
  }
}
