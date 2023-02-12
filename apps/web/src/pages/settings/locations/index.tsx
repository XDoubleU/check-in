import CustomButton from "@/components/CustomButton"
import CustomPagination, { CustomPaginationProps } from "@/components/CustomPagination"
import LocationCard from "@/components/cards/LocationCard"
import AdminLayout from "@/layouts/AdminLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { Location } from "types"
import { GetServerSidePropsContext } from "next"
import { useSession } from "next-auth/react"
import Router, { useRouter } from "next/router"
import { FormEvent, useEffect, useState } from "react"
import { Col, Form, Modal } from "react-bootstrap"
import { getAllLocations } from "api-wrapper"

type LocationListProps = {
  locations: Location[],
  pagination: CustomPaginationProps
}

export default function LocationList({locations, pagination}: LocationListProps) {
  const router = useRouter()

  const [addInfo, setAddInfo] = useState({
    name: "",
    capacity: 20,
    username: "",
    password: "",
    repeatPassword: ""
  })
  const [addFormError, setAddFormError] = useState("")
  const [showAdd, setShowAdd] = useState(false)
  const handleCloseAdd = () => setShowAdd(false)
  const handleShowAdd = () => setShowAdd(true)

  const {data, status} = useSession({
    required: true
  })

  useEffect(() => {
    if (addInfo.password !== addInfo.repeatPassword) {
      setAddFormError("Passwords don't match.")
    } else {
      setAddFormError("")
    }
  }, [addInfo])

  const handleAdd = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      username: addInfo.username,
      password: addInfo.password,
      name: addInfo.name,
      capacity: addInfo.capacity
    }

    const response = await fetch("/api/locations", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data)
    })

    if (response.status < 300) {
      router.replace(router.asPath)
      addInfo.name = ""
      addInfo.capacity = 20
      addInfo.username = ""
      addInfo.password = ""
      addInfo.repeatPassword = ""
      handleCloseAdd()
    }
  }

  if (status == "loading") {
    return <LoadingLayout/>
  }

  if (!data.user.isAdmin) {
    Router.push(`/settings/locations/${data.user.locationId}`)
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title="Locations" user={data.user}>
      <Modal show={showAdd} onHide={handleCloseAdd}>
        <Modal.Body>
          <Modal.Title>Add location</Modal.Title>
          <br/>
          <Form onSubmit={handleAdd}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control type="text" placeholder="Name" value={addInfo.name} onChange={({ target}) => setAddInfo({ ...addInfo, name: target.value })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Capacity</Form.Label>
              <Form.Control type="number" value={addInfo.capacity} onChange={({ target}) => setAddInfo({ ...addInfo, capacity: parseInt(target.value) })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Username</Form.Label>
              <Form.Control type="text" placeholder="Username" value={addInfo.username} onChange={({ target}) => setAddInfo({ ...addInfo, username: target.value })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Password</Form.Label>
              <Form.Control type="password" placeholder="Password" value={addInfo.password} onChange={({ target}) => setAddInfo({ ...addInfo, password: target.value })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Repeat password</Form.Label>
              <Form.Control type="password" placeholder="Repeat password" value={addInfo.repeatPassword} onChange={({ target}) => setAddInfo({ ...addInfo, repeatPassword: target.value })} ></Form.Control>
              <Form.Text className="text-danger">
                {addFormError}
              </Form.Text>
            </Form.Group>
            <br/>
            <CustomButton type="button" style={{"float": "left"}} onClick={handleCloseAdd}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Add</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowAdd}>
          Add
        </CustomButton>
      </Col>

      <br/>

      <div className="min-vh-51">
        {
          (locations === undefined || locations.length == 0) ? "Nothing to see here." : ""
        }

        {
          locations.map((location) => {
            return <LocationCard
              id={location.id}
              key={location.id}
              name={location.name}
              capacity={location.capacity}
              username={location.user.username} />
          })
        }
      </div>

      <CustomPagination current={pagination.current} total={pagination.total} pageSize={pagination.pageSize} />
      
    </AdminLayout>
  )  
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const paginatedLocations = await getAllLocations(parseInt(context.query.page as string))

  return {
    props: {
      locations: paginatedLocations.locations,
      pagination: {
        total: paginatedLocations.totalPages,
        current: paginatedLocations.page
      }
    }
  }
}