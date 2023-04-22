/* eslint-disable no-prototype-builtins */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/naming-convention */
import { format } from "date-fns"
import { type CheckInEntity, type SchoolEntity } from "../entities"
import { DATE_FORMAT } from "types-custom"

export function convertRangeData(
  checkIns: CheckInEntity[],
  schools: SchoolEntity[]
): unknown[] {
  const result: any[] = []

  checkIns.forEach((checkIn) => {
    const datetime = format(new Date(checkIn.createdAt), DATE_FORMAT)
    const index = result.findIndex(
      (e) => format(new Date(e.datetime as number), DATE_FORMAT) === datetime
    )

    let entry: any = {
      datetime: new Date(datetime).getTime(),
      capacity: checkIn.capacity
    }
    if (index > -1) {
      entry = result[index]

      if (checkIn.capacity > entry.capacity) {
        entry.capacity = checkIn.capacity
      }
    }

    schools.forEach((school) => {
      if (entry.hasOwnProperty(school.name)) {
        entry[school.name] = 1
      } else {
        entry[school.name] = 0
      }

      if (checkIn.school.name === school.name) {
        entry[school.name]++
      }
    })

    if (index === -1) {
      result.push(entry)
    }
  })

  return result
}

export function convertDayData(
  checkIns: CheckInEntity[],
  schools: SchoolEntity[]
): unknown[] {
  const result: any[] = []

  checkIns.forEach((checkIn) => {
    const lastEntry = result[result.length - 1]

    const entry: any = {
      datetime: new Date(checkIn.createdAt).getTime(),
      capacity: checkIn.capacity
    }

    schools.forEach((school) => {
      if (lastEntry?.hasOwnProperty(school.name)) {
        entry[school.name] = lastEntry[school.name]
      } else {
        entry[school.name] = 0
      }

      if (checkIn.school.name === school.name) {
        entry[school.name]++
      }
    })

    if (lastEntry) {
      entry[checkIn.school.name] =
        (lastEntry[checkIn.school.name] as number) + 1
    } else {
      entry[checkIn.school.name] = 1
    }

    result.push(entry)
  })

  return result
}

export function convertDatetime(
  data: unknown[],
  formatString: string
): unknown[] {
  const result = []

  const anyData = data as any
  for (const entry of anyData) {
    entry.datetime = format(new Date(entry.datetime as number), formatString)
    result.push(entry)
  }

  return result
}
