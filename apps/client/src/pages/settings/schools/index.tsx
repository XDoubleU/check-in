import CustomButton from "@/components/CustomButton"
import CustomPagination, { CustomPaginationProps } from "@/components/CustomPagination"
import SchoolCard from "@/components/cards/SchoolCard"
import AdminLayout from "@/layouts/AdminLayout"
import LoadingLayout from "@/layouts/LoadingLayout"
import { School } from "@prisma/client"
import { GetServerSidePropsContext } from "next"
import { useSession } from "next-auth/react"
import { useRouter } from "next/router"
import { FormEvent, useState } from "react"
import { Col, Form, Modal } from "react-bootstrap"

type SchoolListProps = {
  schools: School[],
  pagination: CustomPaginationProps
}

export default function SchoolList({schools, pagination}: SchoolListProps) {
  const router = useRouter()

  const [addInfo, setAddInfo] = useState({ name: "" })
  const [showAdd, setShowAdd] = useState(false)
  const handleCloseAdd = () => setShowAdd(false)
  const handleShowAdd = () => setShowAdd(true)

  const {data, status} = useSession({
    required: true
  })

  const handleAdd = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const data = {
      ...addInfo
    }

    const response = await fetch("/api/schools", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data)
    })

    if (response.status < 300) {
      router.replace(router.asPath)
      addInfo.name = ""
      handleCloseAdd()
    }
  }

  if (status == "loading") {
    return <LoadingLayout/>
  }

  return (
    <AdminLayout title="Schools" user={data.user}>
      <Modal show={showAdd} onHide={handleCloseAdd}>
        <Modal.Body>
          <Modal.Title>Add school</Modal.Title>
          <br/>
          <Form onSubmit={handleAdd}>
            <Form.Group className="mb-3">
              <Form.Label>Name</Form.Label>
              <Form.Control type="text" placeholder="Name" value={addInfo.name} onChange={({ target}) => setAddInfo({ ...addInfo, name: target.value })}></Form.Control>
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
          (schools === undefined || schools.length == 0) ? "Nothing to see here." : ""
        }

        {
          schools.map((school) => {
            return <SchoolCard id={school.id} key={school.id} name={school.name} />
          })
        }
      </div>

      <CustomPagination current={pagination.current} total={pagination.total} pageSize={pagination.pageSize} />
      
    </AdminLayout>
  )  
}

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const currentPage = parseInt(context.query.page as string ?? "1")
  const response = await fetch(`${process.env.API_URL}/schools?page=${currentPage}`)
  const jsonResponse = await response.json()

  return {
    props: {
      schools: jsonResponse.schools,
      pagination: {
        total: jsonResponse.totalPages,
        current: jsonResponse.page
      }
    }
  }
}