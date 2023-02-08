import { prisma } from "@/common/prisma"
import { NextApiRequest, NextApiResponse } from "next"

//TODO: cleaner ?
export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  if (req.method === "GET") {
    const { locationId } = req.query

    const location = await prisma.location.findFirst({
      where: {
        id: locationId as string
      }
    })

    if (location === null) {
      res
      .status(404)
      .json({ message: "No location found with this id" })
      return
    }

    const today = new Date()
    today.setHours(0,0,0,0)

    const tomorrow = new Date(today)
    tomorrow.setDate(tomorrow.getDate() + 1)

    const count = await prisma.checkIn.count({
      where: {
        locationId: locationId as string,
        datetime: {
          gte: today,
          lt: tomorrow
        }
      }
    })

    res.status(200).json(count)
  }
  else if (req.method === "POST") {
    const { locationId, schoolId } = req.body

    const location = await prisma.location.findFirst({
      where: {
        id: locationId as string
      }
    })

    if (location === null) {
      res
      .status(404)
      .json({ message: "No location found with this id" })
      return
    }

    const school = await prisma.school.findFirst({
      where: {
        id: parseInt(schoolId)
      }
    })

    if (school === null) {
      res
      .status(404)
      .json({ message: "No school found with this id" })
      return
    }

    const checkIn = await prisma.checkIn.create({
      data: {
        locationId: location.id,
        capacity: location.capacity,
        schoolId: school.id
      }
    })

    res.status(200).json(checkIn) 
  }
  else {
    res.setHeader("Allow", ["POST"])
    res
      .status(405)
      .json({ message: `HTTP method ${req.method} is not supported.` })
  }
}