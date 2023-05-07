import { EntityRepository, EntityManager } from "@mikro-orm/core"
import { InjectRepository } from "@mikro-orm/nestjs"
import { Injectable } from "@nestjs/common"
import { SchoolEntity } from "../entities/school"

@Injectable()
export class SchoolsService {
  private readonly em: EntityManager
  private readonly schoolsRepository: EntityRepository<SchoolEntity>

  public constructor(
    em: EntityManager,
    @InjectRepository(SchoolEntity)
    schoolsRepository: EntityRepository<SchoolEntity>
  ) {
    this.em = em
    this.schoolsRepository = schoolsRepository
  }

  public async getTotalCount(): Promise<number> {
    return this.schoolsRepository.count({
      id: {
        $ne: 1
      }
    })
  }

  public async getAll(): Promise<SchoolEntity[]> {
    return await this.schoolsRepository.findAll()
  }

  public async getAllForLocation(locationId: string): Promise<SchoolEntity[]> {
    const schools = await this.schoolsRepository.findAll({
      populate: ["checkIns.location.id"],
      orderBy: {
        name: "asc"
      }
    })

    schools.sort((a, b) => {
      if (a.id === 1) {
        return -1
      }

      if (b.id === 1) {
        return -1
      }

      const aLocationCheckIns = a.checkIns
        .toArray()
        .filter((checkIn) => checkIn.location.id === locationId).length

      const bLocationCheckIns = b.checkIns
        .toArray()
        .filter((checkIn) => checkIn.location.id === locationId).length

      return aLocationCheckIns > bLocationCheckIns ? -1 : 1
    })

    return schools
  }

  public async getAllPaged(
    page: number,
    pageSize: number
  ): Promise<SchoolEntity[]> {
    return this.schoolsRepository.find(
      {
        id: {
          $ne: 1
        }
      },
      {
        orderBy: {
          name: "asc"
        },
        limit: pageSize,
        offset: (page - 1) * pageSize
      }
    )
  }

  public async getById(id: number): Promise<SchoolEntity | null> {
    return await this.schoolsRepository.findOne({
      id: id
    })
  }

  public async getByName(name: string): Promise<SchoolEntity | null> {
    return await this.schoolsRepository.findOne({
      name: name
    })
  }

  public async create(name: string): Promise<SchoolEntity> {
    const school = new SchoolEntity(name)
    await this.em.persistAndFlush(school)
    return school
  }

  public async update(
    school: SchoolEntity,
    name: string
  ): Promise<SchoolEntity> {
    school.name = name
    await this.em.flush()
    return school
  }

  public async delete(school: SchoolEntity): Promise<SchoolEntity> {
    await this.em.removeAndFlush(school)
    return school
  }
}
