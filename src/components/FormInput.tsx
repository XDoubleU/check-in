import { ChangeEventHandler } from "react"

type FormItemProps = {
  name: string,
  type: string,
  value: string | number,
  onChange: ChangeEventHandler<HTMLInputElement>
}


export default function FormItem({name, type, value, onChange}: FormItemProps){
  return (
    <div className="mb-3">
      <label htmlFor={name} className="form-label">{name}</label>
      <input type={type} value={value} onChange={onChange} className="form-control" id={name} placeholder={name}/>
    </div>
  )
}