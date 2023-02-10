import { prisma } from "@/common/prisma"
import { hashSync } from "bcrypt"
import { NextApiRequest, NextApiResponse } from "next"

//TODO: cleaner ?
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === "GET") {
    const { id } = req.query

    const location = await prisma.location.findFirst({
      where: {
        id: id as string
      },
      include: {
        user: true
      }
    })

    if (location === null) {
      res
      .status(404)
      .json({ message: "No location found with this id" })
      return
    }

    res.status(200).json(location)
  }
  else if (req.method === "PATCH") {
    const { id, username, password, name, capacity } = req.body

    const location = await prisma.location.findFirst({
      where: {
        id: id
      }
    })

    if (location === null) {
      res
      .status(404)
      .json({ message: "No location found with this id" })
      return
    }

    await prisma.user.update({
      where: {
        id: location.userId
      },
      data: {
        username: username as string | undefined,
        passwordHash: hashSync(password, 12)
      }
    })

    const result = await prisma.location.update({
      where: {
        id: id
      },
      data: {
        name: name as string | undefined,
        capacity: capacity === undefined ? undefined : parseInt(capacity)
      }
    })

    res.status(200).json(result)
  }
  else if (req.method === "DELETE") {
    const { id } = req.body

    const location = await prisma.location.findFirst({
      where: {
        id: id
      }
    })

    if (location === null) {
      res
      .status(404)
      .json({ message: "No location found with this id" })
      return
    }

    await prisma.location.delete({
      where: {
        id: location.id
      }
    })

    await prisma.user.delete({
      where: {
        id: location.userId
      }
    })

    res.status(200).json(location)
  }
  else {
    res.setHeader("Allow", ["POST"])
    res
      .status(405)
      .json({ message: `HTTP method ${req.method} is not supported.` })
  }
}