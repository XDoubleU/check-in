import { type CheckInsLocationEntryRawMap } from "api-wrapper/types/apiTypes"
import { type ChartDataEntry, type ChartData } from "./Shared"

export function extractAllSchools(
  entries: CheckInsLocationEntryRawMap
): string[] {
  const key = Object.keys(entries)[0]
  return Object.keys(entries[key].schools)
}

export function convertToChartData(
  entries: CheckInsLocationEntryRawMap
): ChartData {
  let result: ChartData = []

  for (const [key, value] of Object.entries(entries)) {
    const entry: ChartDataEntry = {
      datetime: key,
      capacity: value.capacity
    }

    for (const [schoolKey, schoolValue] of Object.entries(value.schools)) {
      entry[schoolKey] = schoolValue
    }

    result.push(entry)
  }

  result = result.sort(function (a: ChartDataEntry, b: ChartDataEntry) {
    return new Date(a.datetime).getTime() - new Date(b.datetime).getTime()
  })

  return result
}
