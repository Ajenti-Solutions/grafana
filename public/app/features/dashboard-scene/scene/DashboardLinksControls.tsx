import React from 'react';

import { sanitizeUrl } from '@grafana/data/src/text/sanitize';
import { selectors } from '@grafana/e2e-selectors';
import { SceneComponentProps, SceneObjectBase, SceneObjectState } from '@grafana/scenes';
import { DashboardLink } from '@grafana/schema';
import { Tooltip } from '@grafana/ui';
import { linkIconMap } from 'app/features/dashboard/components/LinksSettings/LinkSettingsEdit';
import {
  DashboardLinkButton,
  DashboardLinksDashboard,
} from 'app/features/dashboard/components/SubMenu/DashboardLinksDashboard';
import { getLinkSrv } from 'app/features/panel/panellinks/link_srv';

interface DashboardLinksControlsState extends SceneObjectState {
  links: DashboardLink[];
  dashboardUID: string;
}

export class DashboardLinksControls extends SceneObjectBase<DashboardLinksControlsState> {
  static Component = DashboardLinksControlsRenderer;
}

function DashboardLinksControlsRenderer({ model }: SceneComponentProps<DashboardLinksControls>) {
  const { links, dashboardUID } = model.useState();
  return (
    <>
      {links.map((link: DashboardLink, index: number) => {
        const linkInfo = getLinkSrv().getAnchorInfo(link);
        const key = `${link.title}-$${index}`;

        if (link.type === 'dashboards') {
          return <DashboardLinksDashboard key={key} link={link} linkInfo={linkInfo} dashboardUID={dashboardUID} />;
        }

        const icon = linkIconMap[link.icon];

        const linkElement = (
          <DashboardLinkButton
            icon={icon}
            href={sanitizeUrl(linkInfo.href)}
            target={link.targetBlank ? '_blank' : undefined}
            rel="noreferrer"
            data-testid={selectors.components.DashboardLinks.link}
          >
            {linkInfo.title}
          </DashboardLinkButton>
        );

        return (
          <div key={key} data-testid={selectors.components.DashboardLinks.container}>
            {link.tooltip ? <Tooltip content={linkInfo.tooltip}>{linkElement}</Tooltip> : linkElement}
          </div>
        );
      })}
    </>
  );
}
