import { format } from "date-fns"
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
import { DATE_FORMAT } from "api-wrapper/types/apiTypes"

interface DayChartProps extends FilterProps {
  locationId: string
  dayData: ChartData
  setDayData: Dispatch<SetStateAction<ChartData>>
}

interface FilterProps {
  date: Date
  setDate: Dispatch<SetStateAction<Date>>
}

function Filter({ date, setDate }: FilterProps) {
  return (
    <Row>
      <Col>
        <FormInput
          label="Date"
          type="date"
          value={format(date, DATE_FORMAT)}
          onChange={(e) => setDate(new Date(e.target.value))}
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
    void getDataForDayChart(locationId, date)
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
  }, [date, locationId, setDayData])

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
        xAxisTickFomatter={(datetime: number) =>
          format(new Date(datetime), "HH:mm")
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
