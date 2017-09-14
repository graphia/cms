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
					class="form-control"
					placeholder="Operating Procedures"
					v-model="directory.title"
				/>
			</div>

			<div class="form-group">
				<label for="path">Path name</label>
				<input
					name="path"
					class="form-control"
					placeholder="operating-procedures"
					v-model="directory.path"
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

			<div class="form-group">
				<input
					type="submit"
					class="form-control btn btn-success"
					value="Create Directory"
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

	export default {
		name: "DirectoryNew",
		data() {
			return {
				directory: new CMSDirectory()
			};
		},
		created() {
			this.$store.dispatch("initializeCommit");
		},
		computed: {
			commit() {
				return this.$store.state.commit;
			},
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
					"Welcome",
					`created directory ${this.directory.path}`,
					3
				);

				// redirect to the new directory's index page
				this.redirectToIndex(this.directory.path);
				return;
			},
			redirectToIndex(directory) {
				this.$router.push({name: 'document_index', params: {directory}});
			},
		}
	};
</script>