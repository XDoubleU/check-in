import { Injectable } from "@nestjs/common"
import { WsService } from "../ws/ws.service"
import { EntityRepository, EntityManager } from "@mikro-orm/core"
import {
  CheckInEntity,
  type LocationEntity,
  type SchoolEntity
} from "../entities"
import { InjectRepository } from "@mikro-orm/nestjs"
import { LocationsService } from "../locations/locations.service"

@Injectable()
export class CheckInsService {
  private readonly em: EntityManager
  private readonly checkInsRepository: EntityRepository<CheckInEntity>
  private readonly locationsService: LocationsService
  private readonly wsService: WsService

  public constructor(
    em: EntityManager,
    @InjectRepository(CheckInEntity)
    checkInsRepository: EntityRepository<CheckInEntity>,
    locationsService: LocationsService,
    wsService: WsService
  ) {
    this.em = em
    this.checkInsRepository = checkInsRepository
    this.locationsService = locationsService
    this.wsService = wsService
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
    await this.em.persistAndFlush(checkIn)

    // Need this line for recomputing available spots
    const updatedLocation = await this.locationsService.refresh(location.id)
    this.wsService.addLocationUpdate(updatedLocation)

    return checkIn
  }
}
