export default {

	computed: {

		// quick access to route params

		directory() {
			return this.$route.params.directory;
		},

		filename() {
			return this.$route.params.filename;
		},

		// quick access to stuff in the store

		activeDirectory() {
			return this.$store.state.activeDirectory;
		},

		document() {
			return this.$store.state.activeDocument;
		},

		documents() {
			return this.$store.state.documents;
		},

		directories() {
			return this.$store.state.directories;
		},

		commit() {
			return this.$store.state.commit;
		},

		user() {
			return this.$store.state.user;
		}

	}
};