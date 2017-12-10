<template>
	<div>
		<Breadcrumbs :levels="breadcrumbs"/>

		<Conflict/>

		<section>

			<form id="edit-document-form" @submit="update">
				<h1>{{ heading }}</h1>
				<Editor
					:formID="formID"
					:submitButtonText="submitButtonText"
					:formCancellationRedirectParams="formCancellationRedirectParams"
				/>
			</form>

		</section>

	</div>
</template>

<script lang="babel">
	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import Editor from "./Editor";
	import Conflict from "./Conflict";
	import Accessors from '../Mixins/accessors';

	import checkResponse from "../../javascripts/response.js";
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	export default {
		name: "DocumentEdit",
		data() {
			return {
				markdownLoaded: false,
				formID: "edit-document-form",
				submitButtonText: "Update"
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
		mixins: [Accessors],
		computed: {

			heading() {
				let title = this.document.title;
				if (title) {
					return title;
				} else {
					return "No title";
				}
			},

			formCancellationRedirectParams() {
				return {
					name: 'document_show',
					params: {
						directory: this.directory,
						filename: this.document.filename
					}
				};
			},
			breadcrumbs() {

				let directory_title, filename;

				if (this.document.directory_info) {
					directory_title = this.document.directory_info.title;
					filename = this.document.title;
				};

				return [

					new CMSBreadcrumb(
						directory_title || this.directory,
						"document_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						filename || this.filename,
						"document_show",
						{directory: this.document.path, document: this.document.filename}
					),
					new CMSBreadcrumb(
						"Edit",
						"document_edit",
						{directory: this.directory, document: (filename || this.filename)}
					)
				];
			}
		},
		methods: {

			async update(event) {

				event.preventDefault();

				let response = await this.document.update(this.commit);

				if (!checkResponse(response.status)) {

					if (response.status == 409) {
						this.showConflictModal();
						return;
					};

					// any other error
					throw("could not update document", response);
					return;
				};

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