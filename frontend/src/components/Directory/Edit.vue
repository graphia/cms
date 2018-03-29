<template>

	<div class="edit-directory col-md-12" v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<h1>Edit {{ this.activeDirectory.title }}</h1>

		<!-- edit directory form -->
		<form :id="formID" @submit="updateDirectory">

			<TitleField/>

			<div class="form-group">
				<label class="form-control-label" for="description">Description</label>
				<textarea
					name="description"
					class="form-control"
					v-model="activeDirectory.description"
					placeholder="A set of detailed step-by-step instructions compiled to help workers carry out complex routine operations"
				/>
				<p id="display-text-explanation" class="form-text text-muted">
					The description text is displayed in the summary of document types
					on the published homepage. It should be <em>short</em> and <em>concise</em>.
				</p>
			</div>

			<div>
				<MinimalMarkdownEditor/>
			</div>

			<div class="btn-toolbar" role="toolbar">
				<input
					type="submit"
					class="btn btn-success mx-2"
					value="Update directory"
				/>

				<router-link class="btn btn-secondary" :to="{name: 'directory_index', params: {directory}}">
					Cancel
				</router-link>
			</div>

		</form>
		<!-- /edit directory form -->

	</div>

</template>

<script lang="babel">
	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSDirectory from '../../javascripts/models/directory.js';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import MinimalMarkdownEditor from './Editor';
	import TitleField from './Metadata/TitleField';
	import Accessors from '../Mixins/accessors';

	export default {
		name: "DirectoryEdit",
		data() {
			return {
				markdownLoaded: false,
				formID: "update-directory",
				valid: false,
				title: "Edit directory"
			};
		},
		async created() {
			await this.$store.dispatch("getDirectory", this.directory);
			this.markdownLoaded = true;
		},
		methods: {
			async updateDirectory(event) {
				event.preventDefault();

				await this.$store.dispatch("initializeCommit");

				this.commit.message = `Updated ${this.activeDirectory.title} metadata`;
				this.commit.addDirectory(this.activeDirectory);

				let response = await this.activeDirectory.update(this.commit);

				if (!checkResponse(response.status)) {
					console.error(response.status);
					return;
				};

				let json = await response.json();

				this.$store.commit("setLatestRevision", json.oid);

				// directory updated successfully, show a message
				this.$store.state.broadcast.addMessage(
					"success",
					"Directory Updated",
					`Your changes to ${this.activeDirectory.title} have been saved`,
					3
				);

				// redirect to the new directory's index page
				this.redirectToIndex(this.directory);
				return;
			},
			redirectToIndex(directory) {
				this.$router.push({name: 'directory_index', params: {directory}});
			},

		},
		computed: {

			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						this.activeDirectory.title || this.directory,
						"directory_index",
						{directory: this.directory}
					),
					new CMSBreadcrumb(
						"Edit",
						"directory_edit",
						{directory: this.directory}
					)
				];
			}

		},
		components: {
			MinimalMarkdownEditor,
			TitleField,
			Breadcrumbs
		},
		mixins: [Accessors]
	};
</script>