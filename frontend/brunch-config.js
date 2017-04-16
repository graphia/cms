// See http://brunch.io for documentation.
exports.files = {
	javascripts: {
		joinTo: {
			'javascripts/vendor.js': /^(?!src)/,
			'javascripts/app.js': /^src/
		}
	},
	stylesheets: {
		joinTo: {
			'stylesheets/vendor.css': /^(?!src)/,
			'stylesheets/app.css': /^src/
		}
	}
};

exports.paths = {
	watched: [
		'src'
		/*,'node_modules/simplemde/dist'*/
	]
};

exports.plugins = {
	babel: {
		presets: ['latest'],
		plugins: ['transform-runtime']
	},
	sass: {
		mode: 'native',
		options: {
			includePaths: [
				'node_modules/milligram/src'
			],
			sourceMapEmbed: true,
		}
	},
	vue: {
		extractCSS: true,
		indentedSyntax: true,
		out: 'public/stylesheets/components.css'
	}
};
