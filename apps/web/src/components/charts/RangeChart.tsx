import { getDataForRangeChart } from "api-wrapper"
import { type Dispatch, type SetStateAction, useEffect, useState } from "react"
import { Row, Col } from "react-bootstrap"
import { Bar } from "recharts"
import { COLORS, DataLoading, NoDataFound, SharedComposedChart } from "./Shared"
import { convertDates, extractAllSchools } from "./dataProcessing"
import { DATE_FORMAT } from "types-custom"
import { format } from "date-fns"
import FormInput from "components/forms/FormInput"

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
        <FormInput
          label="Start date"
          type="date"
          value={startDate}
          onChange={(e) => setStartDate(e.target.value)}
        />
      </Col>
      <Col>
        <FormInput
          label="End date"
          type="date"
          value={endDate}
          onChange={(e) => setEndDate(e.target.value)}
        />
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
      <SharedComposedChart
        data={rangeData}
        xAxisTickFomatter={(datetime: Date) => format(datetime, DATE_FORMAT)}
      >
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
      </SharedComposedChart>
    </>
  )
}
