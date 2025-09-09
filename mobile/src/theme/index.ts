import { colors } from './colors';
import { typography } from './typography';
import { spacing } from './spacing';

export { colors, typography, spacing };

export const theme = {
  colors,
  typography,
  spacing,
} as const;
