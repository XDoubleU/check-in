import { Injectable } from "@nestjs/common"
import { WsService } from "../ws/ws.service"
import { LocationsService } from "../locations/locations.service"
import { EntityRepository } from "@mikro-orm/core"
import {
  CheckInEntity,
  type LocationEntity,
  type SchoolEntity
} from "../entities"
import { InjectRepository } from "@mikro-orm/nestjs"

@Injectable()
export class CheckInsService {
  private readonly checkInsRepository: EntityRepository<CheckInEntity>
  private readonly wsService: WsService
  private readonly locationsService: LocationsService

  public constructor(
    @InjectRepository(CheckInEntity)
    checkInsRepository: EntityRepository<CheckInEntity>,
    wsService: WsService,
    locationsService: LocationsService
  ) {
    this.checkInsRepository = checkInsRepository
    this.wsService = wsService
    this.locationsService = locationsService
  }

  public async getAll(
    location: LocationEntity,
    startDate: Date,
    endDate: Date
  ): Promise<CheckInEntity[]> {
    return this.checkInsRepository.find({
      location: {
        id: location.id
      },
      createdAt: {
        $gte: startDate,
        $lte: endDate
      }
    })
  }

  public async create(
    location: LocationEntity,
    school: SchoolEntity
  ): Promise<CheckInEntity> {
    const checkIn = new CheckInEntity(location, school)
    await this.checkInsRepository.persistAndFlush(checkIn)

    const updatedLocation = await this.locationsService.refresh(location.id)

    this.wsService.addLocationUpdate(updatedLocation)

    return checkIn
  }
}
