import React, { useMemo } from 'react';

import {
  Avatar,
  CellProps,
  Column,
  FetchDataFunc,
  Icon,
  InteractiveTable,
  Pagination,
  Stack,
  Text,
  Tooltip,
} from '@grafana/ui';

type AnonUserDTO = {
  id: number;
  login: string;
  email: string;
  name: string;
  avatarUrl: string;
  useragent: string;
  lastSeenAtAge: string;
  ipAddress: string;
};

type Cell<T extends keyof AnonUserDTO = keyof AnonUserDTO> = CellProps<AnonUserDTO, AnonUserDTO[T]>;

interface AnonUsersTableProps {
  users: AnonUserDTO[];
  showPaging?: boolean;
  totalPages: number;
  onChangePage: (page: number) => void;
  currentPage: number;
  fetchData?: FetchDataFunc<AnonUserDTO>;
}

export const AnonUsersTable = ({
  users,
  showPaging,
  totalPages,
  onChangePage,
  currentPage,
  fetchData,
}: AnonUsersTableProps) => {
  const columns: Array<Column<AnonUserDTO>> = useMemo(
    () => [
      {
        id: 'avatarUrl',
        header: '',
        cell: ({ cell: { value } }: Cell<'avatarUrl'>) => value && <Avatar src={value} alt={'User avatar'} />,
      },
      {
        id: 'login',
        header: 'Login',
        cell: ({ cell: { value } }: Cell<'login'>) => value,
        sortType: 'string',
      },
      {
        id: 'useragent',
        header: 'User Agent',
        cell: ({ cell: { value } }: Cell<'useragent'>) => value,
        sortType: 'string',
      },
      {
        id: 'lastSeenAtAge',
        header: 'Last active',
        headerTooltip: {
          content: 'Time since user was seen using Grafana',
          iconName: 'question-circle',
        },
        cell: ({ cell: { value } }: Cell<'lastSeenAtAge'>) => {
          return <>{value && <>{value === '10 years' ? <Text color={'disabled'}>Never</Text> : value}</>}</>;
        },
        sortType: (a, b) => new Date(a.original.lastSeenAtAge!).getTime() - new Date(b.original.lastSeenAtAge!).getTime(),
      },
      {
        id: 'ipAddress',
        header: 'Origin (IP address)',
        cell: ({ cell: { value } }: Cell<'ipAddress'>) => value,
        sortType: 'string',
      },
      {
        id: 'edit',
        header: '',
        cell: ({ row: { original } }: Cell) => {
          return (
            <a href={`admin/users/edit/${original.id}`} aria-label={`Edit team ${original.name}`}>
              <Tooltip content={'Edit user'}>
                <Icon name={'pen'} />
              </Tooltip>
            </a>
          );
        },
      },
    ],
    []
  );
  return (
    <Stack direction={'column'} gap={2}>
      <InteractiveTable columns={columns} data={users} getRowId={(user) => String(user.id)} fetchData={fetchData} />
      {showPaging && (
        <Stack justifyContent={'flex-end'}>
          <Pagination numberOfPages={totalPages} currentPage={currentPage} onNavigate={onChangePage} />
        </Stack>
      )}
    </Stack>
  );
};
