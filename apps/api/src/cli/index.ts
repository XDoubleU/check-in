import { Command } from "commander"
import prompts from "prompts"
import { Role } from "types-custom"
import config, { CheckInEntity, LocationEntity, SchoolEntity, UserEntity } from "mikro-orm-config"
import { MikroORM } from "@mikro-orm/core"

const mikroOptions = {
  ...config,
  entities: [ CheckInEntity, LocationEntity, SchoolEntity, UserEntity ]
}

const program = new Command()

program
  .version("1.0.0")
  .description("A CLI tool for Check-In")

program.command("createadmin")
  .description("Creates an admin user")
  .option("-u, --username <string>", "username")
  .option("-p, --password <string>", "password")
  .action(async (options: { username?: string, password?: string }) => {
    const em = (await MikroORM.init(mikroOptions)).em

    let promptResponse = {
      username: "",
      password: ""
    }

    if (options.username && options.password){
      promptResponse = {
        username: options.username,
        password: options.password
      }
    } else {
      promptResponse = await prompts([
        {
          type: "text",
          name: "username",
          message: "Username?"
        },
        {
          type: "password",
          name: "password",
          message: "Password?"
        }
      ])
    }

    const existingUser = await em.findOne(UserEntity, { username: promptResponse.username })

    if (existingUser) {
      console.log("This username is already used")
      return
    }

    const user = new UserEntity(promptResponse.username, promptResponse.password, Role.Admin)
    await em.persistAndFlush(user)

    console.log("Admin added")
    process.exit(0)
  })

program.parse()