import { getBackendSrv } from '@grafana/runtime';

export interface ServerStat {
  activeAdmins: number;
  activeEditors: number;
  activeSessions: number;
  activeUsers: number;
  activeAnonymousUsers?: number;
  activeViewers: number;
  admins: number;
  alerts: number;
  dashboards: number;
  datasources: number;
  editors: number;
  orgs: number;
  playlists: number;
  snapshots: number;
  stars: number;
  tags: number;
  users: number;
  viewers: number;
}

export const getServerStats = async (): Promise<ServerStat | null> => {
  let resp = await getBackendSrv().get('api/admin/anonstats').then((res) => {
    return res;
  }).catch((err) => {
    console.error(err);
    // FIXME:
    // return null;
  });
  return getBackendSrv()
    .get('api/admin/stats')
    .then((res) => {
      console.log(`resp`);
      console.log(resp);
      console.log(`res`);
      console.log(res);
      if (resp) {
        res.activeAnonymousUsers = resp.length;
      }
      return res;
    })
    .catch((err) => {
      console.error(err);
      return null;
    });
};
