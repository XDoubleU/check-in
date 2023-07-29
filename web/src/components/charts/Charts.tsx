import { useState } from "react"
import { Tab, Tabs } from "react-bootstrap"
import { startOfISOWeek, endOfISOWeek } from "date-fns"
import CustomButton from "components/CustomButton"
import RangeChart from "./RangeChart"
import DayChart from "./DayChart"
import { downloadCSVForDayChart, downloadCSVForRangeChart } from "api-wrapper"
import { type ChartData } from "./Shared"

interface ChartProps {
  locationId: string
}

// eslint-disable-next-line @typescript-eslint/naming-convention
function getDates(): Date[] {
  const date = new Date()
  const weekStart = startOfISOWeek(date)
  const weekEnd = endOfISOWeek(date)

  return [weekStart, weekEnd]
}

// eslint-disable-next-line @typescript-eslint/naming-convention
function getDate(): Date {
  return new Date()
}

// eslint-disable-next-line max-lines-per-function
export default function Charts({ locationId }: ChartProps) {
  const [weekStart, weekEnd] = getDates()

  const [startDate, setStartDate] = useState(weekStart)
  const [endDate, setEndDate] = useState(weekEnd)
  const [date, setDate] = useState(getDate())

  const [rangeData, setRangeData] = useState<ChartData>([])
  const [dayData, setDayData] = useState<ChartData>([])

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
            downloadCSVForRangeChart(locationId, startDate, endDate)
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
            downloadCSVForDayChart(locationId, date)
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
