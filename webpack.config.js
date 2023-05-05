const { VueLoaderPlugin } = require('vue-loader');
const FileLoadPlugin = require('file-loader')
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');
const { DefinePlugin } = require('webpack');

module.exports = (env, argv) => {
  var mode = "production"
  if (argv.mode) {
    mode = "development"
  }

  var config = {
    entry: {
      results: {
        import: './src/results.js',
      },
    },
    mode: mode,
    output: {
      filename: '[name].[contenthash].js',
      path: path.resolve(__dirname, './dist'),
      clean: true
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
            filename: 'images/icon/[name][ext]'
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
      new DefinePlugin({
        __VUE_OPTIONS_API__: true,
        __VUE_PROD_DEVTOOLS__: false,
      }),
      new VueLoaderPlugin(),
      new HtmlWebpackPlugin({
        title: "results.roadrunners.club | Roadrunners Timing System",
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
