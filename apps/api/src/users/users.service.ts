import { Injectable } from "@nestjs/common"
import { InjectRepository } from "@mikro-orm/nestjs"
import { EntityRepository } from "@mikro-orm/core"
import { UserEntity } from "mikro-orm-config"

@Injectable()
export class UsersService {
  constructor(
    @InjectRepository(UserEntity)
    private readonly usersRepository: EntityRepository<UserEntity>
  ) {}

  async getById(id: string): Promise<UserEntity | null> {
    return await this.usersRepository.findOne({
      id: id
    })
  }

  async getByUserName(username: string): Promise<UserEntity | null> {
    return await this.usersRepository.findOne({
      username: username
    })
  }

  async create(username: string, password: string): Promise<UserEntity> {
    const user = new UserEntity(username, password)
    await this.usersRepository.persistAndFlush(user)
    return user
  }

  async update(user: UserEntity, username?: string, password?: string): Promise<UserEntity> {
    user.update(username, password)
    await this.usersRepository.flush()
    return user
  }

  async delete(user: UserEntity): Promise<UserEntity | null> {
    await this.usersRepository.removeAndFlush(user)
    return user
  }
}
