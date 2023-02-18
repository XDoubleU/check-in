import { Injectable } from "@nestjs/common"
import { User } from "types"
import { hashSync } from "bcrypt"
import { PrismaService } from "../prisma.service"

@Injectable()
export class UsersService extends PrismaService {
  async getById(id: string): Promise<User | null> {
    return await this.user.findFirst({
      where: {
        id: id
      }
    })
  }

  async getByUserName(username: string): Promise<User | null> {
    return await this.user.findFirst({
      where: {
        username: username
      }
    })
  }

  async create(username: string, password: string): Promise<User | null> {
    const passwordHash = hashSync(password, 12)

    return await this.user.create({
      data: {
        username: username,
        passwordHash: passwordHash
      }
    })
  }

  async update(user: User, username?: string, password?: string): Promise<User | null> {
    const passwordHash = password === undefined ? undefined : hashSync(password, 12)

    return await this.user.update({
      where: {
        id: user.id
      },
      data: {
        username: username,
        passwordHash: passwordHash
      }
    })
  }

  async delete(user: User): Promise<User | null> {
    return await this.user.delete({
      where: {
        id: user.id
      }
    })
  }
}
