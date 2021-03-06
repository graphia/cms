<template>

	<div class="new-directory" v-title="title">

		<Breadcrumbs :levels="breadcrumbs"/>

		<h1>Create a new directory</h1>

		<!-- new directory form -->
		<form :id="formID" @submit="createDirectory">

			<TitleField/>

			<div class="form-group">
				<label class="form-control-label" for="path">Path</label>
				<input
					name="path"
					class="form-control"
					placeholder="operating-procedures"
					v-model="activeDirectory.path"
					required="true"
					readonly="true"
				/>
			</div>

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

			<div class="btn-toolbar">
				<input
					type="submit"
					class="btn btn-success mx-2"
					value="Create directory"
					:disabled="!valid"
				/>
				<router-link class="btn btn-secondary" :to="{name: 'home'}">
					Cancel
				</router-link>
			</div>

		</form>
		<!-- /new directory form -->

	</div>

</template>

<script lang="babel">
	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSDirectory from '../../javascripts/models/directory.js';
	import slugify from '../../javascripts/utilities/slugify.js';
	import CMSBreadcrumb from '../../javascripts/models/breadcrumb.js';

	import Breadcrumbs from '../Utilities/Breadcrumbs';
	import MinimalMarkdownEditor from './Editor';
	import TitleField from './Metadata/TitleField';
	import Accessors from '../Mixins/accessors';

	export default {
		name: "DirectoryNew",
		data() {
			return {
				markdownLoaded: true,
				formID: "create-directory",
				valid: false,
				title: "Create directory"
			};
		},
		created() {
			this.$store.commit("initializeDirectory");
			this.$store.dispatch("initializeCommit");

			this.$bus.$on("checkMetadata", () => {
				this.validate()
			});
		},
		methods: {
			async createDirectory(event) {
				event.preventDefault();

				this.commit.addDirectory(this.activeDirectory);

				let response = await this.activeDirectory.create(this.commit);

				if (!checkResponse(response.status)) {
					console.error(response.status);
					return;
				};

				let json = await response.json();

				await this.$store.commit("setLatestRevision", json.oid);

				// new directory created successfully, show a message
				this.$store.state.broadcast.addMessage(
					"success",
					"Directory Created",
					`You have created the directory ${this.activeDirectory.title}, it has the path ${this.activeDirectory.path}`,
					3
				);

				// redirect to the new directory's index page
				this.redirectToIndex(this.activeDirectory.path);
				return;
			},
			redirectToIndex(directory) {
				this.$router.push({name: 'directory_index', params: {directory}});
			},
			validate() {
				if (!this.form) {
					this.form = document.getElementById(this.formID);
				};
				this.valid = this.form.checkValidity();
			}
		},
		computed: {
			breadcrumbs() {
				return [
					new CMSBreadcrumb(
						"New Directory",
						"directory_new",
						{}
					)
				];
			}
		},
		watch: {
			"activeDirectory.title": function title() {
				this.activeDirectory.path = slugify(this.activeDirectory.title);
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

<style lang="scss" scoped>
	form {
		margin: 1rem;
	};
</style>
