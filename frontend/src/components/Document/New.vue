<template>
	<div>

		<Breadcrumbs :levels="breadcrumbs"/>

		<Conflict/>

		<h1>{{ heading }}</h1>

		<form id="new-document-form" @submit="create">
			<Editor
				:formID="formID"
				:submitButtonText="submitButtonText"
				:newFile="true"
				:formCancellationRedirectParams="formCancellationRedirectParams"
			/>
		</form>

	</div>

</template>

<script lang="babel">
	import checkResponse from "../../javascripts/response.js";
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import Editor from "../Editor";
	import Conflict from "./Conflict";

	export default {
		name: "DocumentNew",
		data() {
			return {
				formID: "new-document-form",
				submitButtonText: "Create",
			};
		},
		async created() {

			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");
			this.$store.dispatch("getDirectory", this.directory);

			// initialize a fresh new document
			this.$store.commit("initializeDocument", this.directory);


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
			},

			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						this.$store.state.activeDirectory.title || this.directory,
						"document_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						"New Document",
						"document_new",
						{directory: this.directory}
					)
				];
			}

		},
		methods: {
			async create(event) {
				event.preventDefault();

				let response = await this.document.create(this.commit);

				if (!checkResponse(response.status)) {

					if (response.status == 409) {
						this.showConflictModal();
						return;
					};

					// any other error
					throw("could not create document", response);
					return;
				};

				console.debug("Document saved, redirecting to 'document_show'");
				this.redirectToShowDocument(this.document.path, this.document.filename);

			},

			showConflictModal() {
				$("#conflict-warning.modal").modal()
			},

			redirectToShowDocument(directory, filename) {
				this.$router.push({
					name: 'document_show',
					params:{directory, filename}
				});
			}

		},
		components: {
			Editor,
			Breadcrumbs,
			Conflict
		}
	}
</script>