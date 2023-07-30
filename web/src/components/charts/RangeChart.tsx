import { getDataForRangeChart } from "api-wrapper"
import { type Dispatch, type SetStateAction, useEffect, useState } from "react"
import { Row, Col } from "react-bootstrap"
import { Bar } from "recharts"
import {
  COLORS,
  type ChartData,
  DataLoading,
  NoDataFound,
  SharedComposedChart
} from "./Shared"
import { convertToChartData, extractAllSchools } from "./dataProcessing"
import { format } from "date-fns"
import FormInput from "components/forms/FormInput"
import { DATE_FORMAT } from "api-wrapper/types/apiTypes"

interface RangeChartProps extends FilterProps {
  locationId: string
  rangeData: ChartData
  setRangeData: Dispatch<SetStateAction<ChartData>>
}

interface FilterProps {
  startDate: Date
  endDate: Date
  setStartDate: Dispatch<SetStateAction<Date>>
  setEndDate: Dispatch<SetStateAction<Date>>
}

function Filter({ startDate, endDate, setStartDate, setEndDate }: FilterProps) {
  return (
    <Row>
      <Col>
        <FormInput
          label="Start date"
          type="date"
          value={format(startDate, DATE_FORMAT)}
          onChange={(e) => setStartDate(new Date(e.target.value))}
          max={format(endDate, DATE_FORMAT)}
        />
      </Col>
      <Col>
        <FormInput
          label="End date"
          type="date"
          value={format(endDate, DATE_FORMAT)}
          onChange={(e) => setEndDate(new Date(e.target.value))}
          min={format(startDate, DATE_FORMAT)}
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
    void getDataForRangeChart(locationId, startDate, endDate)
      .then((response) => {
        if (
          !response.ok ||
          !response.data ||
          Object.keys(response.data).length === 0
        ) {
          setRangeData([])
          return
        }

        const newdata = convertToChartData(response.data)
        setRangeData(newdata)
        setSchools(extractAllSchools(response.data))
      })
      .then(() => setLoading(false))
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
        xAxisTickFomatter={(datetime: number) =>
          format(new Date(datetime), DATE_FORMAT)
        }
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
