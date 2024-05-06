import image from '@rollup/plugin-image';

export default {
  input: './src/cursors.js',
  output: {
    file: './dist/cursors.js',
    format: 'iife',
  },  
  plugins: [image()],
};
