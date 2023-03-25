import { format } from "date-fns"
import { getDataForDayChart } from "my-api-wrapper"
import { useEffect, useState, type Dispatch, type SetStateAction } from "react"
import { Col, Form, Row } from "react-bootstrap"
import {
  CartesianGrid,
  Legend,
  Line,
  ResponsiveContainer,
  XAxis,
  YAxis,
  Tooltip,
  Area,
  ComposedChart
} from "recharts"
import { convertDates, extractAllSchools } from "./dataProcessing"
import {
  CHART_PROPS,
  COLORS,
  NoDataFound,
  RESPONSIVE_CONTAINER_PROPS
} from "./Shared"

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
        <Form.Group className="mb-3">
          <Form.Label>Date</Form.Label>
          <Form.Control
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
          />
        </Form.Group>
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

  useEffect(() => {
    void getDataForDayChart(locationId, date).then((response) => {
      let data = response.data ?? []
      data = convertDates(data)
      setDayData(response.data ?? [])
      setSchools(extractAllSchools(data))
    })
  }, [date, locationId, setDayData])

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
      <ResponsiveContainer {...RESPONSIVE_CONTAINER_PROPS}>
        <ComposedChart data={dayData} {...CHART_PROPS}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="datetime"
            tickFormatter={(datetime: Date) => format(datetime, "HH:mm")}
          />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line dataKey="capacity" stroke="red" strokeDasharray="3 3" />
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
        </ComposedChart>
      </ResponsiveContainer>
    </>
  )
}
