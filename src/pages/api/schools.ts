import { prisma } from "@/common/prisma"
import { NextApiRequest, NextApiResponse } from "next"

//TODO: cleaner ?
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === "GET") {
    const { page, pageSize, all } = req.query

    if (Boolean(all) === true) {
      const schools = await prisma.school.findMany({
        orderBy: {
          name: "asc"
        }
      })
  
      res.status(200).json(schools)
    }

    const pageInt = page === undefined || Array.isArray(page) ? 1 : parseInt(page)
    const pageSizeInt = pageSize === undefined || Array.isArray(pageSize) ? 4 : parseInt(pageSize)

    const pageSkip = (pageInt - 1) >= 0 ? (pageInt - 1) : 0
    const schools = await prisma.school.findMany({
      where: {
        id: {
          not: 1
        }
      },
      orderBy: {
        name: "asc"
      },
      take: pageSizeInt,
      skip: pageSkip * pageSizeInt
    })

    res.status(200).json(schools)
  }
  else if (req.method === "POST") {
    const { name } = req.body

    const existingSchool = await prisma.school.findFirst({
      where: {
        name: name
      }
    })
    
    if (existingSchool !== null) {
      res
      .status(409)
      .json({ message: `A school with the name ${name} already exists.` })
      return
    }

    const school = await prisma.school.create({
      data: {
        name: name
      }
    })

    res.status(200).json(school)
  }
  else if (req.method === "PATCH") {
    const { id, name } = req.body

    const school = await prisma.school.findFirst({
      where: {
        id: id
      }
    })

    if (school === null) {
      res
      .status(404)
      .json({ message: "No school found with this id" })
      return
    }

    const result = await prisma.school.update({
      where: {
        id: id
      },
      data: {
        name: name
      }
    })

    res.status(200).json(result)
  }
  else if (req.method === "DELETE") {
    const { id } = req.body

    const school = await prisma.school.findFirst({
      where: {
        id: id
      }
    })

    if (school === null) {
      res
      .status(404)
      .json({ message: "No school found with this id" })
      return
    }

    await prisma.school.delete({
      where: {
        id: school.id
      }
    })

    res.status(200).json(school)
  }
  else {
    res.setHeader("Allow", ["POST"])
    res
      .status(405)
      .json({ message: `HTTP method ${req.method} is not supported.` })
  }
}