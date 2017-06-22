<template>
	<div>
		<h1>History</h1>

		<ol>
			<li v-for="item in history">
				{{ item.message }}
			</li>
		</ol>
	</div>
</template>

<script lang="babel">

	import checkResponse from '../javascripts/response.js';

	export default {
		name: "DocumentHistory",

		data() {
			return {
				history: []
			};
		},

		computed: {
			directory() {
				return this.$route.params.directory;
			},
			filename() {
				return this.$route.params.filename;
			}
		},

		async created() {

			var directory = this.directory;
			var filename = this.filename;

			if (!this.$store.state.activeDocument.populated()) {
				await this.$store.dispatch("getDocument", {directory, filename});
			};

			let response = await this.$store.state.activeDocument.log();

			if (!checkResponse(response.status)) {
				throw("Could not retrieve history");
			}

			this.history = await response.json()

		}
	};
</script>