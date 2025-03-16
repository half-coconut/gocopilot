import Spinner from "../../ui/Spinner";
import InterfaceRow from "./InterfaceRow";
import { useInterfaces } from "./useInterfaces";
import Table from "../../ui/Table";
import Menus from "../../ui/Menus";
// import { useSearchParams } from "react-router-dom";
import Empty from "../../ui/Empty";

function InterfaceTable() {
  const { isLoading, interfaceItems, total } = useInterfaces();
  // const [searchParams] = useSearchParams();

  if (isLoading) return <Spinner />;
  if (!total) return <Empty resourceName="cabins" />;

  // 1） FILTER
  // const filterValue = searchParams.get("discount") || "all";

  // let filteredCabins;
  // if (filterValue === "all") filteredCabins = cabins;
  // if (filterValue === "no-discount")
  //   filteredCabins = cabins.filter((cabin) => cabin.discount === 0);
  // if (filterValue === "with-discount")
  //   filteredCabins = cabins.filter((cabin) => cabin.discount > 0);

  // // 2) SORT
  // const sortBy = searchParams.get("sortBy") || "startDate-asc";
  // const [field, direction] = sortBy.split("-");
  // const modifier = direction === "asc" ? 1 : -1;
  // const sortedCabins = filteredCabins.sort(
  //   (a, b) => (a[field] - b[field]) * modifier
  // );

  return (
    <Menus>
      <Table columns="1fr 1fr 1fr  1fr 1fr 1fr 1fr 1fr">
        <Table.Header>
          {/* <div></div> */}
          <div>Name</div>
          <div>Project</div>
          <div>type</div>
          <div>Creator</div>
          <div>Updater</div>
          <div>Ctime</div>
          <div>Utime</div>
          <div></div>
        </Table.Header>

        <Table.Body
          // data={cabins}
          // data={filteredCabins}
          data={interfaceItems}
          render={(interfaceItem) => (
            <InterfaceRow
              interfaceItem={interfaceItem}
              key={interfaceItem.id}
            />
          )}
        />
      </Table>
    </Menus>
  );
}
export default InterfaceTable;
