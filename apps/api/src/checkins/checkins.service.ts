import { Injectable } from "@nestjs/common"
import { CheckIn, Location, School } from "@prisma/client"
import { PrismaService } from "src/prisma.service"

@Injectable()
export class CheckInsService extends PrismaService {
  async create(location: Location, school: School): Promise<CheckIn | null> {
    return await this.checkIn.create({
      data: {
        locationId: location.id,
        capacity: location.capacity,
        schoolId: school.id
      }
    })
  }
}
