import React from 'react';
import { Table } from '@/design.system';
import { ManagedSource, Namespace } from '@/types';
import { SourcesTableRow } from './sources.table.row';
import { SourcesTableHeader } from './sources.table.header';
import { EmptyList } from '@/components/lists';
import { OVERVIEW } from '@/utils';

type TableProps = {
  data: ManagedSource[];
  namespaces?: Namespace[];
  onRowClick: (source: ManagedSource) => void;
  sortSources?: (condition: string) => void;
  filterSourcesByKind?: (kinds: string[]) => void;
  filterSourcesByNamespace?: (namespaces: string[]) => void;
  toggleActionStatus?: (ids: string[], disabled: boolean) => void;
};

export const ManagedSourcesTable: React.FC<TableProps> = ({
  data,
  namespaces,
  onRowClick,
  sortSources,
  toggleActionStatus,
  filterSourcesByKind,
  filterSourcesByNamespace,
}) => {
  const currentPageRef = React.useRef(1);

  function onPaginate(pageNumber: number) {
    currentPageRef.current = pageNumber;
  }

  function renderEmptyResult() {
    return <EmptyList title={OVERVIEW.EMPTY_SOURCE} />;
  }

  function renderTableHeader() {
    return (
      <SourcesTableHeader
        data={data}
        namespaces={namespaces}
        sortSources={sortSources}
        toggleActionStatus={toggleActionStatus}
        filterSourcesByKind={filterSourcesByKind}
        filterSourcesByNamespace={filterSourcesByNamespace}
      />
    );
  }

  return (
    <>
      <Table<ManagedSource>
        data={data}
        renderTableHeader={renderTableHeader}
        onPaginate={onPaginate}
        renderEmptyResult={renderEmptyResult}
        renderTableRows={(item, index) => (
          <SourcesTableRow
            onRowClick={onRowClick}
            data={data}
            item={item}
            index={index}
          />
        )}
      />
    </>
  );
};
