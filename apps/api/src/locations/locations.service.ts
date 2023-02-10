import { Injectable } from "@nestjs/common"
import { Location, User } from "@prisma/client"
import { PrismaService } from "src/prisma.service"

@Injectable()
export class LocationsService extends PrismaService {
  async getTotalCount(): Promise<number> {
    return await this.location.count()
  }
  
  async getAll(page: number, pageSize: number): Promise<Location[]> {
    return await this.location.findMany({
      orderBy: {
        name: "asc"
      },
      take: pageSize,
      skip: (page - 1) * pageSize
    })
  }

  async getById(id: string): Promise<Location | null> {
    return await this.location.findFirst({
      where: {
        id: id
      }
    })
  }

  async getByName(name: string): Promise<Location | null> {
    return await this.location.findFirst({
      where: {
        name: name
      }
    })
  }

  async create(name: string, capacity: number, user: User): Promise<Location | null> {
    return await this.location.create({
      data: {
        name: name,
        capacity: capacity, 
        userId: user.id
      }
    })
  }

  async update(location: Location, name?: string, capacity?: number): Promise<Location | null> {
    return await this.location.update({
      where: {
        id: location.id
      },
      data: {
        name: name,
        capacity: capacity
      }
    })
  }

  async delete(location: Location): Promise<Location | null> {
    return await this.location.delete({
      where: {
        id: location.id
      }
    })
  }
}
