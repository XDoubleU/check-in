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
import FormInput from "components/forms/FormInput"
import { DATE_FORMAT } from "api-wrapper/types/apiTypes"
import moment, { type Moment } from "moment"

interface RangeChartProps extends FilterProps {
  locationIds: string[]
  rangeData: ChartData
  setRangeData: Dispatch<SetStateAction<ChartData>>
}

interface FilterProps {
  startDate: Moment
  endDate: Moment
  setStartDate: Dispatch<SetStateAction<Moment>>
  setEndDate: Dispatch<SetStateAction<Moment>>
}

function Filter({
  startDate,
  endDate,
  setStartDate,
  setEndDate
}: Readonly<FilterProps>) {
  return (
    <Row>
      <Col>
        <FormInput
          label="Start date"
          type="date"
          value={startDate.format(DATE_FORMAT)}
          onChange={(e) => {
            setStartDate(moment(e.target.value))
          }}
          max={endDate.format(DATE_FORMAT)}
        />
      </Col>
      <Col>
        <FormInput
          label="End date"
          type="date"
          value={endDate.format(DATE_FORMAT)}
          onChange={(e) => {
            setEndDate(moment(e.target.value))
          }}
          min={startDate.format(DATE_FORMAT)}
        />
      </Col>
    </Row>
  )
}

// eslint-disable-next-line max-lines-per-function
export default function RangeChart({
  locationIds,
  rangeData,
  startDate,
  endDate,
  setRangeData,
  setStartDate,
  setEndDate
}: Readonly<RangeChartProps>) {
  const [schools, setSchools] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void getDataForRangeChart(locationIds, startDate, endDate)
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
      .then(() => {
        setLoading(false)
      })
  }, [startDate, endDate, setRangeData, locationIds])

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
        xAxisTickFomatter={(datetime: string) =>
          moment.utc(datetime).format(DATE_FORMAT)
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
