<template>
	<div>
		<h1>New Directory</h1>


		<div v-if="numberOfDirectories == 0">
			{{ directories.length }}
			<p>No directories</p>

		</div>

		<div v-else-if="numberOfDirectories > 0">



			<div class="card" v-for="(contents, directory) in directories">

				<div class="card-header">
					<h3 class="card-title">{{ directory | capitalize }}</h3>
				</div>

				<div class="list-group list-group-flush">

					<router-link
						v-for="document in contents"
						class="list-group-item list-group-item-action"
						:to="{name: 'document_show', params: {directory: directory, filename: document.filename}}"
					>
						{{ document.frontmatter.title }}

					</router-link>

				</div>

			</div>
		</div>

		<div class="new-directory">

			<form>

				<div class="input-group">
					<input
						class="form-control"
						placeholder="stories"
						v-model="directory.path"
					/>

					<span class="input-group-btn">
						<input
							type="submit"
							class="form-control btn btn-success"
							value="Create Directory"
							@click="createDirectory"
						/>
					</span>

				</div>
			</form>

		</div>

	</div>
</template>

<script lang="babel">

	import checkResponse from '../javascripts/response.js';
	import config from '../javascripts/config.js';
	import CMSDirectory from '../javascripts/models/directory.js';
	import 'babel-runtime/core-js/object/keys';

	export default {
		name: "Directories",
		data() {
			return {
				directories: {},
				directory: new CMSDirectory()
			}
		},
		created() {
			this.fetchDirectorySummary();

			// set up a fresh new commit
			this.$store.dispatch("initializeCommit");

		},
		computed: {
			numberOfDirectories() {
				let count = Object.keys(this.directories).length;
				console.debug("directory count", count);
				return count;
			},
			commit() {
				return this.$store.state.commit;
			},
		},
		methods: {
			createDirectory(event) {
				event.preventDefault();

				this.directory.create(this.commit);

				console.debug("clicked");
			},
			async fetchDirectorySummary() {

				let path = `${config.api}/summary`

				console.log("fetching directories", path);

				try {
					let response = await fetch(path, {
						mode: "cors",
						method: "GET",
						headers: this.$store.state.auth.authHeader()
					});

					if (!checkResponse(response)) {
						console.error(response);
						return;
					}

					let json = await response.json();
					console.log("got json", json)

					// TODO map the directories into CMSFile objects
					this.directories = json;
					console.log("directories", json);


				}
				catch(error) {
					console.error(error);
				}
			}
		}
	}
</script>