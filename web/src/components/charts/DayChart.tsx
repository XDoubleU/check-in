import { getDataForDayChart } from "api-wrapper"
import { useEffect, useState, type Dispatch, type SetStateAction } from "react"
import { Col, Row } from "react-bootstrap"
import { Area } from "recharts"
import FormInput from "components/forms/FormInput"
import { convertToChartData, extractAllSchools } from "./dataProcessing"
import {
  COLORS,
  type ChartData,
  DataLoading,
  NoDataFound,
  SharedComposedChart
} from "./Shared"
import { DATE_FORMAT, TIME_FORMAT } from "api-wrapper/types/apiTypes"
import moment, { type Moment } from "moment"

interface DayChartProps extends FilterProps {
  locationIds: string[]
  dayData: ChartData
  setDayData: Dispatch<SetStateAction<ChartData>>
}

interface FilterProps {
  date: Moment
  setDate: Dispatch<SetStateAction<Moment>>
}

function Filter({ date, setDate }: FilterProps) {
  return (
    <Row>
      <Col>
        <FormInput
          label="Date"
          type="date"
          value={date.format(DATE_FORMAT)}
          onChange={(e) => setDate(moment(e.target.value))}
        />
      </Col>
      <Col></Col>
    </Row>
  )
}

// eslint-disable-next-line max-lines-per-function
export default function DayChart({
  locationIds,
  dayData,
  date,
  setDayData,
  setDate
}: DayChartProps) {
  const [schools, setSchools] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    void getDataForDayChart(locationIds, date)
      .then((response) => {
        if (
          !response.ok ||
          !response.data ||
          Object.keys(response.data).length === 0
        ) {
          setDayData([])
          return
        }

        const newData = convertToChartData(response.data)
        setDayData(newData)
        setSchools(extractAllSchools(response.data))
      })
      .then(() => setLoading(false))
  }, [date, locationIds, setDayData])

  if (loading) {
    return (
      <>
        <Filter date={date} setDate={setDate} />
        <DataLoading />
      </>
    )
  }

  if (Object.keys(dayData).length === 0) {
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
        xAxisTickFomatter={(datetime: string) =>
          moment.utc(datetime).format(TIME_FORMAT)
        }
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
