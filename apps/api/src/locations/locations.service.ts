import { Injectable } from "@nestjs/common"
import { BaseLocation, CheckIn, Location, User } from "types"
import { PrismaService } from "../prisma.service"
import { SseService } from "../sse/sse.service"

@Injectable()
export class LocationsService extends PrismaService {
  constructor(private readonly sseService: SseService) {
    super()
  }

  async getTotalCount(): Promise<number> {
    return await this.location.count()
  }
  
  async getAll(): Promise<Location[]> {
    const locations = await this.location.findMany({
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      }
    })

    return this.computeLocations(locations)
  }

  async getAllPaged(page: number, pageSize: number): Promise<Location[]> {
    const locations = await this.location.findMany({
      orderBy: {
        name: "asc"
      },
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      },
      take: pageSize,
      skip: (page - 1) * pageSize
    })

    return this.computeLocations(locations)
  }

  async getById(id: string): Promise<Location | null> {
    const location = await this.location.findFirst({
      where: {
        id: id
      },
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      }
    })

    if (!location) {
      return null
    }

    return this.computeLocation(location)
  }

  async getByName(name: string): Promise<Location | null> {
    const location = await this.location.findFirst({
      where: {
        name: name
      },
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      }
    })

    if (!location) {
      return null
    }

    return this.computeLocation(location)
  }

  async create(name: string, capacity: number, user: User): Promise<Location> {
    const location = await this.location.create({
      data: {
        name: name,
        capacity: capacity, 
        userId: user.id
      },
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      }
    })

    return this.computeLocation(location)
  }

  async update(location: Location, name?: string, capacity?: number): Promise<Location> {
    const result = await this.location.update({
      where: {
        id: location.id
      },
      data: {
        name: name,
        capacity: capacity
      },
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      }
    })

    const computedLocation = await this.computeLocation(result)
    this.sseService.addLocationUpdate(computedLocation)
    return computedLocation
  }

  async delete(location: Location): Promise<Location> {
    const result = await this.location.delete({
      where: {
        id: location.id
      },
      select: {
        id: true,
        name: true,
        capacity: true,
        user: {
          select: {
            id: true,
            username: true
          }
        }
      }
    })

    return this.computeLocation(result)
  }

  private getDates(): Date[] {
    const today = new Date()
    today.setHours(0)
    today.setMinutes(0)
    today.setSeconds(0)

    const tomorrow = new Date(today)
    tomorrow.setDate(tomorrow.getDate() + 1)

    return [today, tomorrow]
  }

  private transformLocation(location: BaseLocation, checkInsToday: CheckIn[]): Location {
    const normalizedName = location.name
                            .toLowerCase()
                            .replace(" ", "-")
                            .replace("[^A-Za-z0-9\-]+", "")

    return {
      id: location.id,
      name: location.name,
      normalizedName,
      available: location.capacity - checkInsToday.length,
      capacity: location.capacity,
      user: location.user
    }
  }

  private async computeLocation(location: BaseLocation): Promise<Location> {
    const [today, tomorrow] = this.getDates()

    const checkInsToday = await this.checkIn.findMany({
      where: {
        locationId: location.id,
        datetime: {
          gte: today,
          lt:  tomorrow
        }
      }
    })

    return this.transformLocation(location, checkInsToday)
  }

  private async computeLocations(locations: BaseLocation[]): Promise<Location[]> {
    const [today, tomorrow] = this.getDates()

    const checkInsToday = await this.checkIn.findMany({
      where: {
        locationId: {
          in: locations.map(location => location.id)
        },
        datetime: {
          gte: today,
          lt:  tomorrow
        }
      }
    })

    const result: Location[] = []
    locations.forEach((location) => {
      const checkIns = checkInsToday.filter(checkIn => {
        return checkIn.locationId === location.id
      })
      result.push(this.transformLocation(location, checkIns))
    })

    return result
  }
}
