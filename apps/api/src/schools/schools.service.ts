import { EntityRepository } from "@mikro-orm/core"
import { InjectRepository } from "@mikro-orm/nestjs"
import { Injectable } from "@nestjs/common"
import { SchoolEntity } from "mikro-orm-config"

@Injectable()
export class SchoolsService {
  constructor(
    @InjectRepository(SchoolEntity)
    private readonly schoolsRepository: EntityRepository<SchoolEntity>
  ) {}

  async getTotalCount(): Promise<number> {
    return this.schoolsRepository.count()
  }
  
  async getAll(locationId?: string): Promise<SchoolEntity[]> {
    const schools = await this.schoolsRepository.findAll()

    if (!locationId) {
      return schools
    }

    schools.sort((a, b) => {
      let aLocationCheckIns = 0
      let bLocationCheckIns = 0

      if (a.id !== 1) {
        void a.checkIns.matching({
          where: {
            location: {
              id: locationId
            }
          }
        }).then(data => aLocationCheckIns = data.length)
      }

      if (b.id !== 1) {
        void b.checkIns.matching({
          where: {
            location: {
              id: locationId
            }
          }
        }).then(data => bLocationCheckIns = data.length)
      }

      return (aLocationCheckIns < bLocationCheckIns) ? 1 : -1
    })

    return schools
  }

  async getAllPaged(page: number, pageSize: number): Promise<SchoolEntity[]> {
    return this.schoolsRepository.findAll({
      orderBy: {
        name: "asc"
      },
      limit: pageSize,
      offset: (page - 1) * pageSize
    })
  }

  async getById(id: number): Promise<SchoolEntity | null> {
    return await this.schoolsRepository.findOne({
      id: id
    })
  }

  async getByName(name: string): Promise<SchoolEntity | null> {
    return await this.schoolsRepository.findOne({
      name: name
    })
  }

  async create(name: string): Promise<SchoolEntity> {
    const school = new SchoolEntity(name)
    await this.schoolsRepository.persistAndFlush(school)
    return school
  }

  async update(school: SchoolEntity, name: string): Promise<SchoolEntity> {
    school.name = name
    await this.schoolsRepository.flush()
    return school
  }

  async delete(school: SchoolEntity): Promise<SchoolEntity> {
    await this.schoolsRepository.removeAndFlush(school)
    return school
  }
}
