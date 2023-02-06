import { Command } from "commander"
import prompts from "prompts"
import { hashSync } from "bcrypt"
import { PrismaClient } from "@prisma/client"

const prisma = new PrismaClient()
const program = new Command()

program
  .version("1.0.0")
  .description("A CLI tool for Check-In")

program.command("createadmin")
  .description("Creates an admin user")
  .action(async () => {
    const promptResponse = await prompts([
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

    const existingUser = await prisma.user.findUnique({
      where: {
        username: promptResponse.username
      }
    })

    if (existingUser !== null) {
      console.log("This username is already used")
      return
    }

    const passwordHash = hashSync(promptResponse.password, 10)
    const result = await prisma.user.create({
      data: {
        username: promptResponse.username,
        passwordHash: passwordHash,
        isAdmin: true
      }
    })

    if (result === null || result === undefined) {
      console.log("Something went wrong")
      return
    }

    console.log("Admin added")
    process.exit(0)
  })

program.parse()