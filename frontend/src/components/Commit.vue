<template>
	<div>
		<h1>Commit</h1>

		<pre>
			{{ this.commit.diff }}
		</pre>
	</div>

</template>

<script lang="babel">
	import config from '../javascripts/config.js';
	import store from '../javascripts/store.js';

	export default {
		name: "Commit",
		data() {
			return {
				commit: {}
			}
		},
		computed: {
			hash() {
				return this.$route.params.hash;
			}
		},
		created() {
			console.debug("created");
			this.retrievePatch(this.hash);

		},
		methods: {
			async retrievePatch() {
				console.debug("got here");

				let path = `${config.api}/commits/${this.hash}`;

				let response = await fetch(path, {mode: "cors", headers: store.state.auth.authHeader()});


				let json = await response.json()

				this.commit = json;


			}
		}
	};
</script>

<style lang="scss">

</style>