<template>
	<div>

		<div v-if="numberOfDirectories == 0">
			{{ directories.length }}
			<p>No directories</p>

		</div>

		<!-- listing directories -->
		<div class="row" v-else-if="numberOfDirectories > 0">

			<div class="col-lg-4 mt-4" v-for="directory in directories">

				<div class="card" :class="directory.path" >

					<h4 class="card-header">
						<router-link :to="{name: 'document_index', params: {directory: directory.path}}">
							{{ (directory.info.title || directory.path) | capitalize }}
						</router-link>
					</h4>

					<div class="card-body">
						{{ directory.info.description }}
					</div>

					<!-- listing documents inside a directory -->
					<div class="list-group list-group-flush" v-if="directory.contents.length > 0">

						<router-link
							v-for="document in directory.contents"
							class="list-group-item list-group-item-action"
							:to="{name: 'document_show', params: {directory: directory.path, filename: document.filename}}"
							:data-filename="document.filename"
						>

							{{ document.frontmatter.title }}

						</router-link>

						<div class="card-body">
							<router-link class="btn btn-sm btn-primary" :to="{name: 'document_new', params: {directory: directory}}">
								Create a document
							</router-link>
						</div>


					</div>
					<!-- /listing documents inside a directory -->

					<div class="card-body" v-else-if="directory.contents.length == 0">

						<div class="alert alert-info">
							There's nothing here yet
						</div>

						<router-link class="btn btn-sm btn-primary" :to="{name: 'document_new', params: {directory: directory}}">
							Create a document
						</router-link>

					</div>

				</div>
			</div>

			<DirectoryNew/>

		</div>
		<!-- /listing directories -->

	</div>
</template>

<script lang="babel">

	import checkResponse from '../javascripts/response.js';
	import config from '../javascripts/config.js';
	import CMSDirectory from '../javascripts/models/directory.js';
	import 'babel-runtime/core-js/object/keys';
	import DirectoryNew from '../components/DirectoryNew';

	export default {
		name: "Directories",
		data() {
			return {
				directories: {}
			}
		},
		created() {
			this.fetchDirectorySummary();
		},
		computed: {
			numberOfDirectories() {
				let count = Object.keys(this.directories).length;
				console.debug("directory count", count);
				return count;
			}
		},
		methods: {
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
		},
		components: {
			DirectoryNew
		}
	}
</script>