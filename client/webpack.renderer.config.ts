import type { Configuration } from 'webpack';

import { rules } from './webpack.rules';
import { plugins } from './webpack.plugins';

rules.push({
  test: /\.css$/,
  use: [{ loader: 'style-loader' }, {
    loader: 'css-loader',
    options: {
      modules: {
        auto: (resourcePath: string) => !resourcePath.includes('node_modules'),
        localIdentName: "[name]_[local]__[hash:base64:5]"
      }
    }
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