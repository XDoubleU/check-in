/* eslint-disable @typescript-eslint/no-unsafe-return */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable @typescript-eslint/no-unsafe-member-access */
/* eslint-disable @typescript-eslint/no-unsafe-call */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/naming-convention */

export function extractAllSchools(data: unknown[]): string[] {
  const ignore = ["datetime", "capacity"]
  const result = new Array<string>()

  for (const entry of data) {
    const keys = Object.getOwnPropertyNames(entry)
    for (const key of keys) {
      if (!result.includes(key) && !ignore.includes(key)) {
        result.push(key)
      }
    }
  }

  return result
}

export function convertDates(data: unknown[]): unknown[] {
  const anyData = data as any

  return anyData
    .map((entry: any) => {
      entry.datetime = new Date(entry.datetime)
      return entry
    })
    .sort(function (a: any, b: any) {
      return a.datetime - b.datetime
    }) as unknown[]
}
