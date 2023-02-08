import { CSSProperties, MouseEventHandler } from "react"
import { Button } from "react-bootstrap"

type CustomButtonProps = {
  children: string,
  type?: "button" | "submit" | "reset" | undefined,
  onClick?: MouseEventHandler<HTMLButtonElement>,
  style?: CSSProperties | undefined
}

export default function CustomButton({children, type, onClick, style}: CustomButtonProps) {
  return (
    <Button className="text-white" type={type} onClick={onClick} style={style}>{children}</Button>
  )
}