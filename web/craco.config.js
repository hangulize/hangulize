module.exports = {
  plugins: [{ plugin: require('@semantic-ui-react/craco-less') }],
  babel: { presets: ['@babel/preset-typescript'] },
  webpack: {
    configure: (webpackConfig) => {
      webpackConfig.module.rules.push({
        test: /\.yaml$/,
        type: 'javascript/auto',
        use: 'yaml-loader',
      })
      return webpackConfig
    },
  },
}
