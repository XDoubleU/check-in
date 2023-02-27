import CustomButton from "@/components/CustomButton"
import CustomPagination, { CustomPaginationProps } from "@/components/CustomPagination"
import LocationCard from "@/components/cards/LocationCard"
import AdminLayout from "@/layouts/AdminLayout"
import { Location } from "types"
import { useRouter } from "next/router"
import { FormEvent, useCallback, useEffect, useState } from "react"
import { Col, Form, Modal } from "react-bootstrap"
import { createLocation, getAllLocations } from "api-wrapper"

interface LocationList {
  locations: Location[],
  pagination: CustomPaginationProps
}

export default function LocationList() {
  const router = useRouter()

  const [locationList, setLocationList] = useState<LocationList>({
    locations: [],
    pagination: {
      current: 0,
      total: 0
    }
  })
  const [createInfo, setCreateInfo] = useState({
    name: "",
    capacity: 20,
    username: "",
    password: "",
    repeatPassword: ""
  })
  const [showCreate, setShowCreate] = useState(false)
  const handleCloseCreate = () => setShowCreate(false)
  const handleShowCreate = () => setShowCreate(true)
  const onCloseCreate = useCallback(() => {
    return !showCreate 
  }, [showCreate])

  useEffect(() => {
    if(!router.isReady) return
    const page = router.query.page ? parseInt(router.query.page as string) : undefined
    void getAllLocations(page)
      .then(async (data) => {
        if (!data) {
          await router.push("/signin")
          return
        }

        setLocationList({
          locations: data.locations,
          pagination: {
            current: data.page,
            total: data.totalPages
          }
        })
      })
  }, [onCloseCreate, router])

  const handleCreate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const response = await createLocation(
      createInfo.name,
      createInfo.capacity,
      createInfo.username,
      createInfo.password
    )

    if (response) {
      createInfo.name = ""
      createInfo.capacity = 20
      createInfo.username = ""
      createInfo.password = ""
      createInfo.repeatPassword = ""
      handleCloseCreate()
    } else {
      console.log("ERROR")
    }
  }

  return (
    <AdminLayout title="Locations">
      <Modal show={showCreate} onHide={handleCloseCreate}>
        <Modal.Body>
          <Modal.Title>Create location</Modal.Title>
          <br/>
          <Form onSubmit={() => handleCreate}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control type="text" placeholder="Name" value={createInfo.name} onChange={({ target}) => setCreateInfo({ ...createInfo, name: target.value })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Capacity</Form.Label>
              <Form.Control type="number" value={createInfo.capacity} onChange={({ target}) => setCreateInfo({ ...createInfo, capacity: parseInt(target.value) })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Username</Form.Label>
              <Form.Control type="text" placeholder="Username" value={createInfo.username} onChange={({ target}) => setCreateInfo({ ...createInfo, username: target.value })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Password</Form.Label>
              <Form.Control type="password" placeholder="Password" value={createInfo.password} onChange={({ target}) => setCreateInfo({ ...createInfo, password: target.value })}></Form.Control>
            </Form.Group>
            <Form.Group className="mb-3">
              <Form.Label>Repeat password</Form.Label>
              <Form.Control type="password" placeholder="Repeat password" value={createInfo.repeatPassword} onChange={({ target}) => setCreateInfo({ ...createInfo, repeatPassword: target.value })} ></Form.Control>
              <Form.Text className="text-danger">
                TODO: error
              </Form.Text>
            </Form.Group>
            <br/>
            <CustomButton type="button" style={{"float": "left"}} onClick={handleCloseCreate}>Cancel</CustomButton>
            <CustomButton type="submit" style={{"float": "right"}}>Create</CustomButton>
          </Form>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowCreate}>
          Create
        </CustomButton>
      </Col>

      <br/>

      <div className="min-vh-51">
        {
          (locationList.locations.length == 0) ? "Nothing to see here." : ""
        }

        {
          locationList.locations.map((location) => {
            return <LocationCard
              id={location.id}
              key={location.id}
              name={location.name}
              normalizedName={location.normalizedName}
              capacity={location.capacity}
              username={location.user.username} />
          })
        }
      </div>

      <CustomPagination
        current={locationList.pagination.current} 
        total={locationList.pagination.total} 
      />
      
    </AdminLayout>
  )  
}