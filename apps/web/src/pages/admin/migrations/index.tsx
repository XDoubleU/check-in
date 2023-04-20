import {
  applyMigrationsDown,
  applyMigrationsUp,
  applySeeder
} from "api-wrapper"
import { Col, Row } from "react-bootstrap"
import CustomButton from "components/CustomButton"
import AdminLayout from "layouts/AdminLayout"

// eslint-disable-next-line @typescript-eslint/naming-convention
async function migrationUp() {
  const name = (await applyMigrationsUp()).data

  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  document.getElementById("output-migrations-up")!.innerHTML = name ?? "Error"
}

// eslint-disable-next-line @typescript-eslint/naming-convention
async function migrationDown() {
  const name = (await applyMigrationsDown()).data

  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  document.getElementById("output-migrations-down")!.innerHTML = name ?? "Error"
}

// eslint-disable-next-line @typescript-eslint/naming-convention
async function seed() {
  const ok = (await applySeeder()).ok

  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  document.getElementById("output-seed")!.innerHTML = ok ? "ACK" : "NACK"
}

export default function Migrations() {
  return (
    <AdminLayout title="Migrations">
      <Row>
        <Col md={2}>
          <CustomButton onClick={migrationUp}>Migrations up</CustomButton>
          <br />
          <br />
          <CustomButton onClick={migrationDown}>Migrations down</CustomButton>
          <br />
          <br />
          <CustomButton onClick={seed}>Seed</CustomButton>
        </Col>
        <Col md={2}>
          <p id="output-migrations-up">...</p>
          <br />
          <p id="output-migrations-down">...</p>
          <br />
          <p id="output-seed">...</p>
        </Col>
      </Row>
    </AdminLayout>
  )
}
