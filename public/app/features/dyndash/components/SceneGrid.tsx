import React, { FC } from 'react';
import { SceneView } from './SceneView';
import { SceneVizView } from './VizPanel';
import { SceneItem } from '../models';
import { ComponentView } from './ComponentView';
import { GRID_CELL_HEIGHT } from 'app/core/constants';

export interface Props {
  panels: SceneItem[];
}

export const SceneGrid: FC<Props> = ({ panels }) => {
  return (
    <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
      {panels.map(panel => (
        <div key={panel.id} id={panel.id} style={getSceneItemStyles(panel)}>
          <ScenePanelView panel={panel} key={panel.id} />
        </div>
      ))}
    </div>
  );
};

function getSceneItemStyles(panel: SceneItem) {
  return {
    width: `${(24 / panel.gridPos.w) * 100}%`,
    height: `${panel.gridPos.h * GRID_CELL_HEIGHT}px`,
  }
}

interface PanelProps {
  panel: SceneItem;
}

const ScenePanelView: FC<PanelProps> = ({ panel }) => {
  switch (panel.type) {
    case 'viz':
      return <SceneVizView panel={panel} />;
    case 'scene':
      return <SceneView model={panel} />;
    case 'component':
      return <ComponentView panel={panel} />;
  }
};
