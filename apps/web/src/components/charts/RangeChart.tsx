import { getDataForRangeChart } from "my-api-wrapper"
import { type Dispatch, type SetStateAction, useEffect, useState } from "react"
import { Row, Col, Form } from "react-bootstrap"
import {
  ResponsiveContainer,
  CartesianGrid,
  Tooltip,
  XAxis,
  YAxis,
  Legend,
  Bar,
  ComposedChart,
  Line
} from "recharts"
import {
  CHART_PROPS,
  COLORS,
  DataLoading,
  NoDataFound,
  RESPONSIVE_CONTAINER_PROPS
} from "./Shared"
import { convertDates, extractAllSchools } from "./dataProcessing"
import { DATE_FORMAT } from "types-custom"
import { format } from "date-fns"

interface RangeChartProps extends FilterProps {
  locationId: string
  rangeData: unknown[]
  setRangeData: Dispatch<SetStateAction<unknown[]>>
}

interface FilterProps {
  startDate: string
  endDate: string
  setStartDate: Dispatch<SetStateAction<string>>
  setEndDate: Dispatch<SetStateAction<string>>
}

function Filter({ startDate, endDate, setStartDate, setEndDate }: FilterProps) {
  return (
    <Row>
      <Col>
        <Form.Group className="mb-3">
          <Form.Label>Start date</Form.Label>
          <Form.Control
            type="date"
            value={startDate}
            onChange={(e) => setStartDate(e.target.value)}
          />
        </Form.Group>
      </Col>
      <Col>
        <Form.Group className="mb-3">
          <Form.Label>End date</Form.Label>
          <Form.Control
            type="date"
            value={endDate}
            onChange={(e) => setEndDate(e.target.value)}
          />
        </Form.Group>
      </Col>
    </Row>
  )
}

// eslint-disable-next-line max-lines-per-function
export default function RangeChart({
  locationId,
  rangeData,
  startDate,
  endDate,
  setRangeData,
  setStartDate,
  setEndDate
}: RangeChartProps) {
  const [schools, setSchools] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void getDataForRangeChart(locationId, startDate, endDate).then(
      (response) => {
        let data = response.data ?? []
        data = convertDates(data)
        setRangeData(data)
        setSchools(extractAllSchools(data))
        setLoading(false)
      }
    )
  }, [startDate, endDate, setRangeData, locationId])

  if (loading) {
    return (
      <>
        <Filter
          startDate={startDate}
          endDate={endDate}
          setStartDate={setStartDate}
          setEndDate={setEndDate}
        />
        <DataLoading />
      </>
    )
  }

  // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-explicit-any
  if (rangeData.length === 0) {
    return (
      <>
        <Filter
          startDate={startDate}
          endDate={endDate}
          setStartDate={setStartDate}
          setEndDate={setEndDate}
        />
        <NoDataFound />
      </>
    )
  }

  return (
    <>
      <Filter
        startDate={startDate}
        endDate={endDate}
        setStartDate={setStartDate}
        setEndDate={setEndDate}
      />
      <ResponsiveContainer {...RESPONSIVE_CONTAINER_PROPS}>
        <ComposedChart data={rangeData as never} {...CHART_PROPS}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis
            dataKey="datetime"
            tickFormatter={(datetime: Date) => format(datetime, DATE_FORMAT)}
          />
          <YAxis />
          <Tooltip />
          <Legend />
          <Line dataKey="capacity" stroke="red" strokeDasharray="3 3" />
          {schools.map((school, index) => {
            return (
              <Bar
                key={school}
                dataKey={school}
                stackId="a"
                fill={COLORS[index % 32]}
              />
            )
          })}
        </ComposedChart>
      </ResponsiveContainer>
    </>
  )
}
