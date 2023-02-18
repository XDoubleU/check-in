import { Injectable } from "@nestjs/common"
import { Location, User } from "types"
import { PrismaService } from "../prisma.service"

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
      include: {
        user: true
      },
      take: pageSize,
      skip: (page - 1) * pageSize
    })
  }

  async getById(id: string): Promise<Location | null> {
    return await this.location.findFirst({
      where: {
        id: id
      },
      include: {
        user: true
      }
    })
  }

  async getByUserId(userId: string): Promise<Location | null> {
    return await this.location.findFirst({
      where: {
        userId: userId
      },
      include: {
        user: true
      }
    })
  }

  async getByName(name: string): Promise<Location | null> {
    return await this.location.findFirst({
      where: {
        name: name
      },
      include: {
        user: true
      }
    })
  }

  async create(name: string, capacity: number, user: User): Promise<Location | null> {
    return await this.location.create({
      data: {
        name: name,
        capacity: capacity, 
        userId: user.id
      },
      include: {
        user: true
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
      },
      include: {
        user: true
      }
    })
  }

  async delete(location: Location): Promise<Location | null> {
    return await this.location.delete({
      where: {
        id: location.id
      },
      include: {
        user: true
      }
    })
  }
}
