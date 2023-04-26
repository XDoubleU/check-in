import { Injectable } from "@nestjs/common"
import { EntityRepository } from "@mikro-orm/core"
import { InjectRepository } from "@mikro-orm/nestjs"
import { Role } from "types-custom"
import { UserEntity } from "../entities"

@Injectable()
export class UsersService {
  private readonly usersRepository: EntityRepository<UserEntity>

  public constructor(
    @InjectRepository(UserEntity)
    usersRepository: EntityRepository<UserEntity>
  ) {
    this.usersRepository = usersRepository
  }

  public async getManagerCount(): Promise<number> {
    return this.usersRepository.count({
      roles: {
        $contains: [Role.Manager]
      }
    })
  }

  public async getAllManagersPaged(
    page: number,
    pageSize: number
  ): Promise<UserEntity[]> {
    return this.usersRepository.find(
      {
        roles: {
          $contains: [Role.Manager]
        }
      },
      {
        orderBy: {
          username: "asc"
        },
        limit: pageSize,
        offset: (page - 1) * pageSize
      }
    )
  }

  public async getById(id: string): Promise<UserEntity | null> {
    return await this.usersRepository.findOne({
      id: id
    })
  }

  public async getByUserName(username: string): Promise<UserEntity | null> {
    return await this.usersRepository.findOne({
      username: username
    })
  }

  public async create(
    username: string,
    password: string,
    role: Role
  ): Promise<UserEntity> {
    const user = new UserEntity(username, password, role)
    await this.usersRepository.persistAndFlush(user)
    return user
  }

  public async update(
    user: UserEntity,
    username?: string,
    password?: string
  ): Promise<UserEntity> {
    user.update(username, password)
    await this.usersRepository.flush()
    return user
  }

  public async delete(user: UserEntity): Promise<UserEntity> {
    await this.usersRepository.removeAndFlush(user)
    return user
  }
}
