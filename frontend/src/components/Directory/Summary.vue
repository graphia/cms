<template>
	<div>

		<div class="col col-md-6" v-if="loading">
			<div class="alert alert-warning mx-auto">
				Loading...
			</div>
		</div>

		<div class="col col-md-6" v-else-if="directories.length == 0">
			<div class="jumbotron">
				{{ directories.length }}
				<p>No directories</p>

				<DirectoryNewButton/>
			</div>
		</div>

		<div v-else>

			<!-- listing directories -->
				<div :data-directory="directory.path" :class="directory.path" v-for="(directory, i) in directories" :key="i">
					<h4>
						<router-link :to="{name: 'document_index', params: {directory: directory.path}}">
							{{ (directory.title || directory.path) | capitalize }}
						</router-link>
					</h4>

					<p class="directory-description">{{ directory.description }}</p>

					<IndexList :documents="directory.contents" :includeNewButton="true"/>

				</div>

				<div class="col-lg-4 mt-4">
					<DirectoryNewButton/>
				</div>

			<!-- /listing directories -->
		</div>

	</div>
</template>

<script lang="babel">

	import 'babel-runtime/core-js/object/keys';

	import checkResponse from '../../javascripts/response.js';
	import config from '../../javascripts/config.js';
	import CMSFile from '../../javascripts/models/file.js';
	import CMSDirectory from '../../javascripts/models/directory.js';
	import IndexList from '../Document/Index/List';

	import DirectoryNewButton from './NewButton';

	export default {
		name: "DirectorySummary",
		data() {
			return {
				directories: {},
				loading: true
			}
		},
		created() {
			this.fetchDirectorySummary();
			this.$store.commit("refreshServerInfo");
		},
		computed: {
			numberOfDirectories() {
				let count = Object.keys(this.directories).length;
				return count;
			}
		},
		methods: {
			async fetchDirectorySummary() {

				let path = `${config.api}/summary`

				let response = await fetch(path, {
					method: "GET",
					headers: this.$store.state.auth.authHeader()
				});

				if (!checkResponse(response.status)) {
					console.error(response);
					return;
				}

				let json = await response.json();

				this.directories = json.map((dir) => {

					return new CMSDirectory(
						dir.path,
						dir.info.title,
						dir.info.description,
						dir.info.body,
						dir.contents.map((file) => {
							return new CMSFile(file);
						})
					);
				});

				this.loading = false;

			}
		},
		components: {
			DirectoryNewButton,
			IndexList
		}
	}
</script>

<style lang="scss">
	.translations > .translations-list {
		margin-bottom: 0px;
	}
</style>