<template>
	<div>

		<div v-if="numberOfDirectories == 0">
			{{ directories.length }}
			<p>No directories</p>

		</div>

		<!-- listing directories -->
		<div class="row" v-else-if="numberOfDirectories > 0">

			<div  class="col-lg-4 mt-4"  v-for="(contents, directory) in directories">

				<div class="card">

					<h4 class="card-header">
						<router-link :to="{name: 'document_index', params: {directory: directory}}">
								{{ directory | capitalize }}
						</router-link>
					</h4>

					<!-- listing documents inside a directory -->
					<div class="list-group list-group-flush" v-if="contents.length > 0">

						<router-link
							v-for="document in contents"
							class="list-group-item list-group-item-action"
							:to="{name: 'document_show', params: {directory: directory, filename: document.filename}}"
						>

							{{ document.frontmatter.title }}

						</router-link>

					</div>
					<!-- /listing documents inside a directory -->

					<div class="card-body" v-else-if="contents.length == 0">

						<div class="alert alert-info">
							There's nothing here yet
						</div>

						<router-link class="btn btn-sm btn-primary" :to="{name: 'document_new', params: {directory: directory}}">
							Create a document
						</router-link>

					</div>

				</div>
			</div>

			<!-- new directory form -->
			<div class="new-directory col-md-4">

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
			<!-- /new directory form -->


		</div>
		<!-- /listing directories -->



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
			async createDirectory(event) {
				event.preventDefault();

				let response = await this.directory.create(this.commit);

				if (!checkResponse(response)) {
					console.error(response);
					return;
				}

				// new directory created successfully, show a message
				this.$store.state.broadcast.addMessage(
					"success",
					"Welcome",
					`created directory ${this.directory.path}`,
					3
				);

				// refresh the dir list and initialise a new dir for form
				this.fetchDirectorySummary();
				this.directory = new CMSDirectory()
				return;
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