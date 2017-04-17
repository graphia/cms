// See http://brunch.io for documentation.
exports.files = {
	javascripts: {
		exclude: '**/*.min.js',
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
				'node_modules/bootstrap/scss'
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

exports.npm = {
	globals: {
		$: 'jquery'
	}
}