import { Command, Positional } from 'nestjs-command';
import { Injectable } from '@nestjs/common';
import { UsersService } from './users.service';

@Injectable()
export class UserCommand {
  private readonly usersService: UsersService
  public constructor(usersService: UsersService) {
    this.usersService = usersService
  }

  @Command({
    command: 'createadmin',
    describe: 'create an admin',
  })
  public async create(
    @Positional({
      name: 'username',
      describe: 'the username',
      type: 'string'
    })
    username: string,
    @Positional({
      name: 'password',
      describe: 'the password',
      type: 'string'
    })
    password: string,
  ): Promise<void> {
    await this.usersService.createAdmin(username, password)
  }
}