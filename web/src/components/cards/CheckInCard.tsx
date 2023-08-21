import { Card } from "react-bootstrap"
import DeleteModal from "components/modals/DeleteModal"
import { deleteCheckIn } from "api-wrapper"
import { type ICardProps } from "interfaces/ICardProps"
import { type CheckIn, type User } from "api-wrapper/types/apiTypes"

type CheckInCardProps = ICardProps<CheckIn> & {user: User}

function CheckInDeleteModal({ data, fetchData }: CheckInCardProps) {
  const handleDelete = () => {
    return deleteCheckIn(data.locationId, data.id)
  }

  return (
    <DeleteModal
      handler={handleDelete}
      fetchData={fetchData}
      typeName="checkIn"
    />
  )
}

export default function CheckInCard({ data, user, fetchData }: CheckInCardProps) {
  return (
    <>
      <Card>
        <Card.Body>
          <div className="d-flex flex-row">
            <div>
              <Card.Title>{`${data.createdAt} - ${data.schoolName} (ID: ${data.id})`}</Card.Title>
            </div>
            {user.role === "admin" || user.role === "manager" && (
              <div className="ms-auto">
                <CheckInDeleteModal data={data} user={user} fetchData={fetchData} />
              </div>
            )}
          </div>
        </Card.Body>
      </Card>
      <br />
    </>
  )
}
