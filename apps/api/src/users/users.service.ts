import { Injectable } from "@nestjs/common"
import { BaseUser, User } from "types-custom"
import { compareSync, hashSync } from "bcrypt"
import { PrismaService } from "../prisma.service"

@Injectable()
export class UsersService extends PrismaService {
  async getById(id: string): Promise<User | null> {
    const user = await this.user.findFirst({
      where: {
        id: id
      },
      select: {
        id: true,
        username: true,
        roles: true,
        location: {
          select: {
            id: true
          }
        }
      }
    })

    if (!user) {
      return null
    }

    return this.computeUser(user)
  }

  async getByUserName(username: string): Promise<User | null> {
    const user = await this.user.findFirst({
      where: {
        username: username
      },
      select: {
        id: true,
        username: true,
        roles: true,
        location: {
          select: {
            id: true
          }
        }
      }
    })

    if (!user) {
      return null
    }

    return this.computeUser(user)
  }

  async checkPassword(username: string, password: string): Promise<boolean> {
    const user = await this.user.findFirst({
      where: {
        username: username
      }
    })

    if (!user) {
      return false
    }

    return compareSync(password, user.passwordHash)
  }

  async create(username: string, password: string): Promise<User> {
    const passwordHash = hashSync(password, 12)

    return await this.user.create({
      data: {
        username: username,
        passwordHash: passwordHash
      }
    })
  }

  async update(user: User, username?: string, password?: string): Promise<User> {
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

  private computeUser(user: BaseUser): User {
    return {
      id: user.id,
      username: user.username,
      roles: user.roles,
      locationId: user.location?.id
    }
  }
}
