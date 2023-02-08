import { prisma } from "@/common/prisma"
import { NextApiRequest, NextApiResponse } from "next"
import { hashSync } from "bcrypt"

//TODO: cleaner ?
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === "GET") {
    const { page, pageSize } = req.query
    const pageInt = page === undefined || Array.isArray(page) ? 1 : parseInt(page)
    const pageSizeInt = pageSize === undefined || Array.isArray(pageSize) ? 4 : parseInt(pageSize)

    const pageSkip = (pageInt - 1) >= 0 ? (pageInt - 1) : 0
    const locations = await prisma.location.findMany({
      include: {
        user: true
      }, 
      orderBy: {
        name: "asc"
      },
      take: pageSizeInt,
      skip: pageSkip * pageSizeInt
    })

    res.status(200).json(locations)
  }
  else if (req.method === "POST") {
    const { username, password, name, capacity } = req.body

    const existingUser = await prisma.user.findFirst({
      where: {
        username: username
      }
    })

    if (existingUser !== null) {
      res
        .status(409)
        .json({ message: `A user with the username ${username} already exists.` })
      return
    }

    const existingLocation = await prisma.location.findFirst({
      where: {
        name: name
      }
    })
    
    if (existingLocation !== null) {
      res
      .status(409)
      .json({ message: `A location with the name ${name} already exists.` })
      return
    }

    const passwordHash = hashSync(password, 12)
    const user = await prisma.user.create({
      data: {
        username: username,
        passwordHash: passwordHash,
        isAdmin: false
      }
    })

    const location = await prisma.location.create({
      data: {
        name: name,
        capacity: capacity,
        userId: user.id
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