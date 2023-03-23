import { EntityRepository, QueryOrder } from "@mikro-orm/core"
import { InjectRepository } from "@mikro-orm/nestjs"
import { Injectable } from "@nestjs/common"
import { LocationEntity, type UserEntity } from "mikro-orm-config"
import { SseService } from "../sse/sse.service"

@Injectable()
export class LocationsService {
  private readonly locationsRepository: EntityRepository<LocationEntity>
  private readonly sseService: SseService

  public constructor(
    @InjectRepository(LocationEntity)
    locationsRepository: EntityRepository<LocationEntity>,
    sseService: SseService
  ) {
    this.locationsRepository = locationsRepository
    this.sseService = sseService
  }

  public async getTotalCount(): Promise<number> {
    return await this.locationsRepository.count()
  }

  public async getAll(): Promise<LocationEntity[]> {
    return await this.locationsRepository.findAll()
  }

  public async getAllPaged(
    page: number,
    pageSize: number
  ): Promise<LocationEntity[]> {
    return await this.locationsRepository.findAll({
      orderBy: {
        name: QueryOrder.ASC
      },
      limit: pageSize,
      offset: (page - 1) * pageSize
    })
  }

  public async getById(id: string): Promise<LocationEntity | null> {
    return await this.locationsRepository.findOne({
      id: id
    })
  }

  public async refresh(locationId: string): Promise<LocationEntity> {
    return this.locationsRepository.findOneOrFail(
      { id: locationId },
      { refresh: true }
    )
  }

  public async getByName(name: string): Promise<LocationEntity | null> {
    return await this.locationsRepository.findOne({
      name: name
    })
  }

  public async create(
    name: string,
    capacity: number,
    user: UserEntity
  ): Promise<LocationEntity> {
    const location = new LocationEntity(name, capacity, user)
    await this.locationsRepository.persistAndFlush(location)
    return await this.refresh(location.id)
  }

  public async update(
    location: LocationEntity,
    name?: string,
    capacity?: number
  ): Promise<LocationEntity> {
    location.name = name ?? location.name
    location.capacity = capacity ?? location.capacity

    await this.locationsRepository.flush()

    const updatedLocation = await this.refresh(location.id)

    this.sseService.addLocationUpdate(updatedLocation)

    return location
  }

  public async delete(location: LocationEntity): Promise<LocationEntity> {
    await this.locationsRepository.removeAndFlush(location)
    return location
  }
}
