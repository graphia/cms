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
		'node_modules/simplemde/dist',
		'node_modules/bootstrap/dist',
		'node_modules/octicons',
		'node_modules/text-diff',
		'node_modules/fecha',
		'node_modules/babel-runtime/core-js/object'
	]
};

exports.modules = {
	autoRequire: {
		'cms/javascripts/app.js': ['src/javascripts/initialize']
	}
};

exports.plugins = {
	babel: {
		presets: ['es2015', 'es2016', 'es2017'],
		plugins: [
			'transform-runtime'
		]
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
		$: 'jquery',
		jQuery: 'jquery',
		Tether: 'tether',
		Popper: 'popper.js',
		Bootstrap: 'bootstrap'
	}
};
