import { BarLoader } from "react-spinners"

export default function Loader() {
  const style: React.CSSProperties = {
    position: "fixed",
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)"
  }

  return (
    <div style={style}>
      <BarLoader color="red" />
    </div>
  )
}
