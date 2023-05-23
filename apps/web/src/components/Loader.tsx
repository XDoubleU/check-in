import { BarLoader } from "react-spinners"

export interface LoaderProps {
  message?: string
}

export default function Loader({ message }: LoaderProps) {
  const style: React.CSSProperties = {
    position: "fixed",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    textAlign: "center"
  }

  return (
    <div style={style}>
      <BarLoader color="red" />
      <p>{message}</p>
    </div>
  )
}
