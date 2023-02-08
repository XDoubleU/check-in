import { useEffect, useState } from "react"
import { Card, Form } from "react-bootstrap"
import Link from "next/link"
import DeleteModal from "../modals/DeleteModal"
import UpdateModal from "../modals/UpdateModal"

type LocationCardProps = {
  id: string, 
  name: string,
  capacity: number,
  username: string
}

type LocationUpdateProps = {
  id: string, 
  name: string,
  capacity: number,
  username: string,
  password: string,
  repeatPassword: string
}

export function LocationUpdateModal({id, name, capacity, username}: LocationCardProps) {
  const [updateInfo, setUpdateInfo] = useState({
    id: id,
    name: name,
    capacity: capacity,
    username: username,
    password: "",
    repeatPassword: ""
  })
  const [updateFormError, setUpdateFormError] = useState("")

  useEffect(() => {
    if (updateInfo.password !== updateInfo.repeatPassword) {
      setUpdateFormError("Passwords don't match.")
    } else {
      setUpdateFormError("")
    }
  }, [updateInfo])

  return (
    <UpdateModal<LocationUpdateProps> updateInfo={updateInfo} endpoint={`/api/locations/${id}`} >
      <Form.Group className="mb-3">
        <Form.Label>Name</Form.Label>
        <Form.Control type="text" placeholder="Name" value={updateInfo.name} onChange={({ target}) => setUpdateInfo({ ...updateInfo, name: target.value })}></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Capacity</Form.Label>
        <Form.Control type="number" value={updateInfo.capacity} onChange={({ target}) => setUpdateInfo({ ...updateInfo, capacity: parseInt(target.value) })}></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Username</Form.Label>
        <Form.Control type="text" placeholder="Username" value={updateInfo.username} onChange={({ target}) => setUpdateInfo({ ...updateInfo, username: target.value })}></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Password</Form.Label>
        <Form.Control type="password" placeholder="Password" value={updateInfo.password} onChange={({ target}) => setUpdateInfo({ ...updateInfo, password: target.value })}></Form.Control>
      </Form.Group>
      <Form.Group className="mb-3">
        <Form.Label>Repeat password</Form.Label>
        <Form.Control type="password" placeholder="Repeat password" value={updateInfo.repeatPassword} onChange={({ target}) => setUpdateInfo({ ...updateInfo, repeatPassword: target.value })} ></Form.Control>
        <Form.Text className="text-danger">
          {updateFormError}
        </Form.Text>
      </Form.Group>
    </UpdateModal>
  )
}

export default function LocationCard({id, name, capacity, username}: LocationCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title><Link href={`/settings/locations/${id}`}>{name}</Link> (TODO: Normalized name)</Card.Title>
              <Card.Subtitle className="mb-2 text-muted">{capacity}</Card.Subtitle>
            </div>
            <div className="ms-auto">
              <LocationUpdateModal id={id} name={name} username={username} capacity={capacity} />
              <DeleteModal id={id} name={name} endpoint={`/api/locations/${id}`} />
            </div>
          </div>
        </Card.Body>
      </Card>
      <br/>
    </>
  )
}