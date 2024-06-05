import image from '@rollup/plugin-image';
import { nodeResolve } from '@rollup/plugin-node-resolve';

export default {
  input: './src/cursors.js',
  output: {
    file: './dist/cursors.js',
    format: 'iife',
  },
  plugins: [image(), nodeResolve()],
};
