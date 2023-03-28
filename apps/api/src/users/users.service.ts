import { Injectable } from "@nestjs/common"
import { EntityRepository } from "@mikro-orm/core"
import { UserEntity } from "mikro-orm-config"
import { InjectRepository } from "@mikro-orm/nestjs"
import { Role } from "types-custom"

@Injectable()
export class UsersService {
  private readonly usersRepository: EntityRepository<UserEntity>

  public constructor(
    @InjectRepository(UserEntity)
    usersRepository: EntityRepository<UserEntity>
  ) {
    this.usersRepository = usersRepository
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

  public async createAdmin(username: string, password: string): Promise<UserEntity> {
    const user = new UserEntity(username, password, Role.Admin)
    await this.usersRepository.persistAndFlush(user)
    return user
  }

  public async create(username: string, password: string): Promise<UserEntity> {
    const user = new UserEntity(username, password)
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

  public async delete(user: UserEntity): Promise<UserEntity | null> {
    await this.usersRepository.removeAndFlush(user)
    return user
  }
}
