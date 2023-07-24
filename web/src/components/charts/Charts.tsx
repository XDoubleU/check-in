import { useState } from "react"
import { Tab, Tabs } from "react-bootstrap"
import { startOfISOWeek, endOfISOWeek, format } from "date-fns"
import CustomButton from "components/CustomButton"
import RangeChart from "./RangeChart"
import DayChart from "./DayChart"
import { downloadCsvForDayChart, downloadCsvForRangeChart } from "api-wrapper"
import { DATE_FORMAT } from "api-wrapper/types/apiTypes"

interface ChartProps {
  locationId: string
}

// eslint-disable-next-line @typescript-eslint/naming-convention
function getDates(): string[] {
  const date = new Date()
  const weekStart = format(startOfISOWeek(date), DATE_FORMAT)
  const weekEnd = format(endOfISOWeek(date), DATE_FORMAT)

  return [weekStart, weekEnd]
}

// eslint-disable-next-line @typescript-eslint/naming-convention
function getDate(): string {
  return format(new Date(), DATE_FORMAT)
}

// eslint-disable-next-line max-lines-per-function
export default function Charts({ locationId }: ChartProps) {
  const [weekStart, weekEnd] = getDates()

  const [startDate, setStartDate] = useState(weekStart)
  const [endDate, setEndDate] = useState(weekEnd)
  const [date, setDate] = useState(getDate())

  const [rangeData, setRangeData] = useState<unknown[]>([])
  const [dayData, setDayData] = useState<unknown[]>([])

  return (
    <Tabs
      defaultActiveKey="range"
      unmountOnExit={true}
      mountOnEnter={true}
      className="mb-3"
      fill
    >
      <Tab eventKey="range" title="Range">
        <br />
        <CustomButton
          onClick={(event) => {
            event.preventDefault()
            downloadCsvForRangeChart(locationId, startDate, endDate)
          }}
        >
          Download as CSV
        </CustomButton>
        <br />
        <br />
        <RangeChart
          locationId={locationId}
          rangeData={rangeData}
          startDate={startDate}
          endDate={endDate}
          setRangeData={setRangeData}
          setStartDate={setStartDate}
          setEndDate={setEndDate}
        />
      </Tab>
      <Tab eventKey="day" title="Day">
        <br />
        <CustomButton
          onClick={(event) => {
            event.preventDefault()
            downloadCsvForDayChart(locationId, date)
          }}
        >
          Download as CSV
        </CustomButton>
        <br />
        <br />
        <DayChart
          locationId={locationId}
          dayData={dayData}
          date={date}
          setDayData={setDayData}
          setDate={setDate}
        />
      </Tab>
    </Tabs>
  )
}
