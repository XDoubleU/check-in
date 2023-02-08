import { Prisma } from "@prisma/client"

export type LocationWithUser = Prisma.LocationGetPayload<{
  include: {
    user: true
  }
}>