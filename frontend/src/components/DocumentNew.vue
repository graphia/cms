<template>
	<section>

		<form id="new-document-form" @submit="create">
			<h1>{{ heading }}</h1>
			<Editor
				:formID="formID"
				:submitButtonText="submitButtonText"
				:newFile="true"
				:formCancellationRedirectParams="formCancellationRedirectParams"
			/>
		</form>

	</section>
</template>

<script lang="babel">
	import Editor from "../components/Editor";

	import checkResponse from "../javascripts/response.js";

	export default {
		name: "DocumentNew",
		data() {
			return {
				formID: "new-document-form",
				submitButtonText: "Create",
			};
		},
		async created() {

			console.debug("new doc...");

			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");

			// initialize a fresh new document
			let doc = await this.$store.dispatch("initializeDocument", this.directory);
			this.$store.commit("setActiveDocument", doc);


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

			heading() {
				let title = this.document.title;
				if (title) {
					return title;
				} else {
					return "New Document";
				}
			},

			formCancellationRedirectParams() {
				return {
					name: 'document_index'
				};
			}

		},
		methods: {
			async create(event) {
				event.preventDefault();

				let response = await this.document.create(this.commit);

				if (!checkResponse(response.status)) {
					throw("could not create document");
				};

				console.debug("Document saved, redirecting to 'document_show'");
				this.redirectToShowDocument(this.document.path, this.document.filename);

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