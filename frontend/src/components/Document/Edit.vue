<template>
	<div v-title="heading">
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
	import filenameFromLanguageCode from '../../javascripts/utilities/filename-from-language-code.js';

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

			const filename = filenameFromLanguageCode(this.params.language_code)

			// retrieve the document and make it Active
			await this.$store.dispatch("editDocument", {
				directory: this.params.directory,
				document: this.params.document,
				filename: filename
			});

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

				let dir_title, doc_title;

				// if we have it, use the metadata provided directory and
				// document title
				if (this.document.directory_info) {
					dir_title = this.document.directory_info.title;
					doc_title = this.document.title;
				};
				return [

					new CMSBreadcrumb(
						dir_title || this.directory,
						"directory_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						doc_title || this.params.document,
						"document_show",
						{directory: this.directory, document: this.params.document}
					),
					new CMSBreadcrumb(
						"Edit",
						"document_edit",
						{directory: this.directory, document: this.params.document}
					)
				];
			}
		},
		methods: {

			async update(event) {

				event.preventDefault();

				this.commit.addFile(this.document);

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

				let json = await response.json();

				await this.$store.commit(
					"setLatestRevision",
					json.latest_revision
				);

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