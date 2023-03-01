import { Injectable } from "@nestjs/common"
import { School } from "types"
import { PrismaService } from "../prisma.service"

@Injectable()
export class SchoolsService extends PrismaService {
  async getTotalCount(): Promise<number> {
    return await this.school.count()
  }
  
  async getAll(locationId: string | undefined): Promise<School[]> {
    const schools = await this.school.findMany()
    if (!locationId) {
      return schools
    }

    const checkins = await this.checkIn.findMany({
      where: {
        location: {
          id: locationId
        }
      }
    })
    
    function countCheckIns(school: School): number {
      return checkins.reduce((total,x) => (x.schoolId === school.id ? total+1 : total), 0)
    }
    
    schools.sort((schoolA, schoolB) => (countCheckIns(schoolA) > countCheckIns(schoolB)) ? 1 : -1)

    return schools
  }

  async getAllPaged(page: number, pageSize: number): Promise<School[]> {
    return await this.school.findMany({
      orderBy: {
        name: "asc"
      },
      take: pageSize,
      skip: (page - 1) * pageSize
    })
  }

  async getById(id: number): Promise<School | null> {
    return await this.school.findFirst({
      where: {
        id: id
      }
    })
  }

  async getByName(name: string): Promise<School | null> {
    return await this.school.findFirst({
      where: {
        name: name
      }
    })
  }

  async create(name: string): Promise<School> {
    return await this.school.create({
      data: {
        name: name
      }
    })
  }

  async update(school: School, name: string): Promise<School> {
    return await this.school.update({
      where: {
        id: school.id
      },
      data: {
        name: name
      }
    })
  }

  async delete(school: School): Promise<School> {
    return await this.school.delete({
      where: {
        id: school.id
      }
    })
  }
}
