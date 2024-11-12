import { useState } from "react"
import { Tab, Tabs } from "react-bootstrap"
import CustomButton from "components/CustomButton"
import RangeChart from "./RangeChart"
import DayChart from "./DayChart"
import { downloadCSVForDayChart, downloadCSVForRangeChart } from "api-wrapper"
import { type ChartData } from "./Shared"
import moment, { type Moment } from "moment"

interface ChartProps {
  locationIds: string[]
}

function getDates(): Moment[] {
  const date = new Date()
  const weekStart = moment(date).startOf("isoWeek")
  const weekEnd = moment(date).endOf("isoWeek")

  return [weekStart, weekEnd]
}

function getDate(): Moment {
  return moment(new Date())
}

// eslint-disable-next-line max-lines-per-function
export default function Charts({ locationIds }: ChartProps) {
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
            downloadCSVForRangeChart(locationIds, startDate, endDate)
          }}
        >
          Download as CSV
        </CustomButton>
        <br />
        <br />
        <RangeChart
          locationIds={locationIds}
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
            downloadCSVForDayChart(locationIds, date)
          }}
        >
          Download as CSV
        </CustomButton>
        <br />
        <br />
        <DayChart
          locationIds={locationIds}
          dayData={dayData}
          date={date}
          setDayData={setDayData}
          setDate={setDate}
        />
      </Tab>
    </Tabs>
  )
}
