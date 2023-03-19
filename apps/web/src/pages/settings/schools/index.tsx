import CustomButton from "@/components/CustomButton"
import CustomPagination, {
  type CustomPaginationProps
} from "@/components/CustomPagination"
import SchoolCard from "@/components/cards/SchoolCard"
import AdminLayout from "@/layouts/AdminLayout"
import { type CreateSchoolDto, type School } from "types-custom"
import { createSchool, getAllSchools } from "my-api-wrapper"
import { useRouter } from "next/router"
import { useCallback, useEffect, useState } from "react"
import { Alert, Col, Form, Modal } from "react-bootstrap"
import { type SubmitHandler, useForm } from "react-hook-form"

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

  const {
    register,
    handleSubmit,
    setError,
    reset,
    formState: { errors }
  } = useForm<CreateSchoolDto>()

  const [showCreate, setShowCreate] = useState(false)
  const handleCloseCreate = () => setShowCreate(false)
  const handleShowCreate = () => setShowCreate(true)
  const onCloseCreate = useCallback(() => {
    return !showCreate
  }, [showCreate])

  const fetchData = useCallback(async () => {
    if (!router.isReady) return

    const page = router.query.page
      ? parseInt(router.query.page as string)
      : undefined

    const response = await getAllSchools(page)
    if (!response.ok) {
      await router.push("/signin")
      return
    }

    setSchoolList({
      schools: response.data?.schools ?? Array<School>(),
      pagination: {
        current: response.data?.page ?? 1,
        total: response.data?.totalPages ?? 1
      }
    })
  }, [router])

  useEffect(() => {
    void fetchData()
  }, [fetchData, onCloseCreate])

  const onSubmit: SubmitHandler<CreateSchoolDto> = async (data) => {
    const response = await createSchool(data)

    if (!response.ok) {
      setError("root", {
        message: response.message ?? "Something went wrong"
      })
    } else {
      handleCloseCreate()
      reset()
    }
  }

  return (
    <AdminLayout title="Schools">
      <Modal show={showCreate} onHide={handleCloseCreate}>
        <Modal.Body>
          <Modal.Title>Create school</Modal.Title>
          <br />
          <Form onSubmit={handleSubmit(onSubmit)}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control
                type="text"
                placeholder="Name"
                required
                {...register("name")}
              ></Form.Control>
            </Form.Group>
            {errors.root && <Alert key="danger">{errors.root.message}</Alert>}
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
            <SchoolCard
              id={school.id}
              key={school.id}
              name={school.name}
              fetchData={fetchData}
            />
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
