import React from 'react';
import { Flamebearer } from '@models/flamebearer';
import Dropdown, { MenuItem, SubMenu } from '@ui/Dropdown';
import DiffLegend from './DiffLegend';
import styles from './Header.module.css';
import {
  FlamegraphPalette,
  DefaultPalette,
  ColorBlindPalette,
} from './colorPalette';
import { DiffLegendPaletteDropdown } from './DiffLegendPaletteDropdown';

interface HeaderProps {
  format: Flamebearer['format'];
  units: Flamebearer['units'];

  palette: FlamegraphPalette;
  setPalette: (p: FlamegraphPalette) => void;
  ExportData: () => React.ReactElement;
}
export default function Header(props: HeaderProps) {
  const { format, units, ExportData, palette, setPalette } = props;

  const unitsToFlamegraphTitle = {
    objects: 'amount of objects in RAM per function',
    bytes: 'amount of RAM per function',
    samples: 'CPU time per function',
  };

  const getTitle = () => {
    switch (format) {
      case 'single': {
        return (
          <div>
            <div
              className={`${styles.row} ${styles['flamegraph-title']}`}
              role="heading"
              aria-level={2}
            >
              {unitsToFlamegraphTitle[units] && (
                <>Frame width represents {unitsToFlamegraphTitle[units]}</>
              )}
            </div>
          </div>
        );
      }

      case 'double': {
        return (
          <>
            <div className={styles.row} role="heading" aria-level={2}>
              Base graph: left - Comparison graph: right
            </div>
            <DiffLegendPaletteDropdown
              palette={palette}
              onChange={setPalette}
            />
          </>
        );
      }

      default:
        throw new Error(`unexpected format ${format}`);
    }
  };

  const title = getTitle();

  return (
    <div className={styles['flamegraph-header']}>
      <div>{title}</div>
      <div className={styles.buttons}>
        <ExportData />
      </div>
    </div>
  );
}
