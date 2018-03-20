<template>

	<div class="edit-directory" v-title="title">

		<h4>
			Edit {{ this.activeDirectory.title }}
		</h4>

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

			<div class="form-group">
				<input
					type="submit"
					class="form-control btn btn-success"
					value="Update directory"
				/>
			</div>

		</form>
		<!-- /edit directory form -->

	</div>

</template>

<script lang="babel">
	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSDirectory from '../../javascripts/models/directory.js';

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
			}
		},
		components: {
			MinimalMarkdownEditor,
			TitleField
		},
		mixins: [Accessors]
	};
</script>