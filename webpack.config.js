const { VueLoaderPlugin } = require('vue-loader');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

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
            clean: true,
        },
        module: {
            rules: [
                {
                    test: /\.(jpg|png)$/,
                    use: {
                        loader: 'url-loader',
                    },
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
            new VueLoaderPlugin(),
            new HtmlWebpackPlugin({
                title: "results.roadrunners.club | Roadrunners Timing System",
                filename: "index.html",
                template: 'src/index.ejs',
                favicon: 'src/assets/images/favicon.ico',
                meta: {
                    viewport: "initial-scale=1, maximum-scale=1",
                }
            }),
            new MiniCssExtractPlugin({
                filename: "[name].[contenthash].css"
            }),
        ],
    }

    return config
};
