<template>
	<div v-title="heading">

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
	import Editor from "./Editor";
	import Conflict from "./Conflict";
	import Accessors from '../Mixins/accessors';

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
			},

		},
		methods: {
			async create(event) {
				event.preventDefault();

				this.commit.addFile(this.document);

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

				this.redirectToShowDocument(
					this.document.path,
					this.document.document,
					this.document.language
				);

			},

			showConflictModal() {
				$("#conflict-warning.modal").modal()
			},

			redirectToShowDocument(directory, document, language_code) {

				let params = {directory, document};
				let enabled = this.$store.state.translationEnabled;
				let isDefault = (language_code !== this.$store.state.defaultLanguage)

				if (enabled && isDefault) {
					params.language_code = language_code;
				};

				this.$router.push({name: 'document_show', params});
			}

		},
		mixins: [Accessors],
		components: {
			Editor,
			Breadcrumbs,
			Conflict
		}
	}
</script>