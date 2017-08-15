<template>
	<section>

		<form id="edit-document-form" @submit="update">
			<h1>{{ heading }}</h1>
			<Editor/>
		</form>

	</section>
</template>

<script lang="babel">
	import Editor from "../components/Editor";

	export default {
		name: "DocumentEdit",
		data() {
			return {
				markdownLoaded: false,
				form: null,
				formID: "edit-document-form"
			};
		},
		async created() {
			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");

			// retrieve the document and add it to vuex's store
			await this.$store.dispatch("editDocument", {directory: this.directory, filename: this.filename});

			// FIMXE use the bus 🚌
			this.markdownLoaded = true;

		},
		computed: {

			// quick access to things in the store
			document() {
				return this.$store.state.activeDocument;
			},
			commit() {
				return this.$store.state.commit;
			},

			// quick access to route params
			directory() {
				return this.$route.params.directory;
			},
			filename() {
				return this.$route.params.filename;
			},

			heading() {
				let title = this.document.title;
				if (title) {
					return title;
				} else {
					return "No title";
				}
			}
		},
		methods: {
			update(event) {
				event.preventDefault();

				this.document.update(this.commit)
					.then((response) => {
						console.debug("Document saved, redirecting to 'document_show'");
						this.redirectToShowDocument(this.directory, this.filename);
					});
			},
			redirectToShowDocument(directory, filename) {
				this.$router.push({
					name: 'document_show',
					params:{directory, filename}
				});
			}
		},
		components: {
			Editor
		}
	}
</script>