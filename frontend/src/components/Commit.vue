<template>
	<div>
		<h1>Commit</h1>

		<pre v-html="this.full_diff">
		</pre>

		<ol>
			<li v-for="(item, key, index) in files">
				<CommitFile :path='key' :files='item'/>
			</li>
		</ol>


	</div>

</template>

<script lang="babel">
	import config from '../javascripts/config.js';
	import store from '../javascripts/store.js';
	import checkResponse from '../javascripts/response.js';
	import CommitFile from './Commit/File';

	export default {
		name: "Commit",
		data() {
			return {
				commit: {},
				full_diff: null
			}
		},
		computed: {
			hash() {
				return this.$route.params.hash;
			},
			files() {
				return this.commit.files;
			}
		},
		async created() {
			console.debug("created");
			await this.retrievePatch(this.hash);
			this.setupDiff();
		},
		components: {
			CommitFile
		},
		methods: {
			async retrievePatch() {
				console.debug("got here");

				let path = `${config.api}/commits/${this.hash}`;

				let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()});

				if (!checkResponse(response.status)) {
					throw("Could not retrieve changeset");
				}

				let json = await response.json()

				this.commit = json;


			},
			setupDiff() {
				console.log("setting up diff");


			}
		}
	};
</script>