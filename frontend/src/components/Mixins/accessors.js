export default {

	computed: {

		// quick access to route params

		directory() {
			return this.$route.params.directory;
		},

		filename() {
			return this.$route.params.filename;
		},

		params() {
			return this.$route.params;
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

		currentUser() {
			return this.$store.state.user;
		}

	}
};