import { PlusIcon } from "@radix-ui/react-icons";

import { Button, DataTable, Flex, useTable } from "@raystack/apsara";
import { useNavigate, useParams } from "react-router-dom";
import PageHeader from "~/components/page-header";

const defaultPageHeader = {
  title: "Organizations",
  breadcrumb: [],
};

export const OrganizationsTokenHeader = ({
  header = defaultPageHeader,
  ...props
}: any) => {
  const navigate = useNavigate();
  const { filteredColumns, table } = useTable();
  const isFiltered = filteredColumns.length > 0;
  let { organisationId, billingaccountId } = useParams();

  return (
    <>
      <PageHeader
        title={header.title}
        breadcrumb={header.breadcrumb}
        {...props}
      >
        {isFiltered ? <DataTable.ClearFilter /> : <DataTable.FilterOptions />}
        <DataTable.ViewOptions />
        <DataTable.GloabalSearch placeholder="Search transaction..." />
        <Button
          variant="secondary"
          onClick={() =>
            navigate(
              `/organisations/${organisationId}/billingaccounts/${billingaccountId}/tokens/add`
            )
          }
          style={{ width: "100%" }}
        >
          <Flex
            direction="column"
            align="center"
            style={{ paddingRight: "var(--pd-4)" }}
          >
            <PlusIcon />
          </Flex>
          Add tokens
        </Button>
      </PageHeader>
    </>
  );
};
