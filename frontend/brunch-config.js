// See http://brunch.io for documentation.
exports.files = {
	javascripts: {
		exclude: '**/*.min.js',
		joinTo: {
			'cms/javascripts/vendor.js': /^(?!src)/,
			'cms/javascripts/app.js': /^src/
		}
	},
	stylesheets: {
		joinTo: {
			'cms/stylesheets/vendor.css': /^(?!src)/,
			'cms/stylesheets/app.css': /^src/
		}
	}
};

exports.paths = {
	watched: [
		'src',
		'node_modules/simplemde/dist'
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
		out: 'public/cms/stylesheets/components.css'
	}
};

exports.npm = {
	globals: {
		$: 'jquery'
	}
};

exports.server = {
	indexPath: '/cms/index.html'
};
