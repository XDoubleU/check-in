import { type CheckInsLocationEntryRawMap } from "api-wrapper/types/apiTypes"
import { type ChartDataEntry, type ChartData } from "./Shared"

export function extractAllSchools(
  entries: CheckInsLocationEntryRawMap
): string[] {
  const key = parseInt(Object.keys(entries)[0])
  return Object.keys(entries[key].schools).sort()
}

export function convertToChartData(entries: CheckInsLocationEntryRawMap): ChartData {
  let result: ChartData = []

  for (const [key, value] of Object.entries(entries)) {
    const entry: ChartDataEntry = {
      datetime: parseInt(key),
      capacity: value.capacity
    }

    for (const [schoolKey, schoolValue] of Object.entries(value.schools)) {
      entry[schoolKey] = schoolValue
    }

    result.push(entry)
  }

  result = result.sort(function (a: ChartDataEntry, b: ChartDataEntry) {
    return a.datetime - b.datetime
  })

  return result
}
