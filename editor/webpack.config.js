const path = require('path');
const requirejs = require('requirejs');

requirejs.config({
  basePath: './src',
});

module.exports = {
  entry: {
    index: './src/index.js',
  },

  output: {
    filename: '[name].js',
    path: path.resolve(__dirname, 'dist'),
  },

  externals: {
    hangulize: 'hangulize',
  },

  optimization: {
    minimize: false,
  },

  module: {
    rules: [{
      test: /\.css$/,
      use: ['style-loader', 'css-loader'],
    }],
  },
};
