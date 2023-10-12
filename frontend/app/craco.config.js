const { whenTest } = require('@craco/craco')

module.exports = {
	babel: {
		plugins: [
			...whenTest(() => [['@babel/plugin-transform-modules-commonjs', { allowTopLevelThis: true }]], [])
		]
	}
}