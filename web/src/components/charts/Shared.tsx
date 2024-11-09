import { type ReactNode } from "react"
import {
  CartesianGrid,
  ComposedChart,
  Legend,
  Line,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from "recharts"
import Loader from "components/Loader"

export interface ChartDataEntry {
  [name: string]: number | string
  datetime: string
  capacity: number
}
export type ChartData = ChartDataEntry[]

const RESPONSIVE_CONTAINER_PROPS = {
  width: "100%",
  height: 500,
  aspect: 3
}

const CHART_PROPS = {
  width: 500,
  height: 300,
  margin: {
    top: 20,
    right: 30,
    left: 20,
    bottom: 5
  }
}

// colors used in google charts
export const COLORS = [
  "#3366cc",
  "#dc3912",
  "#ff9900",
  "#109618",
  "#990099",
  "#0099c6",
  "#dd4477",
  "#66aa00",
  "#b82e2e",
  "#316395",
  "#994499",
  "#22aa99",
  "#aaaa11",
  "#6633cc",
  "#e67300",
  "#8b0707",
  "#651067",
  "#329262",
  "#5574a6",
  "#3b3eac",
  "#b77322",
  "#16d620",
  "#b91383",
  "#f4359e",
  "#9c5935",
  "#a9c413",
  "#2a778d",
  "#668d1c",
  "#bea413",
  "#0c5922",
  "#743411"
]

export function NoDataFound() {
  return (
    <ResponsiveContainer {...RESPONSIVE_CONTAINER_PROPS}>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          height: "100%"
        }}
      >
        <h2>No data found</h2>
      </div>
    </ResponsiveContainer>
  )
}

export function DataLoading() {
  return (
    <ResponsiveContainer {...RESPONSIVE_CONTAINER_PROPS}>
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          height: "100%"
        }}
      >
        <Loader message="Fetching chart data." />
      </div>
    </ResponsiveContainer>
  )
}

interface SharedComposedChartProps {
  data: ChartData
  xAxisTickFomatter: (datetime: string) => string
  children: ReactNode
}

export function SharedComposedChart({
  data,
  xAxisTickFomatter,
  children
}: Readonly<SharedComposedChartProps>) {
  return (
    <ResponsiveContainer {...RESPONSIVE_CONTAINER_PROPS}>
      <ComposedChart data={data} {...CHART_PROPS}>
        <CartesianGrid strokeDasharray="3 3" />
        <Tooltip labelFormatter={xAxisTickFomatter} />
        <Legend />
        <YAxis />
        <XAxis dataKey="datetime" tickFormatter={xAxisTickFomatter} />
        <Line dataKey="capacity" stroke="red" strokeDasharray="3 3" />
        {children}
      </ComposedChart>
    </ResponsiveContainer>
  )
}
