const webpack = require('webpack');
const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const extractCSS = new ExtractTextPlugin('allstyles.css');


module.exports = {
    entry: {'main': './src/app.js'},
    output: {
        path: path.resolve(__dirname, './root'),
        filename: 'bundle.js',
        publicPath: 'root/'
    },
    module: {
        loaders: [
            {
                test: /\.html$/,
                loader: 'file-loader?name=[name].[ext]',
            },
            {
                test: /\.jsx?$/,
                exclude: /node_modules/,
                loader: 'babel-loader',
            options: {
                presets: ['babel-preset-react', 'babel-preset-env']
            }
        },
        {
          test: /\.css$/, use: extractCSS.extract(['css-loader?minimize'])
        }
    ],
},
plugins: [
    extractCSS,
    new webpack.ProvidePlugin({
        Popper: ['popper.js', 'default']
    }),
]
};
