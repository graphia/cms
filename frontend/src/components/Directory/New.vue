<template>

	<div class="new-directory">

		<h4>
			Create a new directory
		</h4>

		<!-- new directory form -->
		<form @submit="createDirectory">

			<div class="form-group">
				<label for="title">Title</label>
				<input
					name="title"
					type="text"
					class="form-control"
					placeholder="Operating Procedures"
					v-model="directory.title"
					required="true"
				/>
			</div>

			<div class="form-group">
				<label for="path">Path</label>
				<input
					name="path"
					class="form-control"
					placeholder="operating-procedures"
					v-model="directory.path"
					required="true"
				/>
			</div>

			<div class="form-group">
				<label for="description">Description</label>
				<textarea
					name="description"
					class="form-control"
					v-model="directory.description"
					placeholder="A set of detailed step-by-step instructions compiled to help workers carry out complex routine operations"
				/>
			</div>

			<div>
				<MinimalMarkdownEditor/>
			</div>

			<div class="form-group">
				<input
					type="submit"
					class="form-control btn btn-success"
					value="Create directory"
				/>
			</div>

		</form>
		<!-- /new directory form -->

	</div>

</template>

<script lang="babel">
	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSDirectory from '../../javascripts/models/directory.js';
	import MinimalMarkdownEditor from './Editor';

	export default {
		name: "DirectoryNew",
		data() {
			return {
				//directory: new CMSDirectory(),
				markdownLoaded: true
			};
		},
		created() {
			this.$store.commit("initializeDirectory");
			this.$store.dispatch("initializeCommit");
		},
		computed: {
			directory() {
				return this.$store.state.activeDirectory;
			},
			commit() {
				return this.$store.state.commit;
			}
		},
		methods: {
			async createDirectory(event) {
				event.preventDefault();

				let response = await this.directory.create(this.commit);

				if (!checkResponse(response.status)) {
					console.error(response.status);
					return;
				}

				// new directory created successfully, show a message
				this.$store.state.broadcast.addMessage(
					"success",
					"Directory Created",
					`You have created the directory ${this.directory.title}, it has the path ${this.directory.path}`,
					3
				);

				// redirect to the new directory's index page
				this.redirectToIndex(this.directory.path);
				return;
			},
			redirectToIndex(directory) {
				this.$router.push({name: 'document_index', params: {directory}});
			},
		},
		components: {
			MinimalMarkdownEditor
		}
	};
</script>