import { createLocation, getAllLocationsPaged, getUser } from "api-wrapper"
import { useForm } from "react-hook-form"
import CreateModal from "components/modals/CreateModal"
import FormInput from "components/forms/FormInput"
import ListViewLayout, { type List } from "layouts/ListViewLayout"
import { useCallback, useState } from "react"
import LocationCard from "components/cards/LocationCard"
import { type ICreateModalProps } from "interfaces/ICreateModalProps"
import {
  type Role,
  type CreateLocationDto,
  type Location,
  type PaginatedLocationsDto,
  type User
} from "api-wrapper/types/apiTypes"
import { AuthRedirecter } from "contexts/authContext"
import UserInputs from "components/forms/UserInputs"
import TimeZoneInput from "components/forms/TimeZoneInput"

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
      {/* jscpd:ignore-start */}
      <FormInput
        label="Name"
        type="text"
        placeholder="Name"
        required
        register={register("name")}
        errors={errors.name}
      />
      <FormInput
        label="Capacity"
        type="number"
        placeholder={10}
        required
        register={register("capacity")}
        errors={errors.capacity}
      />
      <TimeZoneInput register={register("timeZone")} />
      <UserInputs
        required={true}
        register={register}
        watch={watch}
        errors={errors}
      />
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
  const redirects = new Map<Role, string>([["default", "/settings"]])
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
    async (responseData: PaginatedLocationsDto) => {
      const locationsWithUsernames: List<LocationWithUsername> = {
        data: [],
        pagination: {
          current: responseData.pagination.current,
          total: responseData.pagination.total
        }
      }

      for (const location of responseData.data) {
        const username = ((await getUser(location.userId)).data as User)
          .username

        ;(locationsWithUsernames.data as LocationWithUsername[]).push({
          id: location.id,
          name: location.name,
          normalizedName: location.normalizedName,
          capacity: location.capacity,
          timeZone: location.timeZone,
          username: username,
          available: location.available,
          yesterdayFullAt: location.yesterdayFullAt
        })
      }

      return locationsWithUsernames
    },
    []
  )

  return (
    <AuthRedirecter redirects={redirects}>
      <ListViewLayout
        title="Locations"
        form={form}
        list={locationList}
        setList={setLocationList}
        apiCall={getAllLocationsPaged}
        preprocessList={preprocessList}
        createModal={CreateLocationModal}
        card={LocationCard}
      />
    </AuthRedirecter>
  )
}
