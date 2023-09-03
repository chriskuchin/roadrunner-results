const { VueLoaderPlugin } = require('vue-loader');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const { DefinePlugin } = require('webpack');

const { GenerateSW } = require('workbox-webpack-plugin');

const statusPlugin = {
  fetchDidSucceed: ({ response }) => {
    if (response.status >= 500) {
      // Throwing anything here will trigger fetchDidFail.
      throw new Error('Server error.');
    }
    // If it's not 5xx, use the response as-is.
    return response;
  },
};

module.exports = (env, argv) => {
  var mode = "production"
  if (argv.mode) {
    mode = "development"
  }

  var config = {
    watch: true,
    entry: {
      results: {
        import: './src/results.js',
      },
    },
    mode: mode == "development",
    output: {
      filename: '[name].[contenthash].js',
      path: path.resolve(__dirname, './dist'),
      clean: true,
      publicPath: '/',
    },
    module: {
      rules: [
        {
          test: /\.json/,
          type: 'asset/resource',
          generator: {
            filename: '[name][ext]'
          }
        },
        {
          test: /\.(jpg|png)$/,
          type: 'asset/resource',
          generator: {
            filename: 'images/[name][ext]'
          }
        },
        {
          test: (value => value.includes("icon-")),
          type: 'asset/resource',
          generator: {
            filename: 'images/icons/[name][ext]'
          }
        },
        {
          test: /\.vue$/,
          loader: 'vue-loader'
        },
        {
          test: /\.css$/,
          use: [
            'vue-style-loader',
            'css-loader'
          ]
        },
        {
          test: /\.s[ac]ss$/i,
          use: [
            MiniCssExtractPlugin.loader,
            "css-loader",
            "sass-loader",
          ],
        },
      ]
    },
    plugins: [
      new GenerateSW({
        skipWaiting: true,
        clientsClaim: true,
        runtimeCaching: [{
          urlPattern: new RegExp('/api/.*'),
          method: 'POST',
          handler: 'NetworkOnly',
          options: {
            backgroundSync: {
              name: 'api-retry',
              options: {
                maxRetentionTime: 24 * 60,
              },
            },
            plugins: [statusPlugin]
          }
        }, {
          urlPattern: new RegExp('/api/.*'),
          method: 'GET',
          handler: 'NetworkFirst'
        }]
      }),
      new DefinePlugin({
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false,
      }),
      new VueLoaderPlugin(),
      new HtmlWebpackPlugin({
        title: "rslts.run | Roadrunners Timing System",
        filename: "index.html",
        template: 'src/index.ejs',
        favicon: 'src/assets/images/favicon.ico',
        meta: {
          viewport: "initial-scale=1, maximum-scale=1",
        },
        templateParameters: {
          mode: mode,
        }
      }),
      new MiniCssExtractPlugin({
        filename: "[name].[contenthash].css"
      }),
    ],
  }

  return config
};
