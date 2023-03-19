import { Injectable } from "@nestjs/common"
import { SseService } from "../sse/sse.service"
import { LocationsService } from "../locations/locations.service"
import { EntityRepository } from "@mikro-orm/core"
import {
  CheckInEntity,
  type LocationEntity,
  type SchoolEntity
} from "mikro-orm-config"
import { InjectRepository } from "@mikro-orm/nestjs"

@Injectable()
export class CheckInsService {
  private readonly checkInsRepository: EntityRepository<CheckInEntity>
  private readonly sseService: SseService
  private readonly locationsService: LocationsService

  public constructor(
    @InjectRepository(CheckInEntity)
    checkInsRepository: EntityRepository<CheckInEntity>,
    sseService: SseService,
    locationsService: LocationsService
  ) {
    this.checkInsRepository = checkInsRepository
    this.sseService = sseService
    this.locationsService = locationsService
  }

  public async create(
    location: LocationEntity,
    school: SchoolEntity
  ): Promise<CheckInEntity | null> {
    const checkIn = new CheckInEntity(location, school)
    await this.checkInsRepository.persistAndFlush(checkIn)

    const updatedLocation = await this.locationsService.getById(location.id)
    if (!updatedLocation) {
      return null
    }

    this.sseService.addLocationUpdate(updatedLocation)

    return checkIn
  }
}
