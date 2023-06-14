const path = require('path')
const HTMLWebpackPlugin = require('html-webpack-plugin')
const {CleanWebpackPlugin} = require('clean-webpack-plugin')
module.exports= {
    context: path.resolve(__dirname,'src'),
    mode:'development',
    entry: {
        main: './frond.ts',
    },
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js'],
    },
    output: {
        filename: '[name].[contenthash].ts',
        path: path.resolve(__dirname, 'dist')
    },
    devServer: {
        port: 8888
    },
    plugins: [
        new HTMLWebpackPlugin({
            template: './zadanie5.html'
        }),
        new CleanWebpackPlugin()
    ],
}
