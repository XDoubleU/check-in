import { format } from "date-fns"
import { getDataForDayChart } from "api-wrapper"
import { useEffect, useState, type Dispatch, type SetStateAction } from "react"
import { Col, Row } from "react-bootstrap"
import { Area } from "recharts"
import FormInput from "components/forms/FormInput"
import { convertDates, extractAllSchools } from "./dataProcessing"
import { COLORS, DataLoading, NoDataFound, SharedComposedChart } from "./Shared"

interface DayChartProps extends FilterProps {
  locationId: string
  dayData: unknown[]
  setDayData: Dispatch<SetStateAction<unknown[]>>
}

interface FilterProps {
  date: string
  setDate: Dispatch<SetStateAction<string>>
}

function Filter({ date, setDate }: FilterProps) {
  return (
    <Row>
      <Col>
        <FormInput
          label="Date"
          type="date"
          value={date}
          onChange={(e) => setDate(e.target.value)}
        />
      </Col>
      <Col></Col>
    </Row>
  )
}

// eslint-disable-next-line max-lines-per-function
export default function DayChart({
  locationId,
  dayData,
  date,
  setDayData,
  setDate
}: DayChartProps) {
  const [schools, setSchools] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void getDataForDayChart(locationId, date).then((response) => {
      let data = response.data ?? []
      data = convertDates(data)
      setDayData(response.data ?? [])
      setSchools(extractAllSchools(data))
      setLoading(false)
    })
  }, [date, locationId, setDayData])

  if (loading) {
    return (
      <>
        <Filter date={date} setDate={setDate} />
        <DataLoading />
      </>
    )
  }

  if (dayData.length === 0) {
    return (
      <>
        <Filter date={date} setDate={setDate} />
        <NoDataFound />
      </>
    )
  }

  return (
    <>
      <Filter date={date} setDate={setDate} />
      <SharedComposedChart
        data={dayData}
        xAxisTickFomatter={(datetime: Date) => format(datetime, "HH:mm")}
      >
        {schools.map((school, index) => {
          return (
            <Area
              key={school}
              stackId="a"
              dataKey={school}
              stroke={COLORS[index % 32]}
              fill={COLORS[index % 32]}
            />
          )
        })}
      </SharedComposedChart>
    </>
  )
}
