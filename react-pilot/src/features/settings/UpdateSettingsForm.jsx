import Form from "../../ui/Form.jsx";
import FormRow from "../../ui/FormRow";
import Input from "../../ui/Input";
// import { useSettings } from "./useSettings.js";
// import Spinner from "../../ui/Spinner.jsx";
// import { useUpdateSetting } from "./useUpdateSetting.js";

function UpdateSettingsForm() {
  // const {
  //   isLoading,
  //   // error,
  //   settings: {
  //     minBookingLength,
  //     maxBookingLength,
  //     maxGuestsPerBooking,
  //     breakfastPrice,
  //   } = {},
  // } = useSettings();

  // const { isUpdating, updateSetting } = useUpdateSetting();

  // if (isLoading) return <Spinner />;

  // function handleUpdate(e, field) {
  //   const { value } = e.target;
  //   if (!value) return;
  //   updateSetting({ [field]: value });
  // }

  const data = {
    minUsers: 5,
    maxUsers: 90,
    spawnRate: 8,
    runTime: 15,
  };

  return (
    <Form>
      <FormRow label="Minimum Concurrent Users">
        <Input
          type="number"
          id="min-users"
          // disabled={isUpdating}
          defaultValue={data.minUsers}
          // onBlur={(e) => handleUpdate(e, "minBookingLength")}
        />
      </FormRow>

      <FormRow label="Maximum Concurrent Users">
        <Input
          type="number"
          id="max-users"
          // disabled={isUpdating}
          defaultValue={data.maxUsers}
          // onBlur={(e) => handleUpdate(e, "maxBookingLength")}
        />
      </FormRow>

      <FormRow label="Spawn Rate">
        <Input
          type="number"
          id="spawn-rate"
          // disabled={isUpdating}
          defaultValue={data.spawnRate}
          // onBlur={(e) => handleUpdate(e, "maxGuestsPerBooking")}
        />
      </FormRow>

      <FormRow label="Run Time">
        <Input
          type="number"
          id="run-time"
          // disabled={isUpdating}
          defaultValue={data.runTime}
          // onBlur={(e) => handleUpdate(e, "breakfastPrice")}
        />
      </FormRow>
    </Form>
  );
}

export default UpdateSettingsForm;
