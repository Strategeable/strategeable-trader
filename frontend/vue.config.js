/** @type {import('@vue/cli-service').ProjectOptions} */
module.exports = {
  chainWebpack(config) {
    config.module
      .rule('js')
      .use('babel-loader')
      .tap(() => {
        return {
          rootMode: "upward"
        }
      })
  }
};
