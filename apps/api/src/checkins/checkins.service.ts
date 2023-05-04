import { Injectable } from "@nestjs/common"
import { WsService } from "../ws/ws.service"
import { EntityRepository, EntityManager } from "@mikro-orm/core"
import {
  CheckInEntity,
  type LocationEntity,
  type SchoolEntity
} from "../entities"
import { InjectRepository } from "@mikro-orm/nestjs"

@Injectable()
export class CheckInsService {
  private readonly em: EntityManager
  private readonly checkInsRepository: EntityRepository<CheckInEntity>
  private readonly wsService: WsService

  public constructor(
    em: EntityManager,
    @InjectRepository(CheckInEntity)
    checkInsRepository: EntityRepository<CheckInEntity>,
    wsService: WsService
  ) {
    this.em = em
    this.checkInsRepository = checkInsRepository
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

    this.wsService.addLocationUpdate(location)

    return checkIn
  }
}
