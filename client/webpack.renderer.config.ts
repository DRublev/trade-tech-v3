import type { Configuration } from 'webpack';

import path from 'path';
import { rules } from './webpack.rules';
import { plugins } from './webpack.plugins';

rules.push({
  // exclude: /(node_modules|\.webpack|\@radix\-ui)/,
  exclude: path.resolve('./node_modules/'),
  test: /\.css$/,
  use: [{ loader: 'style-loader' }, {
    loader: 'css-loader', 
    options: {
      importLoaders: 1,
      modules: true,
      // localIdentName: "[name]__[local]___[hash:base64:5]"
    },
  }],
});

export const rendererConfig: Configuration = {
  module: {
    rules,
  },
  plugins,
  resolve: {
    extensions: ['.js', '.ts', '.jsx', '.tsx', '.css'],
  },
  devServer: { open: false }
};