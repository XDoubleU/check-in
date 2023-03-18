import { type CSSProperties, type MouseEventHandler } from "react"
import { Button } from "react-bootstrap"

interface CustomButtonProps {
  children: string
  type?: "button" | "submit" | "reset"
  onClick?: MouseEventHandler<HTMLButtonElement>
  style?: CSSProperties
  className?: string
  value?: string | number
}

export default function CustomButton({
  children,
  type,
  onClick,
  style,
  className,
  value
}: CustomButtonProps) {
  return (
    <Button
      className={`${className ?? ""} text-white`}
      type={type}
      onClick={onClick}
      style={style}
      value={value}
    >
      {children}
    </Button>
  )
}
