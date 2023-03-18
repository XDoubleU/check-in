import CustomButton from "@/components/CustomButton"
import CustomPagination, {
  type CustomPaginationProps
} from "@/components/CustomPagination"
import SchoolCard from "@/components/cards/SchoolCard"
import AdminLayout from "@/layouts/AdminLayout"
import { type School } from "types-custom"
import { createSchool, getAllSchools } from "my-api-wrapper"
import { useRouter } from "next/router"
import { type FormEvent, useCallback, useEffect, useState } from "react"
import { Col, Form, Modal } from "react-bootstrap"

interface SchoolList {
  schools: School[]
  pagination: CustomPaginationProps
}

// eslint-disable-next-line max-lines-per-function
export default function SchoolList() {
  const router = useRouter()

  const [schoolList, setSchoolList] = useState<SchoolList>({
    schools: [],
    pagination: {
      current: 0,
      total: 0
    }
  })
  const [createInfo, setCreateInfo] = useState({ name: "" })
  const [showCreate, setShowCreate] = useState(false)
  const handleCloseCreate = () => setShowCreate(false)
  const handleShowCreate = () => setShowCreate(true)
  const onCloseCreate = useCallback(() => {
    return !showCreate
  }, [showCreate])

  useEffect(() => {
    if (!router.isReady) return
    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined
    void getAllSchools(page).then(async (data) => {
      if (!data) {
        await router.push("/signin")
        return
      }

      setSchoolList({
        schools: data.schools,
        pagination: {
          current: data.page,
          total: data.totalPages
        }
      })
    })
  }, [onCloseCreate, router])

  const handleCreate = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const response = await createSchool(createInfo.name)

    if (response) {
      createInfo.name = ""
      handleCloseCreate()
    } else {
      console.log("ERROR")
    }
  }

  return (
    <AdminLayout title="Schools">
      <Modal show={showCreate} onHide={handleCloseCreate}>
        <Modal.Body>
          <Modal.Title>Create school</Modal.Title>
          <br />
          <Form onSubmit={() => handleCreate}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control
                type="text"
                placeholder="Name"
                value={createInfo.name}
                onChange={({ target }) =>
                  setCreateInfo({ ...createInfo, name: target.value })
                }
              ></Form.Control>
            </Form.Group>
            <br />
            <CustomButton
              type="button"
              style={{ float: "left" }}
              onClick={handleCloseCreate}
            >
              Cancel
            </CustomButton>
            <CustomButton type="submit" style={{ float: "right" }}>
              Create
            </CustomButton>
          </Form>
        </Modal.Body>
      </Modal>

      <Col size={2}>
        <CustomButton onClick={handleShowCreate}>Create</CustomButton>
      </Col>

      <br />

      <div className="min-vh-51">
        {schoolList.schools.length == 0 ? "Nothing to see here." : ""}

        {schoolList.schools.map((school) => {
          return (
            <SchoolCard id={school.id} key={school.id} name={school.name} />
          )
        })}
      </div>

      <CustomPagination
        current={schoolList.pagination.current}
        total={schoolList.pagination.total}
      />
    </AdminLayout>
  )
}
