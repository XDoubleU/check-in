import { type CheckInsGraphDto } from "api-wrapper/types/apiTypes"
import { type ChartDataEntry, type ChartData } from "./Shared"

export function extractAllSchools(entries: CheckInsGraphDto): string[] {
  return Object.keys(entries.valuesPerSchool)
}

export function convertToChartData(entries: CheckInsGraphDto): ChartData {
  let result: ChartData = []

  for (let i = 0; i < entries.dates.length; i++) {
    const entry: ChartDataEntry = {
      datetime: entries.dates[i],
      capacity: Object.values(entries.capacitiesPerLocation).reduce(
        (acc, val) => acc + val[i],
        0
      )
    }

    for (const [schoolKey, schoolValues] of Object.entries(
      entries.valuesPerSchool
    )) {
      entry[schoolKey] = schoolValues[i]
    }

    result.push(entry)
  }

  result = result.sort(function (a: ChartDataEntry, b: ChartDataEntry) {
    return new Date(a.datetime).getTime() - new Date(b.datetime).getTime()
  })

  return result
}
