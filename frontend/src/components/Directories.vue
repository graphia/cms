<template>
	<div>
		<h1>New Directory</h1>


		<div v-if="numberOfDirectories == 0">
			{{ directories.length }}
			<p>No directories</p>

			<button class="btn btn-primary" @click="createDirectory">
				Create Directory
			</button>
		</div>

		<ul v-else-if="numberOfDirectories > 0">
			<li v-for="(contents, directory) in directories">
				{{ directory }}
			</li>
		</ul>

	</div>
</template>

<script lang="babel">

	import checkResponse from '../javascripts/response.js';
	import config from '../javascripts/config.js';
	import __Object from 'babel-runtime/core-js/object/keys';

	//import _object from '../../../node_modules/babel-runtime/core-js/object/keys';

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
			createDirectory(event) {
				event.preventDefault();

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