import {
  type GetAllPaginatedLocationDto,
  type CreateLocationDto,
  type Location
} from "types-custom"
import { Form } from "react-bootstrap"
import { createLocation, getAllLocations, getUser } from "api-wrapper"
import { useForm } from "react-hook-form"
import CreateModal from "components/modals/CreateModal"
import FormInput from "components/forms/FormInput"
import ListViewLayout, { type List } from "layouts/ListViewLayout"
import { useCallback, useState } from "react"
import LocationCard from "components/cards/LocationCard"
import { type ICreateModalProps } from "interfaces/ICreateModalProps"

type CreateLocationForm = CreateLocationDto & { repeatPassword?: string }

export type CreateLocationModalProps = ICreateModalProps<CreateLocationForm>

// eslint-disable-next-line max-lines-per-function
function CreateLocationModal({ form, fetchData }: CreateLocationModalProps) {
  const {
    register,
    watch,
    formState: { errors }
  } = form

  return (
    <CreateModal<CreateLocationDto, Location>
      form={form}
      handler={createLocation}
      fetchData={fetchData}
      typeName="location"
    >
      <FormInput
        label="Name"
        type="text"
        placeholder="Name"
        required
        register={register("name")}
      />
      <FormInput
        label="Capacity"
        type="number"
        placeholder={10}
        required
        register={register("capacity")}
      />
      <FormInput
        label="Username"
        type="text"
        placeholder="Username"
        required
        register={register("username")}
      />
      <FormInput
        label="Password"
        type="password"
        placeholder="Password"
        required
        autocomplete="new-password"
        register={register("password")}
      />
      {/* jscpd:ignore-start */}
      <Form.Group className="mb-3">
        <Form.Label>Repeat password</Form.Label>
        <Form.Control
          type="password"
          placeholder="Repeat password"
          autoComplete="new-password"
          required
          isInvalid={!!errors.repeatPassword}
          {...register("repeatPassword", {
            validate: (val: string | undefined) => {
              if (watch("password") != val) {
                return "Your passwords do no match"
              }
              return undefined
            }
          })}
        ></Form.Control>
        <Form.Control.Feedback type="invalid">
          {errors.repeatPassword?.message}
        </Form.Control.Feedback>
      </Form.Group>
      {/* jscpd:ignore-end */}
    </CreateModal>
  )
}

export type LocationWithUsername = Omit<Location, "userId"> & {
  username: string
}

type LocationList = List<LocationWithUsername>

// eslint-disable-next-line max-lines-per-function
export default function LocationListView() {
  const [locationList, setLocationList] = useState<LocationList>({
    data: undefined,
    pagination: {
      current: 0,
      total: 0
    }
  })

  const form = useForm<CreateLocationForm>()

  // If preprocessList doesn't use useCallback it will be called infinitely
  const preprocessList = useCallback(
    async (responseData: GetAllPaginatedLocationDto) => {
      const locationsWithUsernames: List<LocationWithUsername> = {
        data: [],
        pagination: {
          current: responseData.pagination.current,
          total: responseData.pagination.total
        }
      }

      if (!locationsWithUsernames.data) {
        return locationsWithUsernames
      }

      for (const location of responseData.data) {
        const username = (await getUser(location.userId)).data?.username

        locationsWithUsernames.data.push({
          id: location.id,
          name: location.name,
          normalizedName: location.normalizedName,
          capacity: location.capacity,
          username: username ?? "",
          available: location.available,
          checkIns: location.checkIns,
          yesterdayFullAt: location.yesterdayFullAt
        })
      }

      return locationsWithUsernames
    },
    []
  )

  return (
    <ListViewLayout
      title="Locations"
      form={form}
      list={locationList}
      setList={setLocationList}
      apiCall={getAllLocations}
      preprocessList={preprocessList}
      createModal={CreateLocationModal}
      card={LocationCard}
    />
  )
}
