import { EntityRepository, QueryOrder } from "@mikro-orm/core"
import { InjectRepository } from "@mikro-orm/nestjs"
import { Injectable } from "@nestjs/common"
import { LocationEntity, UserEntity } from "mikro-orm-config"
import { SseService } from "../sse/sse.service"

@Injectable()
export class LocationsService {
  constructor(
    @InjectRepository(LocationEntity)
    private readonly locationsRepository: EntityRepository<LocationEntity>,
    private readonly sseService: SseService
  ) {}

  async getTotalCount(): Promise<number> {
    return await this.locationsRepository.count()
  }
  
  async getAll(): Promise<LocationEntity[]> {
    return await this.locationsRepository.findAll()
  }

  async getAllPaged(page: number, pageSize: number): Promise<LocationEntity[]> {
    return await this.locationsRepository.findAll({
      orderBy: {
        name: QueryOrder.ASC
      },
      limit: pageSize,
      offset: (page - 1) * pageSize,
    })
  }

  async getById(id: string): Promise<LocationEntity | null> {
    return await this.locationsRepository.findOne({
      id: id
    })
  }

  async getByName(name: string): Promise<LocationEntity | null> {
    return await this.locationsRepository.findOne({
      name: name
    })
  }

  async create(name: string, capacity: number, user: UserEntity): Promise<LocationEntity> {
    const location = new LocationEntity(name, capacity, user)
    await this.locationsRepository.persistAndFlush(location)
    return location
  }

  async update(location: LocationEntity, name?: string, capacity?: number): Promise<LocationEntity> {
    location.name = name ?? location.name
    location.capacity = capacity ?? location.capacity

    await this.locationsRepository.flush()

    this.sseService.addLocationUpdate(location)

    return location
  }

  async delete(location: LocationEntity): Promise<LocationEntity> {
    await this.locationsRepository.removeAndFlush(location)
    return location
  }
}
